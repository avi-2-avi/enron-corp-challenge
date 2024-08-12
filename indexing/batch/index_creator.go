package batch

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

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
