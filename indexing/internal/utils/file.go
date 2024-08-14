package utils

import (
	"bufio"
	"fmt"
	"indexing/internal/models"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const maxBufferSize = 400 * 1024 // 400KB

func WalkFiles(root string, filePaths chan<- string) {

	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			fmt.Println("error walking through the directory:", err)
			return err
		}

		if !info.IsDir() {
			filePaths <- path
		}
		return nil
	})
	if err != nil {
		fmt.Println("error walking through the directory:", err)
	}

	close(filePaths)
}

func ProcessEmailFile(id int, wg *sync.WaitGroup, filePaths <-chan string, results chan<- models.Document) {
	defer wg.Done()

	for path := range filePaths {
		if isFileTooLarge(path) {
			fmt.Printf("worker %d: file %s exceeded buffer size of %d KB, skipping\n", id, path, maxBufferSize/1024)
			continue
		}

		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("worker %d error opening file: %v\n", id, err)
			continue
		}

		reader := bufio.NewReaderSize(file, maxBufferSize)

		var date, from, to, subject, fileName string
		var content strings.Builder
		var foundNewLine bool

		fileName = filepath.Base(path)
		if fileName == ".DS_Store" {
			file.Close()
			continue
		}

		for {
			line, err := reader.ReadString('\n')
			line = strings.TrimSpace(line)

			if err != nil {
				if err.Error() != "EOF" {
					fmt.Printf("worker %d error reading file: %v\n", id, err)
				}
				break
			}

			switch {
			case strings.HasPrefix(line, "Date:"):
				date = extractValue(line, "Date: ")
			case strings.HasPrefix(line, "From:") && !strings.HasPrefix(line, "X-From: "):
				from = extractValue(line, "From: ")
			case strings.HasPrefix(line, "To:") && !strings.HasPrefix(line, "X-To: "):
				to = extractValue(line, "To: ")
			case strings.HasPrefix(line, "Subject:"):
				subject = extractValue(line, "Subject: ")
			default:
				if foundNewLine {
					content.WriteString(line + "\n")
				} else if line == "" {
					foundNewLine = true
				}
			}
		}
		file.Close()

		indexNdLine := map[string]string{"_index": "emails"}
		dataNdLine := map[string]string{"path": path, "date": date, "from": from, "to": to, "subject": subject, "content": content.String()}
		doc := models.Document{Index: indexNdLine, Data: dataNdLine}

		results <- doc
	}
}

func isFileTooLarge(path string) bool {
	fileInfo, err := os.Stat(path)
	if err != nil {
		fmt.Printf("error getting file info: %v\n", err)
		return false
	}
	return fileInfo.Size() > maxBufferSize
}

func extractValue(line, prefix string) string {
	parts := strings.SplitN(line, prefix, 2)
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}
