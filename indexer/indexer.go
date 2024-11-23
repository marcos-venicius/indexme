package indexer

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

type Indexer struct {
	verboseOutput        bool
	baseDirectory        string
	ignoredFileNames     map[string]struct{}
	ignoredFolderNames   map[string]struct{}
	ignoredFiles         int
	ignoredFolders       int
	indexedFiles         int
	indexFoldersSync     sync.WaitGroup
	abspath              string
	db                   *sql.DB
	folderTermsFrequency map[string]int
}

func NewIndexer(baseDirectory string) (*Indexer, error) {
	// TODO: check if the baseDirectory has never indexed before
	// TODO: add the baseDirectory to the database

	homedir, err := os.UserHomeDir()

	if err != nil {
		panic(err)
	}

	databaseLocation := filepath.Join(homedir, "indexme.db")

	db, err := sql.Open("sqlite3", databaseLocation)

	if err != nil {
		return nil, err
	}

	indexer := &Indexer{
		verboseOutput:        false,
		baseDirectory:        baseDirectory,
		ignoredFileNames:     make(map[string]struct{}, 0),
		ignoredFolderNames:   make(map[string]struct{}, 0),
		ignoredFiles:         0,
		ignoredFolders:       0,
		indexedFiles:         0,
		indexFoldersSync:     sync.WaitGroup{},
		abspath:              "",
		db:                   db,
		folderTermsFrequency: make(map[string]int),
	}

	err = indexer.DbSetup()

	if err != nil {
		defer db.Close()

		return nil, err
	}

	return indexer, nil
}

func (i *Indexer) addTokenToFolderTermsFrequency(token string, freq int) {
	if count, ok := i.folderTermsFrequency[token]; ok {
		i.folderTermsFrequency[token] += count + freq
	} else {
		i.folderTermsFrequency[token] = freq
	}
}

func (i *Indexer) SetVerboseMode() *Indexer {
	i.verboseOutput = true

	return i
}

func (i *Indexer) IgnoreFolderName(folderName string) *Indexer {
	i.ignoredFolderNames[folderName] = struct{}{}

	return i
}

func (i *Indexer) IgnoreFileName(fileName string) *Indexer {
	i.ignoredFileNames[fileName] = struct{}{}

	return i
}

func (i *Indexer) Index() error {
	defer i.db.Close()

	abs, err := absolutePath(i.baseDirectory)

	if err != nil {
		return err
	}

  if folder := i.GetFolderByAbsPath(abs); folder != nil {
    return errors.New("You already indexed this folder")
  }

	i.abspath = abs

	base := filepath.Base(abs)

	if err := i.CreateFolder(base, abs); err != nil {
		return nil
	}

	i.indexFoldersSync.Add(1)

	go i.indexFolder(abs)

	i.indexFoldersSync.Wait()

	folder := i.GetFolderByAbsPath(i.abspath)

	for token, freq := range i.folderTermsFrequency {
		i.AddFolderTermsFrequency(folder.id, token, freq)
	}

	if i.verboseOutput {
		fmt.Println()
	}

	fmt.Printf("%d files ignored\n", i.ignoredFiles)
	fmt.Printf("%d folders ignored\n", i.ignoredFolders)
	fmt.Printf("%d indexed files\n", i.indexedFiles)

	return nil
}
