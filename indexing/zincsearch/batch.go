package zincsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"indexing/models"
	"net/http"
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

	req, err := http.NewRequest("POST", url, &buffer)
	if err != nil {
		return fmt.Errorf("error creating HTTP request: %v", err)
	}

	setHeaders(req, authHeader, "text/plain")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Printf("Error sending HTTP request: %v\n", err)
		return err
	}
	defer resp.Body.Close()

	if err := handleResponse(resp); err != nil {
		fmt.Printf("failed %v\nRetrying inmediately...", resp.Status)
	}
	return nil
}
