package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"indexing/models"
	"net/http"
	"time"
)

const (
	maxRetries = 5
	retryDelay = time.Second * 2
)

func SendBatch(batch []models.Document, url, authHeader string) error {
	var buffer bytes.Buffer
	for _, doc := range batch {
		jsonDoc, err := json.Marshal(doc.Data)
		if err != nil {
			return fmt.Errorf("error marshalling document: %v", err)
		}
		buffer.WriteString(string(jsonDoc) + "\n")
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {

		req, err := http.NewRequest("POST", url, &buffer)
		if err != nil {
			return fmt.Errorf("error creating HTTP request: %v", err)
		}

		req.Header.Set("Authorization", authHeader)
		req.Header.Set("Content-Type", "text/plain")
		req.Header.Set("Accept", "application/json")

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			fmt.Printf("Attempt %d: Error sending HTTP request: %v\n", attempt, err)
			time.Sleep(retryDelay)
			continue
		}
		defer resp.Body.Close()

		if resp.StatusCode == http.StatusOK {
			return nil
		}

		fmt.Printf("Attempt %d: Error response from server: %v\n", attempt, resp.Status)
		fmt.Println("Retry in", retryDelay)
		time.Sleep(retryDelay)
	}
	return fmt.Errorf("failed to send batch after %d attempts", maxRetries)
}

func IndexCreator(url, authHeader string) {
	payload := map[string]interface{}{
		"name":         "emails",
		"storage_type": "disk",
		"shard_num":    1,
		"mappings": map[string]interface{}{
			"properties": map[string]interface{}{
				"date": map[string]interface{}{
					"type":  "text",
					"index": true,
				},
				"from": map[string]interface{}{
					"type":     "keyword",
					"index":    true,
					"sortable": true,
				},
				"to": map[string]interface{}{
					"type":     "keyword",
					"index":    true,
					"sortable": true,
				},
				"subject": map[string]interface{}{
					"type":     "keyword",
					"index":    true,
					"sortable": true,
				},
				"content": map[string]interface{}{
					"type":          "text",
					"index":         true,
					"highlightable": true,
				},
				"path": map[string]interface{}{
					"type":  "text",
					"index": true,
				},
			},
		},
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		fmt.Println("error marshalling payload to JSON: ", err)
		return
	}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("error creating HTTP request: ", err)
		return
	}

	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending HTTP request: ", err)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("Index created successfully")
	} else {
		fmt.Println("error response from server: ", resp.Status)
		return
	}
}
