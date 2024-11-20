package main

// TODO: save the indexed documents to a sqlite database
// TODO: allow the user to search for a term using the CLI

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
	"unicode"
	"unicode/utf8"
)

type Tokens struct {
	tokens map[string]int // map[term]frequency
}

type IndexDb struct {
	globalTermFrequency   map[string]int            // map[term]frequency
	documentTermFrequency map[string]map[string]int // map[document(filepath)]map[term]frequency
}

type Tree struct {
	verbose        bool
	root           string
	filesCount     int
	ignoredFolders []string
	ignoredFiles   []string
	// used to sync the folders.
	// each sub folder has your own goroutine
	sync sync.WaitGroup
	// this is used to make the program process one file by time,
	// but later it will be improved to have just only one file in memory at time and every time a file is indexed
	// the database will be updated and the content from the current file will be discarded from the memory
	// cause today we have all the indexed tokens in memory, even though we work with one file at time, at the end all the N files tokens will be in the memory
	filesContentMu sync.Mutex
	// this is a temporary structure.
	// it'll be replaced by a sqlite db
	indexDb IndexDb
}

func newTokens() *Tokens {
	return &Tokens{
		tokens: make(map[string]int),
	}
}

func newTree(root string, verbose bool) *Tree {
	abs, err := filepath.Abs(root)

	if err != nil {
		perror(err.Error())
	}

	return &Tree{
		root:           abs,
		verbose:        verbose,
		filesCount:     0,
		ignoredFolders: []string{".idea", ".git"},
		ignoredFiles:   []string{},
		sync:           sync.WaitGroup{},
		filesContentMu: sync.Mutex{},
		indexDb: IndexDb{
			globalTermFrequency:   make(map[string]int),
			documentTermFrequency: make(map[string]map[string]int),
		},
	}
}

func (t *Tokens) add(token string) {
	if count, ok := t.tokens[token]; ok {
		t.tokens[token] = count + 1
	} else {
		t.tokens[token] = 1
	}
}

func (t *Tree) isFolderIgnored(name string) bool {
	for _, ig := range t.ignoredFolders {
		if ig == name {
			return true
		}
	}

	return false
}

func isAlphaNumeric(b byte) bool {
	return unicode.IsLetter(rune(b)) || unicode.IsDigit(rune(b)) || rune(b) == '_'
}

func isWhitespace(b byte) bool {
	return rune(b) == '\t' || rune(b) == '\n' || rune(b) == ' ' || rune(b) == '\r'
}

// tokenize Tokenizes the content and return the frequency of each token
func (t *Tree) tokenize(content []byte) map[string]int {
	tokens := newTokens()

	start := 0
	cursor := 0
	size := len(content)

	for cursor < size-1 {
		for isWhitespace(content[cursor]) && cursor < size-1 {
			cursor++
		}

		start = cursor

		if isAlphaNumeric(content[cursor]) {
			for isAlphaNumeric(content[cursor]) && cursor < size-1 {
				cursor++
			}

			token := string(content[start:cursor])

			tokens.add(token)
			continue
		}

		c := content[cursor]

		cursor++

		if cursor < size && content[cursor] == c {
			for c == content[cursor] && cursor < size-1 {
				cursor++
			}

			tokens.add(string(content[start:cursor]))

		} else {
			tokens.add(string(content[start:cursor]))
		}
	}

	return tokens.tokens
}

// this function checks if the first line has valid utf8 chars, if not, it's considered a binary  file
func isBinaryFile(path string) (bool, error) {
	file, err := os.Open(path)

	if err != nil {
		return false, err
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	scanner.Scan()

	firstLine := scanner.Text()

	return !utf8.ValidString(firstLine), nil
}

func (t *Tree) readFile(parent string, file os.DirEntry) {
	name := file.Name()
	path := filepath.Join(parent, name)

	info, err := file.Info()

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file \"%s\": %s\n", path, err.Error())
		return
	}

	size := info.Size()

	if t.verbose {
		fmt.Printf("Reading %d bytes from \"%s\"\n", size, path)
	}

	content, err := os.ReadFile(path)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read file \"%s\": %s\n", path, err.Error())
		return
	}

	t.filesContentMu.Lock()
	defer t.filesContentMu.Unlock()

	if t.verbose {
		fmt.Printf("Preparing \"%s\"\n", path)
	}

	tokens := t.tokenize(content)

	t.indexDocument(path, tokens)

	t.filesCount++
}

func (t *Tree) indexDocument(path string, termFrequency map[string]int) {
	if t.verbose {
		fmt.Printf("Indexing \"%s\"\n", path)
	}

	t.indexDb.documentTermFrequency[path] = termFrequency

	for k, v := range termFrequency {
		if count, ok := t.indexDb.globalTermFrequency[k]; ok {
			t.indexDb.globalTermFrequency[k] = count + v
		} else {
			t.indexDb.globalTermFrequency[k] = v
		}
	}

	if t.verbose {
		fmt.Printf("%d tokens indexed in \"%s\"\n", len(termFrequency), path)
	}
}

func (t *Tree) searchFiles(base string) {
	defer t.sync.Done()

	if t.verbose {
		fmt.Printf("Looking up at \"%s\"\n", base)
	}

	entries, err := os.ReadDir(base)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Could not read \"%s\"\n", base)
		return
	}

	for _, entry := range entries {
		name := entry.Name()
		path := filepath.Join(base, name)

		if entry.IsDir() {

			if t.isFolderIgnored(name) {
				if t.verbose {
					fmt.Printf("Ignoring \"%s\" directory\n", path)
				}
				continue
			}

			t.sync.Add(1)
			go t.searchFiles(filepath.Join(base, name))
		} else {
			isBinary, err := isBinaryFile(path)

			if err != nil {
				fmt.Printf("Ignoring \"%s\" file due to: %s\n", path, err.Error())
				continue
			}

			if isBinary {
				if t.verbose {
					fmt.Printf("Ignoring binary file \"%s\"\n", path)
				}
			} else {
				t.readFile(base, entry)
			}
		}
	}
}

func (t *Tree) readFiles() {
	t.sync.Add(1)
	go t.searchFiles(t.root)

	t.sync.Wait()

	if t.verbose {
		fmt.Printf("\n\n")
	}

	if t.filesCount > 0 {
		fmt.Printf("%d files indexed successfully\n", t.filesCount)
		fmt.Printf("%d tokens indexed\n", len(t.indexDb.globalTermFrequency))
	} else {
		fmt.Println("no files indexed")
	}
}

func usage(exitcode int) {
	flag.Usage()

	os.Exit(exitcode)
}

func perror(format string, a ...any) {
	fmt.Fprintf(os.Stderr, format+"\n", a...)

	os.Exit(1)
}

func main() {
	folder := flag.String("in", "", "-in . # this will make the tool index the current directory")
	verbose := flag.Bool("v", false, "verbose output")

	flag.Parse()

	if folder == nil {
		usage(1)
	}

	if *folder == "" {
		perror("please, provide a folder path")
	}

	state, err := os.Stat(*folder)

	if err != nil {
		perror(err.Error())
	}

	if !state.IsDir() {
		perror("\"%s\" is not a directory", *folder)
	}

	t := newTree(*folder, *verbose)

	start := time.Now()

	t.readFiles()

	elapsed := time.Since(start)

	if t.verbose {
		fmt.Printf("Time taken: %s\n", elapsed)
	}
}
