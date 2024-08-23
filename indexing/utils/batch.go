package utils

import (
	"fmt"
	"indexing/models"
	zinc "indexing/zincsearch"
)

func ProcessResults(results <-chan models.Document, done chan<- struct{}, authHeader string, maxBatchLines int, apiURL string) {
	url := apiURL + "/api/emails/_multi"

	batchList := make([]models.Document, 0, maxBatchLines)

	for result := range results {
		batchList = append(batchList, result)

		if len(batchList) >= maxBatchLines {
			err := zinc.SendBatch(batchList, url, authHeader)
			if err != nil {
				fmt.Println("error sending batch:", err)
			}
			batchList = batchList[:0]
		}
	}

	if len(batchList) > 0 {
		err := zinc.SendBatch(batchList, url, authHeader)
		if err != nil {
			fmt.Println("error sending final batch:", err)
		}
	}

	done <- struct{}{}
}
