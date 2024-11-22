package indexer

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"unicode/utf8"
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

func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	defer file.Close()

	maxSize := 512
	bytes := make([]byte, 0, maxSize)

	reader := bufio.NewReader(file)

	for {
		if maxSize <= 0 {
			break
		}

		c, err := reader.ReadByte()

		if err != nil {
			if err == io.EOF {
				break
			}

			return false, err
		}

		if c == '\n' {
			break
		}

		bytes = append(bytes, c)

		maxSize--
	}

	return !utf8.Valid(bytes), nil
}
