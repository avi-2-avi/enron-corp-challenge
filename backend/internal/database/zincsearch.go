package database

import (
	"backend/internal/errors"
	"bytes"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var zincBaseURL string
var zincUsername string
var zincPassword string

func InitZincSearch() {
	zincBaseURL = os.Getenv("ZINC_BASE_URL")
	zincUsername = os.Getenv("ZINC_USERNAME")
	zincPassword = os.Getenv("ZINC_PASSWORD")

	if zincBaseURL == "" || zincUsername == "" || zincPassword == "" {
		log.Fatal("ZincSearch configuration is missing")
	}
}

func GetIndexByID(index string, id string) ([]byte, *errors.AppError) {
	url := fmt.Sprintf("%s/api/%s/_doc/%s", zincBaseURL, index, id)

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return nil, errors.NewAppError("Failed to create request", http.StatusInternalServerError, err)
	}

	req.SetBasicAuth(zincUsername, zincPassword)
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewAppError("Failed to execute request", http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.NewAppError(fmt.Sprintf("ZincSearch error: %s", string(body)), resp.StatusCode, nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewAppError("Failed to read response body", http.StatusInternalServerError, err)
	}

	return body, nil
}

func GetIndexDocuments(index string, query []byte) ([]byte, *errors.AppError) {
	url := fmt.Sprintf("%s/api/%s/_search", zincBaseURL, index)

	req, err := http.NewRequest("POST", url, bytes.NewBuffer([]byte(query)))
	if err != nil {
		return nil, errors.NewAppError("Failed to create request", http.StatusInternalServerError, err)
	}

	req.SetBasicAuth(zincUsername, zincPassword)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.NewAppError("Failed to execute request", http.StatusInternalServerError, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, errors.NewAppError(fmt.Sprintf("ZincSearch error: %s", string(body)), resp.StatusCode, nil)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.NewAppError("Failed to read response body", http.StatusInternalServerError, err)
	}

	return body, nil
}
