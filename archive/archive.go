package archive

import (
	"log"
	"archive/tar"
	"io/ioutil"
	"fmt"
	"io"
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

	// Create the tar reader
	tarReader := tar.NewReader(bytes.NewReader(archivedFileBytes))

	// Iterate through the files in the archive
	for {
		fileHeader, err := tarReader.Next()
		if err == io.EOF {
			// end of tar archive
			break
		}
		if err != nil {
			log.Fatalln(err)
		}

		uncompressedFilePath := outputDirPath + "/" + fileHeader.Name

		switch fileHeader.Typeflag {

		case tar.TypeDir:
			log.Printf("Unarchiving folder: %s \n", uncompressedFilePath)
			err := os.MkdirAll(uncompressedFilePath, 0755)
			if err != nil {
				return "", fmt.Errorf("Error occured whilst unarchiving folder: %s", err.Error())
			}

		case tar.TypeReg:
			log.Printf("Unarchiving file: %s \n", uncompressedFilePath)
			err := writeFileFromArchive(uncompressedFilePath, tarReader)
			if err != nil {
				return "", fmt.Errorf("Error occured whilst unarchiving file: %s", err.Error())
			}

		default:
			return "", fmt.Errorf("Unexpected File Type '%s' whilst unarchiving file: '%s'", fileHeader.Typeflag, fileHeader.Name)
		}
	}

	log.Printf("Unarchived files successfully to: %s", outputDirPath)

	unarchiveDirPath = outputDirPath
	return
}

func writeFileFromArchive(uncompressedFilePath string, tarReader io.Reader) (err error) {

	fileWriter, err := os.Create(uncompressedFilePath)
	defer fileWriter.Close()
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err.Error())
	}

	if _, err := io.Copy(fileWriter, tarReader); err != nil {
		return fmt.Errorf("Error occured whilst writing to file: %s", err.Error())
	}

	return
}