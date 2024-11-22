package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marcos-venicius/indexme/indexer"
)

func main() {
	folder := flag.String("in", "", "index a folder")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	idx, err := indexer.NewIndexer(*folder)

	if err != nil {
		fmt.Printf("Could not initialize the indexer: %s\n", err)
		return
	}

	idx.IgnoreFolderName(".git").IgnoreFolderName(".idea").IgnoreFolderName("venv")

	if *verbose {
		idx.SetVerboseMode()
	}

	defer idx.Close()

	err = idx.Index()

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
