/*
* Input: Directory path containing email files in maildir format
* Input structure: /maildir/<user>/<folder>/<email_file>
* Output: NDJSON file containing the formatted email data
 */

package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

func worker(wg *sync.WaitGroup, filePaths <-chan string, results chan<- map[string]interface{}) {
	defer wg.Done()

	indexNdLine := map[string]interface{}{"index": map[string]string{"_index": "emails"}}

	for path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Println("Error opening file:", err)
			continue
		}
		defer file.Close()

		reader := bufio.NewReader(file)

		var foundNewLine bool
		var date, from, to, subject, content, fileName string

		fileName = strings.Split(path, "/")[len(strings.Split(path, "/"))-1]

		if fileName == ".DS_Store" {
			file.Close()
			continue
		}

		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				file.Close()
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

		newNdLine := map[string]interface{}{"path": path, "date": date, "from": from, "to": to, "subject": subject, "content": content}

		results <- indexNdLine
		results <- newNdLine

		file.Close()
	}
}

func main() {
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("Usage: ./indexer <directory_path>")
		return
	}
	root := flag.Arg(0) + "/maildir"
	println("Root directory:", root)

	var wg sync.WaitGroup
	filePaths := make(chan string, 100)               // Buffer size
	results := make(chan map[string]interface{}, 100) // Buffer size
	done := make(chan struct{})

	// One worker node since with more than one node, the order of output writing is not guaranteed
	wg.Add(1)
	go worker(&wg, filePaths, results)

	go func() {
		wg.Wait()
		close(results)
	}()

	go func() {
		filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				fmt.Println("Error walking through the directory:", err)
				return err
			}

			if !info.IsDir() {
				filePaths <- path
			}
			return nil
		})

		close(filePaths)
	}()

	go func() {
		outputFile, err := os.Create("output.ndjson")
		if err != nil {
			fmt.Println("Error creating output file:", err)
			return
		}
		defer outputFile.Close()

		for result := range results {
			jsonResult, err := json.Marshal(result)
			if err != nil {
				fmt.Println("Error marshalling NDJSON:", err)
				continue
			}

			outputFile.Write(jsonResult)
			outputFile.WriteString("\n")
		}

		done <- struct{}{}
	}()

	<-done
	fmt.Println("Done formatting the data in output.ndjson")
}
