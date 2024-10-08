package worker

import (
	"indexing/models"
	"sync"
)

func StartWokers(numWorkers int, wg *sync.WaitGroup, filePaths chan string, results chan models.Document) {
	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go ProcessEmailFile(i, wg, filePaths, results)
	}
}
