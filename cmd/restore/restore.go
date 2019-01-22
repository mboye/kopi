package main

import (
	"crypto/md5"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
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
	flag.Usage = printUsage
	flag.Parse()

	if flag.NArg() != 2 {
		log.Error("Path argument missing")
		printUsage()
		os.Exit(1)
	}

	inputDir := flag.Arg(0)
	outputDir := flag.Arg(1)
	log.WithField("destination", outputDir).Info("Beginning to store files")

	restoreFile := func(file *model.File) (err error) {
		if file.Mode.IsDir() {
			return restoreDir(file, outputDir)
		} else {
			return restoreFile(file, inputDir, outputDir)
		}
	}

	if err := forEachFileOnStdin(restoreFile); err != nil {
		log.Fatal(err)
	}
}

func restoreDir(file *model.File, outputDir string) error {
	outputPath := fmt.Sprintf("%s/%s", outputDir, file.Path)
	log.WithFields(log.Fields{
		"path": outputPath,
		"mode": file.Mode}).Debug("restoring directory")
	return os.MkdirAll(outputPath, file.Mode)
}

func restoreFile(file *model.File, inputDir, outputDir string) error {
	log.WithFields(log.Fields{
		"path": file.Path,
		"mode": file.Mode}).Debug("restoring file")

	if file.Blocks == nil || len(file.Blocks) == 0 {
		return errors.New("cannot restore file without blocks")
	}

	outputPath := fmt.Sprintf("%s/%s", outputDir, file.Path)
	parentDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(parentDir, 0755); err != nil {
		log.WithError(err).Error("failed to create parent directory of file")
		return err
	}

	outputFile, err := os.OpenFile(outputPath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, file.Mode)
	if err != nil {
		return err
	}
	defer outputFile.Close()

	progress := 0
	for _, block := range file.Blocks {
		progress++
		log.WithFields(
			log.Fields{
				"progress":     progress,
				"max_progress": len(file.Blocks)}).Info("restoring block")

		restoreBlock := func() error {
			blockPath := fmt.Sprintf("%s/%s/%s.block", inputDir, block.Hash[:2], block.Hash)
			blockFile, err := os.Open(blockPath)
			if err != nil {
				log.WithError(err).Error("failed to open block file")
				return err
			}
			defer blockFile.Close()

			blockReader := io.LimitReader(blockFile, block.Size)
			blockData, err := ioutil.ReadAll(blockReader)
			if err != nil {
				log.WithError(err).Error("failed to read block data")
				return err
			}
			blockSize := len(blockData)

			if blockSize != int(block.Size) {
				log.WithFields(
					log.Fields{
						"actual_size":   blockSize,
						"expected_size": block.Size,
						"block_path":    blockPath}).Error("corrupt block detected")

				return fmt.Errorf("corrupt block: %s", blockPath)
			}

			hasher := md5.New()
			_, err = hasher.Write(blockData)
			if err != nil {
				return err
			}
			hash := fmt.Sprintf("%x", hasher.Sum(nil))

			if hash != block.Hash {
				log.WithFields(
					log.Fields{
						"actual_hash":   hash,
						"expected_hash": block.Hash}).Error("corrupt block detected")

				return fmt.Errorf("corrupt block: %s", blockPath)
			}

			if _, err = outputFile.Write(blockData); err != nil {
				log.WithError(err).Error("failed to restore block")
				return err
			}

			return nil
		}

		if err := restoreBlock(); err != nil {
			return err
		}
	}

	log.WithFields(
		log.Fields{
			"path":          file.Path,
			"restored_path": outputPath}).Info("file restored")
	return nil
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
	fmt.Fprintf(os.Stderr, "Usage: %s [OPTIONS] <store dir> <destination dir>\n\n", commandName)
	fmt.Fprintf(os.Stderr, "Pass index lines on STDIN.\n")
	fmt.Fprintln(os.Stderr, "\nOptions:")
	flag.PrintDefaults()
}