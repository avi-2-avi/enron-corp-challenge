package worker

import (
	"bufio"
	"fmt"
	"indexing/models"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

const maxBufferSize = 400 * 1024 // 400KB

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

		reader := bufio.NewReader(file)

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
				break
			}

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
				content.WriteString(line + "\n")
			} else {
				if line == "" {
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
		return true
	}
	return fileInfo.Size() > maxBufferSize
}
