package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/marcos-venicius/indexme/indexer"
)

func main() {
	folder := flag.String("i", "", "index a folder")
	search := flag.String("s", "", "search for a term")
	verbose := flag.Bool("v", false, "verbose mode")

	flag.Parse()

	idx, err := indexer.NewIndexer(*folder)

	if err != nil {
		fmt.Printf("Could not initialize the indexer: %s\n", err)
		return
	}

	if *verbose {
		idx.SetVerboseMode()
	}

	if *search != "" {
		result, err := idx.Search(*search, 10)

		if err != nil {
			fmt.Printf("Could not initialize the indexer: %s\n", err)
			return
		} else {
			for i, file := range result {
				fmt.Printf("%02d %f %s\n", i+1, file.Score, file.Path)
			}
		}
	} else if *folder != "" {
		idx.IgnoreFolderName(".git").IgnoreFolderName(".idea").IgnoreFolderName("venv")

		err = idx.Index()
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err)
	}
}
