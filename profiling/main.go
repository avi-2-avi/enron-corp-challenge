package main

import (
	"bufio"
	"bytes"
	"encoding/base64"
	"encoding/json"
	"flag"
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

const (
	maxBatchLines = 1000
	maxRetries    = 5
	retryDelay    = time.Second * 2
)

type Document struct {
	Index map[string]string `json:"index"`
	Data  map[string]string `json:"data"`
}

func worker(id int, wg *sync.WaitGroup, filePaths <-chan string, results chan<- Document) {
	defer wg.Done()

	for path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("Worker %d error opening file: %v\n", id, err)
			continue
		}

		reader := bufio.NewReader(file)

		var foundNewLine bool
		var date, from, to, subject, content, fileName string

		fileName = filepath.Base(path)

		if fileName == ".DS_Store" {
			file.Close()
			continue
		}

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			line = strings.TrimSpace(line)

			if strings.HasPrefix(line, "Date:") {
				parts := strings.SplitN(line, "Date: ", 2)
				if len(parts) > 1 {
					date = parts[1]
				}
			} else if strings.HasPrefix(line, "From:") && !strings.HasPrefix(line, "X-From: ") {
				parts := strings.SplitN(line, "From: ", 2)
				if len(parts) > 1 {
					from = parts[1]
				}
			} else if strings.HasPrefix(line, "To:") && !strings.HasPrefix(line, "X-To: ") {
				parts := strings.SplitN(line, "To: ", 2)
				if len(parts) > 1 {
					to = parts[1]
				}
			} else if strings.HasPrefix(line, "Subject:") {
				parts := strings.SplitN(line, "Subject: ", 2)
				if len(parts) > 1 {
					subject = parts[1]
				}
			}

			if foundNewLine {
				content += line + "\n"
			} else {
				if line == "" {
					foundNewLine = true
				}
			}
		}

		indexNdLine := map[string]string{"_index": "emails"}
		dataNdLine := map[string]string{"path": path, "date": date, "from": from, "to": to, "subject": subject, "content": content}
		doc := Document{Index: indexNdLine, Data: dataNdLine}

		results <- doc

		file.Close()
	}
}

func sendBatch(batch []Document, url, authHeader string) error {
	var buffer bytes.Buffer
	for _, doc := range batch {
		jsonDoc, err := json.Marshal(doc.Data)
		if err != nil {
			return fmt.Errorf("Error marshalling document: %v", err)
		}
		buffer.WriteString(string(jsonDoc) + "\n")
	}

	for attempt := 1; attempt <= maxRetries; attempt++ {

		req, err := http.NewRequest("POST", url, &buffer)
		if err != nil {
			return fmt.Errorf("Error creating HTTP request: %v", err)
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
	return fmt.Errorf("Failed to send batch after %d attempts", maxRetries)
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: ./indexer <directory_path>")
		return
	}
	root := flag.Arg(0)
	fmt.Println("Root directory:", root)

	go func() {
		fmt.Println(http.ListenAndServe("localhost:6060", nil))
	}()

	var wg sync.WaitGroup
	filePaths := make(chan string, 100) // Buffer size
	results := make(chan Document, 100) // Buffer size
	done := make(chan struct{})

	const numWorkers = 1

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go worker(i, &wg, filePaths, results)
	}

	go func() {
		wg.Wait()
		close(results)
		fmt.Println("All workers have finished")
	}()

	go func() {
		err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error walking through the directory:", err)
				return err
			}

			if !info.IsDir() {
				filePaths <- path
			}
			return nil
		})
		if err != nil {
			fmt.Println("Error walking through the directory:", err)
		}

		close(filePaths)
	}()

	go func() {
		url := "http://localhost:4080/api/emails/_multi"
		auth := base64.StdEncoding.EncodeToString([]byte("admin:Pass123!!!"))
		authHeader := "Basic " + auth

		var batch []Document
		var batchSize int

		for result := range results {
			batch = append(batch, result)
			batchSize++

			if batchSize > maxBatchLines {
				err := sendBatch(batch, url, authHeader)
				if err != nil {
					fmt.Println("Error sending batch:", err)
				}
				batch = batch[:0]
				batchSize = 0
			}
		}

		if len(batch) > 0 {
			err := sendBatch(batch, url, authHeader)
			if err != nil {
				fmt.Println("Error sending final batch:", err)
			}
		}

		done <- struct{}{}
	}()

	<-done
	fmt.Println("Done sending the data to ZincSearch")
}
