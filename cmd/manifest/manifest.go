package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"time"

	_ "github.com/mboye/kopi/loglevel"
	"github.com/mboye/kopi/manifest"
	log "github.com/sirupsen/logrus"
)

type manifestHeader struct {
	ID          string    `json:"ID"`
	Date        time.Time `json:"date"`
	Description string    `json:"description"`
}

var encoder = json.NewEncoder(os.Stdout)
var readFlags, writeFlags *flag.FlagSet

func main() {
	writeFlags = flag.NewFlagSet("write", flag.ExitOnError)
	encrypt := writeFlags.Bool("encrypt", false, "Encrypt manifest contents using AES-256")
	description := writeFlags.String("description", "", "Manifest description e.g. monthly backup 2019/1")

	readFlags = flag.NewFlagSet("read", flag.ExitOnError)
	decrypt := readFlags.Bool("decrypt", false, "Decrypt manifest contents using AES-256")

	if len(os.Args) < 3 {
		printCommandUsage()
		os.Exit(1)
	}

	subcommand := os.Args[1]
	backupPath := os.Args[2]
	var err error

	switch subcommand {
	case "read":
		if len(os.Args) < 4 {
			log.Error("Manifest ID parameter missing")
			printReadSubcommandUsage()
			os.Exit(1)
		}
		manifestID := os.Args[3]
		readFlags.Parse(os.Args[4:])
		reader, err := manifest.NewReader(backupPath, *decrypt, manifestID)
		if err != nil {
			log.Fatal(err)
		}

		err = reader.Execute()
	case "write":
		writeFlags.Parse(os.Args[3:])
		writer, err := manifest.NewWriter(backupPath, *encrypt, *description)
		if err != nil {
			log.Fatal(err)
		}

		err = writer.Execute()
	default:
		printCommandUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func printCommandUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s <read|write> <backup dir> [OPTIONS]\n", commandName)
}

func printWriteSubcommandUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s write <destination dir> [OPTIONS]\n\n", commandName)
	fmt.Fprintf(os.Stderr, "Pass index lines on STDIN.\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	writeFlags.PrintDefaults()
}

func printReadSubcommandUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s read <source dir> <manifest ID> [OPTIONS]\n", commandName)
	fmt.Fprintln(os.Stderr, "\nOptions:")
	readFlags.PrintDefaults()
}
