package archive

import (
	_ "archive/tar"
	"log"
	"archive/tar"
	"io/ioutil"
	"fmt"
	"io"
	_ "os"
	"bytes"
	"os"
)

// TODO NO-OP
func ArchiveFile(inputDirPath string, outputFilePath string) (archivedFilePath string) {

	log.Printf("Archiving files from: %s to: %s", inputDirPath, outputFilePath)

	archivedFilePath = outputFilePath
	return
}

func UnarchiveFile(inputFilePath string, outputDirPath string) (unarchiveDirPath string, err error) {

	log.Printf("Unarchiving files from: %s to: %s", inputFilePath, outputDirPath)


	// Read in archived file from path
	archivedFileBytes, err := ioutil.ReadFile(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured finding archived file: %s", err.Error())
	}


	// Go from bytes > Reader
	// FIXME Do we really need this?
	bytesReader := bytes.NewReader(archivedFileBytes)


	// Create the tar reader
	tarReader := tar.NewReader(bytesReader)


	// Iterate through the files in the archive.
	for {
		fileHeader, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		// Write out the damn file or folder
		uncompressedFilePath := outputDirPath + "/" + fileHeader.Name


		switch fileHeader.Typeflag {

		case tar.TypeDir:
			log.Printf("Folder name %s \n", uncompressedFilePath)
			os.MkdirAll(uncompressedFilePath, 666)

		case tar.TypeReg:
			log.Printf("File name %s \n", uncompressedFilePath)
			fileWriter, _ := os.Create(uncompressedFilePath)
			if _, err := io.Copy(fileWriter, tarReader); err != nil {
				return "", fmt.Errorf("Error occured whilst unarchiving files: %s", err.Error())
			}
		}
	}

	log.Printf("Unarchived files succesfully to: %s", outputDirPath)

	unarchiveDirPath = outputDirPath
	return
}
