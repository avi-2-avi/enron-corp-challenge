package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"indexing/models"
	"indexing/worker"
	batch "indexing/zincsearch"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const (
	maxBatchLines = 1000
	apiURL        = "http://localhost:4080"
	numWorkers    = 10
)

func main() {
	startTime := time.Now()

	prof := flag.Bool("prof", false, "Start pprof server")
	flag.Parse()

	if *prof {
		go startPprofServer()
	}

	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: ./indexer <directory_path>")
		return
	}
	root := flag.Arg(0)
	fmt.Println("Root directory:", root)

	auth := base64.StdEncoding.EncodeToString([]byte("admin:Pass123!!!"))
	authHeader := "Basic " + auth

	// Create index
	indexURL := apiURL + "/api/index"
	batch.IndexCreator(indexURL, authHeader)

	fmt.Println("Starting to index the files...")

	var wg sync.WaitGroup
	filePaths := make(chan string, 100)        // Buffer size
	results := make(chan models.Document, 100) // Buffer size, docs
	done := make(chan struct{})

	worker.StartWokers(numWorkers, &wg, filePaths, results)
	go walkFiles(root, filePaths)
	go processResults(results, done, authHeader, maxBatchLines)
	go worker.WaitForWorkers(&wg, results)

	<-done

	elapsedTime := time.Since(startTime)
	fmt.Printf("Done sending the data to ZincSearch. Total execution time: %s\n", elapsedTime)
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

func processResults(results <-chan models.Document, done chan<- struct{}, authHeader string, maxBatchLines int) {
	url := apiURL + "/api/emails/_multi"

	batchList := make([]models.Document, 0, maxBatchLines)

	for result := range results {
		batchList = append(batchList, result)

		if len(batchList) >= maxBatchLines {
			err := batch.SendBatch(batchList, url, authHeader)
			if err != nil {
				fmt.Println("error sending batch:", err)
			}
			batchList = batchList[:0]
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

func startPprofServer() {
	fmt.Println("Starting pprof server on http://localhost:6060")
	if err := http.ListenAndServe("localhost:6060", nil); err != nil {
		fmt.Println("error starting pprof server:", err)
	}
}
