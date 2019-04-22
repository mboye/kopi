package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/scanner"
)

func main() {

	recursive := flag.Bool("recursive", true, "Index path recursively")
	initial := flag.Bool("init", false, "Initial index. Mark all files as modified.")
	withProgress := flag.Bool("progress", true, "Print indexing progress.")
	flag.Parse()

	rootPath := flag.Arg(0)
	s, err := scanner.New(rootPath, *recursive, *initial, *withProgress)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Execute(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Printf("Usage: %s <path>\n", commandName)
}
