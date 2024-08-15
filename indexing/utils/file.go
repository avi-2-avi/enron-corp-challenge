package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

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
