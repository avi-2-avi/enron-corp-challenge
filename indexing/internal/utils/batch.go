package utils

import (
	"fmt"
	"indexing/internal/database"
	"indexing/internal/models"
)

func ProcessResults(results <-chan models.Document, done chan<- struct{}, authHeader string, maxBatchLines int) {
	batchList := make([]models.Document, 0, maxBatchLines)

	for result := range results {
		batchList = append(batchList, result)

		if len(batchList) >= maxBatchLines {
			if err := database.SendBatch(batchList, authHeader); err != nil {
				fmt.Println("error sending batch:", err)
			}

			batchList = batchList[:0]
		}
	}

	if len(batchList) > 0 {
		if err := database.SendBatch(batchList, authHeader); err != nil {
			fmt.Println("error sending final batch:", err)
		}
	}

	done <- struct{}{}
}
