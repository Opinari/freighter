package archive

import (
	"log"
	"archive/tar"
	"io/ioutil"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"bufio"
)

type Archiver interface {
	Archive(inputDirPath string) (archivedFilePath string, err error)
	Unarchive(inputFilePath string, outputDirPath string) (unarchiveDirPath string, err error)
	IsSupported(fileExt string) bool
}

type TarArchiver struct {
}

func (a *TarArchiver) Archive(inputDirPath string) (archivedFilePath string, err error) {

	outputFilePath := inputDirPath + ".tar"

	log.Printf("Archiving files from: %s to: %s", inputDirPath, outputFilePath)

	// Create output file (tar)
	tarFile, err := os.Create(outputFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured opening output tar file: %s", err.Error())
	}

	// Create a new tar writer
	tarWriter := tar.NewWriter(tarFile)
	defer tarWriter.Close()

	// Recurse through the input root directory, adding each file / folder to the archive
	err = filepath.Walk(inputDirPath, func(filePath string, fileInfo os.FileInfo, err error) error {

		if err != nil {
			return fmt.Errorf("Error parsing input directory: %s", err.Error())
		}

		// Build file header
		fileHeader, err := tar.FileInfoHeader(fileInfo, "")
		if err != nil {
			return fmt.Errorf("Error reading file header: %s", err.Error())
		}

		fileHeader.Name, err = filepath.Rel(inputDirPath, filePath)
		if err != nil {
			return fmt.Errorf("Error determining file archive path: %s", err.Error())
		}

		// Write Header
		log.Printf("Writing header entry with name: %s, size: %d", fileHeader.Name, fileHeader.Size)
		if err := tarWriter.WriteHeader(fileHeader); err != nil {
			return fmt.Errorf("Error occured writing file header: %s", err.Error())
		}

		// Write file body if its a regular file, i.e not Dir / Symlink etc
		if fileInfo.Mode().IsRegular() {

			// Read file contents
			fileBody, err := ioutil.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("Error occured reading file body: %s", err.Error())
			}

			// Write Body
			if _, err := tarWriter.Write(fileBody); err != nil {
				return fmt.Errorf("Error occured writing file body: %s", err.Error())
			}
		}

		return nil
	})
	if err != nil {
		return "", fmt.Errorf("Error processing files to archive: %s", err.Error())
	}

	log.Printf("Files were archived to: '%s'", outputFilePath)

	archivedFilePath = outputFilePath
	return
}

func (a *TarArchiver) Unarchive(inputFilePath string, outputDirPath string) (unarchiveDirPath string, err error) {

	log.Printf("Unarchiving files from: %s to: %s", inputFilePath, outputDirPath)

	// Open archived file
	archivedFile, err := os.Open(inputFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured finding archived file: %s", err.Error())
	}
	defer archivedFile.Close()

	// Create the tar reader
	reader := bufio.NewReader(archivedFile)
	tarReader := tar.NewReader(reader)

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

		uncompressedFilePath := filepath.Join(outputDirPath, fileHeader.Name)

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

		case tar.TypeSymlink:
			log.Printf("Unarchiving symlink: %s \n", uncompressedFilePath)
			err := writeFileFromArchive(uncompressedFilePath, tarReader)
			if err != nil {
				return "", fmt.Errorf("Error occured whilst unarchiving symlink: %s", err.Error())
			}

		default:
			return "", fmt.Errorf("Unexpected File Type '%d' whilst unarchiving file: '%s'", fileHeader.Typeflag, fileHeader.Name)
		}
	}

	log.Printf("Unarchived files successfully to: %s", outputDirPath)

	unarchiveDirPath = outputDirPath
	return
}

func (a *TarArchiver) IsSupported(fileExt string) bool {
	if fileExt == ".tar" {
		return true
	}
	return false
}

func writeFileFromArchive(uncompressedFilePath string, tarReader io.Reader) (err error) {

	fileWriter, err := os.Create(uncompressedFilePath)
	if err != nil {
		return fmt.Errorf("Error creating file: %s", err.Error())
	}
	defer fileWriter.Close()

	if _, err := io.Copy(fileWriter, tarReader); err != nil {
		return fmt.Errorf("Error occured whilst writing to file: %s", err.Error())
	}

	return
}

func NewTarArchiver() Archiver {
	return &TarArchiver{}
}
