package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/restorer"

	_ "github.com/mboye/kopi/loglevel"
	log "github.com/sirupsen/logrus"
)

func main() {
	dryRun := flag.Bool("dry-run", false, "Dry run. Only verify that index is restorable.")
	decrypt := flag.Bool("decrypt", false, "Decrypt blocks using AES-256 while restoring")
	progressInterval := flag.Int("progress", 10, "Progres printing interval in seconds. An interval of zero disables printing.")
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() != 2 {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	if *dryRun {
		log.Info("Dry run mode enabled")
	}

	inputDir := flag.Arg(0)
	outputDir := flag.Arg(1)
	restorer, err := restorer.New(inputDir, outputDir, *dryRun, *decrypt, *progressInterval)

	err = restorer.Execute()
	if err != nil {
		log.Fatal(err)
	}
}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <store dir> <destination dir>\n\n", commandName)
	fmt.Fprintf(os.Stderr, "Pass index lines on STDIN.\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}
