package database

import (
	"fmt"
	"indexing/internal/models"
	"net/http"
)

const (
	maxRetries = 3
	apiURL     = "http://localhost:4080"
	indexURL   = apiURL + "/api/index"
)

func CreateIndex(authHeader string) error {
	payload, err := createIndexPayload()
	if err != nil {
		return fmt.Errorf("error marshalling payload to JSON: %v", err)
	}

	req, err := createRequest(apiURL, payload, authHeader, "application/json")
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("error sending HTTP request: %v", err)
	}

	defer resp.Body.Close()

	return handleResponse(resp)
}

func SendBatch(batch []models.Document, authHeader string) error {
	buffer, err := prepareRequestBody(batch)
	if err != nil {
		return fmt.Errorf("error preparing request body: %v", err)
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {
		err := sendRequestWithRetry(buffer, authHeader, attempt)
		if err == nil {
			return nil
		}
	}
	return fmt.Errorf("failed to send batch with %d files after %d attempts", len(batch), maxRetries)
}
