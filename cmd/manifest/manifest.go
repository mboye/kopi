package main

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"time"

	"github.com/dustin/go-humanize"
	"github.com/mboye/kopi/input"
	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/security"
	"github.com/mboye/kopi/util"
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
	util.SetLogLevel()

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
		err = readManifest(backupPath, *decrypt, manifestID)
	case "write":
		writeFlags.Parse(os.Args[3:])
		err = writeManifest(backupPath, *encrypt, *description)
	default:
		printCommandUsage()
		os.Exit(1)
	}

	if err != nil {
		log.Fatal(err)
	}
}

func readManifest(sourceDir string, decrypt bool, id string) error {
	manifestPath := fmt.Sprintf("%s/manifests/%s", sourceDir, id)
	manifestFile, err := os.Open(manifestPath)
	if err != nil {
		return fmt.Errorf("failed to open manifest: %s", err.Error())
	}
	defer manifestFile.Close()

	securityContext, err := security.NewContext(sourceDir, decrypt)
	if err != nil {
		return fmt.Errorf("failed to create security context: %s", err.Error())
	}

	encodedData, err := ioutil.ReadAll(manifestFile)
	if err != nil {
		return fmt.Errorf("failed to decompress manifest: %s", err.Error())
	}

	compressedData, err := securityContext.Decode(encodedData)
	if err != nil {
		return fmt.Errorf("failed to decode manifest: %s", err.Error())
	}

	decompressor, err := gzip.NewReader(bytes.NewBuffer(compressedData))
	if err != nil {
		return fmt.Errorf("failed to create decompressor: %s", err.Error())
	}
	defer decompressor.Close()

	decoder := json.NewDecoder(decompressor)
	decoder.DisallowUnknownFields()

	header := manifestHeader{}
	if err := decoder.Decode(&header); err != nil {
		return fmt.Errorf("failed to decode manifest header: %s", err.Error())
	}
	log.WithFields(log.Fields{
		"date":        header.Date,
		"id":          header.ID,
		"description": header.Description,
	}).Info("read manifest header")

	for decoder.More() {
		file := model.File{}
		if err := decoder.Decode(&file); err != nil {
			return fmt.Errorf("failed to decode file: %s", err.Error())
		}

		encoder.Encode(file)
	}

	return nil
}

func writeManifest(destinationDir string, encrypt bool, description string) error {
	now := time.Now().UTC()
	manifestFilename := fmt.Sprintf("%d/%02d/%02d/%d.manifest",
		now.Year(), now.Month(), now.Day(),
		now.Unix())

	header := manifestHeader{
		ID:          manifestFilename,
		Date:        now,
		Description: description}

	securityContext, err := security.NewContext(destinationDir, encrypt)
	if err != nil {
		log.WithField("error", err).Fatal("failed to create security context")
	}

	compressedManifest := bytes.NewBuffer(nil)
	compressor := gzip.NewWriter(compressedManifest)

	addToManifest := func(item interface{}) error {
		data, err := json.Marshal(item)
		if err != nil {
			return err
		}

		_, err = io.CopyN(compressor, bytes.NewBuffer(data), int64(len(data)))
		if err != nil {
			return err
		}

		_, err = compressor.Write([]byte{'\n'})
		return err
	}

	var fileCount, byteCount int64
	addFileToManifest := func(file *model.File) error {
		fileCount++
		byteCount += file.Size
		return addToManifest(file)
	}

	addToManifest(header)

	err = input.ProcessFilesWithProgress(addFileToManifest, 1)
	if err != nil {
		log.WithError(err).Error("failed to create compressed manifest")
		return err
	}
	compressor.Close()

	log.WithField("size", humanize.Bytes(uint64(compressedManifest.Len()))).Info("compressed manifest created")

	encodedManifest, err := securityContext.Encode(compressedManifest.Bytes())
	if err != nil {
		return err
	}
	log.WithField("size", len(encodedManifest)).Debug("encoded compressed manifest")

	outputPath := fmt.Sprintf("%s/manifests/%s", destinationDir, manifestFilename)
	if err := os.MkdirAll(filepath.Dir(outputPath), 0755); err != nil {
		log.Error("failed to create manifest directory")
		return err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0644)
	if err != nil {
		return err
	}
	defer outputFile.Close()
	log.WithField("path", outputFile.Name()).Debug("opened manifest file")

	bytesWritten, err := io.CopyN(outputFile, bytes.NewBuffer(encodedManifest), int64(len(encodedManifest)))
	if err != nil {
		return fmt.Errorf("failed to save manifest: %s", err.Error())
	}

	log.WithField("bytes_written", bytesWritten).Debug("wrote manifest")

	log.WithFields(log.Fields{
		"id":    header.ID,
		"files": fileCount,
		"bytes": humanize.Bytes(uint64(byteCount))}).Info("created manifest")

	return nil
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
