package zincsearch

import (
	"fmt"
	"net/http"
)

func setHeaders(req *http.Request, authHeader, contentType string) {
	req.Header.Set("Authorization", authHeader)
	req.Header.Set("Content-Type", contentType)
	req.Header.Set("Accept", "application/json")
}

func handleResponse(resp *http.Response) error {
	if resp.StatusCode == http.StatusOK {
		return nil
	}
	return fmt.Errorf("error response from server: %v", resp.Status)
}
