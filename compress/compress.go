package compress

import (
	"compress/gzip"
	"log"
	"fmt"
	"bufio"
	"os"
	"io"
)

// TODO Fix this
const tmpUncompressedFileName = "/tmp.tar"

// TODO NO-OP
func CompressFile(inputDirPath string, outputFilePath string) (archivedFilePath string) {

	log.Printf("Compressing file from: %s to: %s", inputDirPath, outputFilePath)

	archivedFilePath = outputFilePath
	return
}

func UncompressFile(compressedFilePath string, outputDirPath string) (uncompressDirPath string, err error) {

	log.Printf("Uncompressing file from: %s to: %s", compressedFilePath, outputDirPath)

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
	uncompressedFilePath := outputDirPath + tmpUncompressedFileName
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

	// TODO Doing a double close here
	err = uncompressedFile.Close()
	if err != nil {
		return "", fmt.Errorf("Error closing uncompressed file: %s", err.Error())
	}

	log.Printf("Uncompressed file successfully to: %s", uncompressedFilePath)

	uncompressDirPath = uncompressedFilePath
	return nil
}
