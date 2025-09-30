package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

func CreateFolders(basePath string, folders []string) error {
	for _, folder := range folders {
		path := filepath.Join(basePath, folder)

		if err := os.MkdirAll(path, 0755); err != nil {
			fmt.Printf("❌ Failed to create folder %q: %v\n", path, err)
		} else {
			fmt.Printf("✅ Ensured folder %q exists\n", path)
		}
	}
	return nil
}
