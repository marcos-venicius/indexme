package indexer

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
)

func absolutePath(folderPath string) (string, error) {
	abs, err := filepath.Abs(folderPath)

	if err != nil {
		return "", err
	}

	stat, err := os.Stat(abs)

	if err != nil {
		return "", err
	}

	if !stat.IsDir() {
		return "", errors.New(fmt.Sprintf("\"%s\" is not a valid directory", folderPath))
	}

	return abs, nil
}
