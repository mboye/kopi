package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	_ "github.com/mboye/kopi/loglevel"
	"github.com/mboye/kopi/storer"
	log "github.com/sirupsen/logrus"
)

var encoder = json.NewEncoder(os.Stdout)

func main() {
	maxBlockSize := flag.Int64("maxBlockSize", 1024*1024*10, "Split files into blocks of this size")
	encrypt := flag.Bool("encrypt", false, "Encrypt stored blocks using AES-256")
	progressInterval := flag.Uint("progress", 10, "Progres printing interval in seconds. An interval of zero disables printing.")
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() != 1 {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	outputDir := flag.Arg(0)

	s, err := storer.New(outputDir, *maxBlockSize, *encrypt, *progressInterval)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Execute(); err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <destination dir>\n\n", commandName)
	fmt.Fprintf(os.Stderr, "Pass index lines on STDIN.\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}
