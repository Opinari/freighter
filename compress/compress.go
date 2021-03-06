package compress

import (
	"compress/gzip"
	"log"
	"fmt"
	"bufio"
	"os"
	"io"
	"path/filepath"
	"strings"
)

type Compressor interface {
	CompressFile(uncompressedFilePath string) (outputFilePath string, err error)
	UncompressFile(compressedFilePath string, uncompressedFileDir string) (outputFilePath string, err error)
	IsSupported(fileExt string) bool
}

type GzipCompressor struct {
}

func (c *GzipCompressor) CompressFile(uncompressedFilePath string) (outputFilePath string, err error) {

	compressedFilePath := uncompressedFilePath + ".gz"

	log.Printf("Compressing file from: %s to: %s", uncompressedFilePath, compressedFilePath)

	// Create compressed file
	compressedFile, err := os.Create(compressedFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured opening compressed file: %s", err.Error())
	}
	defer compressedFile.Close()

	// Open GzipWriter
	gzipWriter := gzip.NewWriter(compressedFile)
	defer gzipWriter.Close()

	// Open uncompressed file
	uncompressedFile, err := os.Open(uncompressedFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured finding uncompressed file: %s", err.Error())
	}
	defer uncompressedFile.Close()

	reader := bufio.NewReader(uncompressedFile)
	buffer := make([]byte, 4096)
	for {
		n, err := reader.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("Error occured during reading uncompressed file: %s", err.Error())
		}

		if n == 0 {
			break
		}

		_, err = gzipWriter.Write(buffer[:n])
		if err != nil {
			return "", fmt.Errorf("Error writing compressed file: %s", err.Error())
		}
	}

	err = gzipWriter.Flush()
	if err != nil {
		return "", fmt.Errorf("Error flushing to compressed file: %s", err.Error())
	}

	log.Printf("File was compressed to: '%s'", compressedFilePath)

	outputFilePath = compressedFilePath
	return
}

func (c *GzipCompressor) UncompressFile(compressedFilePath string, uncompressedFileDir string) (outputFilePath string, err error) {

	log.Printf("Uncompressing file from: %s to: %s", compressedFilePath, uncompressedFileDir)

	// Open compressed file
	compressedFile, err := os.Open(compressedFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured finding compressed file: %s", err.Error())
	}
	defer compressedFile.Close()

	// Open GzipReader
	gzipReader, err := gzip.NewReader(bufio.NewReader(compressedFile))
	if err != nil {
		return "", fmt.Errorf("Error occured during uncompression of file: %s", err.Error())

	}
	defer gzipReader.Close()

	// Open uncompressed file to write to
	_, compressedFileName := filepath.Split(compressedFilePath)
	uncompressedFilePath := uncompressedFileDir + "/"+ strings.TrimSuffix(compressedFileName, ".gz")

	uncompressedFile, err := os.Create(uncompressedFilePath)
	if err != nil {
		return "", fmt.Errorf("Error creating uncompressed file: %s", err.Error())
	}
	defer uncompressedFile.Close()

	writer := bufio.NewWriter(uncompressedFile)
	buffer := make([]byte, 4096)
	for {
		n, err := gzipReader.Read(buffer)
		if err != nil && err != io.EOF {
			return "", fmt.Errorf("Error occured during reading compressed file: %s", err.Error())
		}

		if n == 0 {
			break
		}

		_, err = writer.Write(buffer[:n])
		if err != nil {
			return "", fmt.Errorf("Error writing uncompressed file: %s", err.Error())
		}
	}

	err = writer.Flush()
	if err != nil {
		return "", fmt.Errorf("Error flushing to uncompressed file: %s", err.Error())

	}
	log.Printf("Uncompressed file successfully to: %s", uncompressedFilePath)

	return uncompressedFilePath, nil
}

func (c *GzipCompressor) IsSupported(fileExt string) bool {
	if fileExt == ".gz" {
		return true
	}
	return false
}

func NewGzipCompressor() Compressor {
	return &GzipCompressor{}
}
