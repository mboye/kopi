package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mboye/kopi/model"
	"github.com/mboye/kopi/util"
	log "github.com/sirupsen/logrus"
)

type fileHandlerFunc func(file *model.File) error

var encoder = json.NewEncoder(os.Stdout)

func main() {
	util.SetLogLevel()
	maxBlockSize := flag.Int64("maxBlockSize", 1024*1024*10, "Split files into blocks of this size")
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() != 1 {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	outputDir := flag.Arg(0)
	log.WithField("destination", outputDir).Info("Beginning to store files")

	filterAndStoreFile := func(file *model.File) error {
		if file.Mode.IsDir() {
			file.Modified = false
			encoder.Encode(file)
			return nil
		}

		if !file.Modified {
			log.WithField("path", file.Path).Debug("Skipping unmodified file")
			return nil
		}

		file.Modified = false
		return storeFile(file, outputDir, *maxBlockSize)
	}

	if err := forEachFileOnStdin(filterAndStoreFile); err != nil {
		log.Fatal(err)
	}
}

func storeFile(file *model.File, outputDir string, maxBlockSize int64) error {
	if err := refreshFileMetadata(file); err != nil {
		log.Fatal(err)
	}

	logger := log.WithField("path", file.Path)

	if inputFile, err := os.Open(file.Path); err != nil {
		return err
	} else {
		defer inputFile.Close()
		fileReader := io.LimitReader(inputFile, file.Size)

		logger.WithField("max_offset", file.Size).Debug("Read file")

		fileOffset := int64(0)
		for fileOffset < file.Size {
			blockReader := io.LimitReader(fileReader, maxBlockSize)
			blockOffset := fileOffset

			hasher := md5.New()
			blockBuffer := make([]byte, maxBlockSize)

			bytesRead, err := blockReader.Read(blockBuffer)
			if err != nil {
				return err
			}

			logger.WithField("bytes_read", bytesRead).Debug("Read file")

			fileOffset += int64(bytesRead)

			if int64(bytesRead) != maxBlockSize && fileOffset != file.Size {
				return errors.New("Incomplete buffer read")
			}

			blockSize := bytesRead
			hasher.Write(blockBuffer[:blockSize])
			hash := fmt.Sprintf("%x", hasher.Sum(nil))
			block := model.Block{Hash: hash, Offset: blockOffset, Size: int64(blockSize)}

			outputPath := fmt.Sprintf("%s/%s/%s.block", outputDir, hash[:2], hash)
			_, err = os.Stat(outputPath)
			if err == nil {
				logger.WithFields(log.Fields{"hash": hash, "offset": blockOffset, "size": blockSize}).Debug("Reusing existing block")
				file.AddBlock(block)
				continue
			}

			blockDirPath := fmt.Sprintf("%s/%s", outputDir, hash[:2])
			err = os.MkdirAll(blockDirPath, 0755)
			if err != nil {
				return err
			}

			outputFile, err := os.Create(outputPath)
			if err != nil {
				return err
			}
			logger.WithField("output_path", outputPath).Debug("Created output file")

			bytesWritten, err := outputFile.Write(blockBuffer[:blockSize])
			if err != nil {
				outputFile.Close()
				return err
			}

			if bytesWritten != blockSize {
				return errors.New("Incomplete block write")
			}

			outputFile.Close()
			file.AddBlock(block)
			logger.WithFields(log.Fields{"hash": hash, "offset": blockOffset, "size": blockSize}).Debug("Block created")

			logger.WithField("fileOffset", fileOffset).Debug("File offset")
		}

		logger.Debug("File read completed")

		if err := encoder.Encode(file); err != nil {
			return err
		}
	}

	return nil
}

func refreshFileMetadata(file *model.File) error {
	if fileInfo, err := os.Stat(file.Path); err != nil {
		log.WithError(err).Error("Failed to get file metadata")
		return err
	} else {
		if file.Size != fileInfo.Size() {
			log.WithFields(log.Fields{"path": file.Path, "expected_size": file.Size, "actual_size": fileInfo.Size()}).Warn("File size has changed")
		}
		file.Size = fileInfo.Size()

		if file.ModifiedTime.UTC() != fileInfo.ModTime().UTC() {
			log.WithFields(log.Fields{"path": file.Path, "expected_modtime": file.ModifiedTime, "actual_modtime": fileInfo.ModTime()}).Warn("File modification time has changed")
		}
		file.ModifiedTime = fileInfo.ModTime()

		if file.Mode != fileInfo.Mode() {
			log.WithFields(log.Fields{"path": file.Path, "expected_mode": file.Mode, "actual_mode": fileInfo.Mode()}).Warn("File mode has changed")
		}
		file.Mode = fileInfo.Mode()
	}

	return nil
}

func forEachFileOnStdin(handler fileHandlerFunc) error {
	decoder := json.NewDecoder(os.Stdin)

	for decoder.More() {
		file := &model.File{}
		if err := decoder.Decode(file); err != nil {
			return fmt.Errorf("Failed to decode file: %s", err.Error())
		}

		if err := handler(file); err != nil {
			return err
		}

	}
	return nil
}

func printUsage() {
	commandName := filepath.Base(os.Args[0])
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <destination dir>\n\n", commandName)
	fmt.Fprintf(os.Stderr, "Pass index lines on STDIN.\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}