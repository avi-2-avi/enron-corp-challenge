package main

import (
	"fmt"
	"indexing/models"
	"indexing/utils"
	"indexing/worker"
	zinc "indexing/zincsearch"
	_ "net/http/pprof"
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

	root := utils.GetRootDirectory()
	fmt.Println("Root directory:", root)

	authHeader := utils.GetAuthHeader()

	indexURL := apiURL + "/api/index"
	zinc.CreateIndex(indexURL, authHeader)

	fmt.Println("Starting to index the files...")
	var wg sync.WaitGroup

	filePaths := make(chan string, 100)        // Buffer size
	results := make(chan models.Document, 100) // Buffer size, docs
	done := make(chan struct{})

	worker.StartWokers(numWorkers, &wg, filePaths, results)
	go utils.WalkFiles(root, filePaths)
	go utils.ProcessResults(results, done, authHeader, maxBatchLines, apiURL)
	go worker.WaitForWorkers(&wg, results)

	<-done

	elapsedTime := time.Since(startTime)
	fmt.Printf("Done sending the data to ZincSearch. Total execution time: %s\n", elapsedTime)
}
