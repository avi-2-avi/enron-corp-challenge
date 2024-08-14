package main

import (
	"fmt"
	"indexing/internal/database"
	"indexing/internal/models"
	"indexing/internal/utils"
	_ "net/http/pprof"
	"sync"
	"time"
)

const (
	maxBatchLines = 1000
	numWorkers    = 10
)

func main() {
	startTime := time.Now()

	root := utils.GetRootDirectory()
	fmt.Println("Root directory found:", root)

	authHeader := utils.GetAuthHeader()

	database.CreateIndex(authHeader)

	filePaths := make(chan string, 200)                  // Buffer size for file paths
	results := make(chan models.Document, maxBatchLines) // Buffer size for batch
	done := make(chan struct{})

	var wg sync.WaitGroup
	fmt.Println("Starting to index the files...")

	utils.StartWokers(numWorkers, &wg, filePaths, results)
	go utils.WalkFiles(root, filePaths)
	go utils.ProcessResults(results, done, authHeader, maxBatchLines)
	go utils.WaitForWorkers(&wg, results)

	<-done

	elapsedTime := time.Since(startTime)
	fmt.Printf("Done sending the data to ZincSearch. Total execution time: %s\n", elapsedTime)
}
