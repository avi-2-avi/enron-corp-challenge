package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"indexing/batch"
	"indexing/models"
	"indexing/worker"
	"os"
	"path/filepath"
	"sync"
)

const maxBatchLines = 1000

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: ./indexer <directory_path>")
		return
	}
	root := flag.Arg(0)
	fmt.Println("Root directory:", root)

	var wg sync.WaitGroup
	filePaths := make(chan string, 100)        // Buffer size
	results := make(chan models.Document, 100) // Buffer size
	done := make(chan struct{})

	const numWorkers = 10

	worker.StartWokers(numWorkers, &wg, filePaths, results)
	go worker.WaitForWorkers(&wg, results)
	go walkFiles(root, filePaths)
	go processResults(results, done)

	<-done
	fmt.Println("Done sending the data to ZincSearch")
}

func walkFiles(root string, filePaths chan<- string) {
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
}

func processResults(results <-chan models.Document, done chan<- struct{}) {
	url := "http://localhost:4080/api/emails/_multi"
	auth := base64.StdEncoding.EncodeToString([]byte("admin:Pass123!!!"))
	authHeader := "Basic " + auth

	var batchList []models.Document
	var batchSize int

	for result := range results {
		batchList = append(batchList, result)
		batchSize++

		if batchSize > maxBatchLines {
			err := batch.SendBatch(batchList, url, authHeader)
			if err != nil {
				fmt.Println("error sending batch:", err)
			}
			batchList = batchList[:0]
			batchSize = 0
		}
	}

	if len(batchList) > 0 {
		err := batch.SendBatch(batchList, url, authHeader)
		if err != nil {
			fmt.Println("error sending final batch:", err)
		}
	}

	done <- struct{}{}
}
