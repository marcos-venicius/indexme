package indexer

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/marcos-venicius/indexme/tokenizer"
)

func (i *Indexer) indexFolder(absBase string) {
	defer i.indexFoldersSync.Done()

	baseFolderName := filepath.Base(absBase)

	if i.isIgnoredFolder(baseFolderName) {
		i.ignoredFolders++

		if i.verboseOutput {
			fmt.Printf("Ignoring folder \"%s\"\n", absBase)
		}

		return
	}

	if i.verboseOutput {
		fmt.Printf("Looking up directory \"%s\"\n", absBase)
	}

	entries, err := os.ReadDir(absBase)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read directory \"%s\": %s\n", absBase, err.Error())

		return
	}

	for _, entry := range entries {
		name := entry.Name()
		abs := filepath.Join(absBase, name)

		if entry.IsDir() {
			i.indexFoldersSync.Add(1)

			go i.indexFolder(abs)
		} else {
			i.indexFile(abs)
		}
	}
}

func (i *Indexer) indexFile(abspath string) {
	baseFileName := filepath.Base(abspath)

	if i.isIgnoredFile(baseFileName) {
		i.ignoredFiles++

		if i.verboseOutput {
			fmt.Printf("Ignoring file \"%s\"\n", abspath)
		}

		return
	}

	if i.verboseOutput {
		fmt.Printf("Indexing file \"%s\"\n", abspath)
	}

	file, err := os.OpenFile(abspath, os.O_RDONLY, 0600)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file \"%s\": %s\n", abspath, err)

		return
	}

	tokens := tokenizer.Tokenize(file)

  // TODO: get tokens frequency to this document
  // TODO: update global term frequency for this base path
  // TODO: insert this tokens frequency in the database

	i.indexedFiles++

	if i.verboseOutput {
		fmt.Printf("Indexed %d tokens in \"%s\"\n", len(tokens), abspath)
	}
}

func (i *Indexer) isIgnoredFile(filename string) bool {
	_, ok := i.ignoredFileNames[filename]

	return ok
}

func (i *Indexer) isIgnoredFolder(folder string) bool {
	_, ok := i.ignoredFolderNames[folder]

	return ok
}
