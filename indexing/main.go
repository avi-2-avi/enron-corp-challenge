package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxBatchLines = 1000
	maxRetries    = 5
	retryDelay    = time.Second * 2
)

type Document struct {
	Index map[string]string `json:"index"`
	Data  map[string]string `json:"data"`
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: ./indexer <directory_path>")
		return
	}
	root := flag.Arg(0)
	fmt.Println("Root directory:", root)

	var wg sync.WaitGroup
	filePaths := make(chan string, 100) // Buffer size
	results := make(chan Document, 100) // Buffer size
	done := make(chan struct{})

	const numWorkers = 10

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go Worker(i, &wg, filePaths, results)
	}

	go func() {
		wg.Wait()
		close(results)
		fmt.Println("All workers have finished")
	}()

	go func() {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("error walking through the directory:", err)
				return err
			}

			if !info.IsDir() {
				filePaths <- path
			}
			return nil
		})
		if err != nil {
			fmt.Println("error walking through the directory:", err)
		}

		close(filePaths)
	}()

	go func() {
		url := "http://localhost:4080/api/emails/_multi"
		auth := base64.StdEncoding.EncodeToString([]byte("admin:Pass123!!!"))
		authHeader := "Basic " + auth

		var batch []Document
		var batchSize int

		for result := range results {
			batch = append(batch, result)
			batchSize++

			if batchSize > maxBatchLines {
				err := SendBatch(batch, url, authHeader)
				if err != nil {
					fmt.Println("error sending batch:", err)
				}
				batch = batch[:0]
				batchSize = 0
			}
		}

		if len(batch) > 0 {
			err := SendBatch(batch, url, authHeader)
			if err != nil {
				fmt.Println("error sending final batch:", err)
			}
		}

		done <- struct{}{}
	}()

	<-done
	fmt.Println("Done sending the data to ZincSearch")
}
