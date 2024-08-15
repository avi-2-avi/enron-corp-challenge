package zincsearch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

func CreateIndex(url, authHeader string) {
	payload, _ := createIndexPayload()

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(payload))
	if err != nil {
		fmt.Println("error creating HTTP request: ", err)
		return
	}

	setHeaders(req, authHeader, "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending HTTP request: ", err)
		return
	}

	defer resp.Body.Close()

	handleResponse(resp)
}

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
