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

func ProcessEmailFile(id int, wg *sync.WaitGroup, filePaths <-chan string, results chan<- models.Document) {
	defer wg.Done()

	for path := range filePaths {
		file, err := os.Open(path)
		if err != nil {
			fmt.Printf("worker %d error opening file: %v\n", id, err)
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
		doc := models.Document{Index: indexNdLine, Data: dataNdLine}

		results <- doc

		file.Close()
	}
}
