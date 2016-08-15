package compress

import (
	"compress/gzip"
	"log"
	"io/ioutil"
	"fmt"
	"bytes"
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

	// Read in compressed file from path
	compressedFile, err := ioutil.ReadFile(compressedFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured finding compressed file: %s", err.Error())
	}


	// Go from bytes > Reader
	// FIXME Do we really need this?
	bytesReader := bytes.NewReader(compressedFile)


	// Uncompress
	gzipReader, err := gzip.NewReader(bytesReader)
	gzipReader.Close()


	// Go from Reader > bytes
	// FIXME Do we really need this?
	buf := new(bytes.Buffer)
	buf.ReadFrom(gzipReader)
	uncompressedFileBytes := buf.Bytes()


	// TODO defer file close?
	// TODO Parse filename
	// Write response to file
	uncompressedFilePath := outputDirPath + tmpUncompressedFileName
	err = ioutil.WriteFile(uncompressedFilePath, uncompressedFileBytes, 0644)
	if err != nil {
		return "", fmt.Errorf("Error occured during uncompression of file: %s", err.Error())
	}


	log.Printf("Uncompressed file succesfully to: %s", uncompressedFilePath)

	uncompressDirPath = uncompressedFilePath
	return
}
