package database

import (
	"bytes"
	"encoding/json"
	"fmt"
	"indexing/internal/models"
	"net/http"
)

// General request functions

func setRequestHeaders(req *http.Request, authHeader string, contentType string) {
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
}

func createRequest(url string, payload []byte, authHeader, contentType string) (*http.Request, error) {
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		return nil, fmt.Errorf("error creating HTTP request: %w", err)
	}

	setRequestHeaders(req, authHeader, contentType)
	return req, nil
}

func handleResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("error response from server: %v", resp.Status)
}

// Index specific functions

func createIndexPayload() ([]byte, error) {
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
		return nil, fmt.Errorf("error marshalling payload to JSON: %w", err)
	}
	return jsonData, nil
}

// Document specific functions

func prepareRequestBody(batch []models.Document) (bytes.Buffer, error) {
	var buffer bytes.Buffer

	for _, doc := range batch {
		jsonDoc, err := json.Marshal(doc.Data)
		if err != nil {
			return buffer, fmt.Errorf("error marshalling document: %v", err)
		}
		buffer.WriteString(string(jsonDoc) + "\n")
	}
	return buffer, nil
}

func sendRequestWithRetry(buffer bytes.Buffer, authHeader string, attempt int) error {
	url := apiURL + "/api/emails/_multi"
	payload := buffer.Bytes()

	req, err := createRequest(url, payload, authHeader, "text/plain")
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}
	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		return fmt.Errorf("attempt %d failed: %v", attempt, err)
	}

	return nil
}
