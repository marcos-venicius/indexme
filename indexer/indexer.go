package indexer

import (
	"bufio"
	"bytes"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"sync"

	"github.com/marcos-venicius/indexme/tokenizer"
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

type IndexedDocument struct {
	Path  string
	Score float32
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

	abs, err := absolutePath(baseDirectory)

	if err != nil {
		defer db.Close()

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
		abspath:              abs,
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

	if folder := i.GetFolderByAbsPath(i.abspath); folder != nil {
		return errors.New("You already indexed this folder")
	}

	base := filepath.Base(i.abspath)

	if err := i.CreateFolder(base, i.abspath); err != nil {
		return nil
	}

	i.indexFoldersSync.Add(1)

	go i.indexFolder(i.abspath)

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

func (i *Indexer) Search(term string, top int) ([]IndexedDocument, error) {
	var folder *FolderTable

	if folder = i.GetFolderByAbsPath(i.abspath); folder == nil {
		return nil, errors.New("Please, before search index this folder!")
	}

	bytesReader := bytes.NewReader([]byte(term))
	reader := bufio.NewReader(bytesReader)

	tokens := tokenizer.Tokenize(reader)

	results := make([]IndexedDocument, 0)

	folderTermsFrequency := i.GetFolderTermsFrequency(i.abspath)
	documentTermsFrequency := i.GetDocumentTermsFrequency(folder.id)

	for document, docFreq := range documentTermsFrequency {
		var score float32 = 0

		for _, token := range tokens {
			if freq, ok := docFreq[token]; ok {
				score += float32(freq) / float32(folderTermsFrequency[token])
			}
		}

		if score <= 0 {
			continue
		}

		results = append(results, IndexedDocument{
			Path:  document,
			Score: score,
		})
	}

	sort.Slice(results, func(i, j int) bool {
		return results[i].Score > results[j].Score
	})

	return results[:min(top, len(results))], nil
}
