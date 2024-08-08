package worker

import (
	"fmt"
	"indexing/models"
	"sync"
)

func WaitForWorkers(wg *sync.WaitGroup, results chan models.Document) {
	wg.Wait()
	close(results)
	fmt.Println("All workers have finished")
}
