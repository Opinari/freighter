package client

import (
	"log"
	"os"
	"github.com/opinari/freighter/archive"
	"github.com/opinari/freighter/compress"
	"github.com/opinari/freighter/dropbox"
	"io/ioutil"
	"strconv"
)

func RestoreFile(restoreFilePath string, remoteFilePath string) {

	// Create Restore File Path
	os.MkdirAll(restoreFilePath, 0755)

	//  Download File
	downloadedFilePath, err := dropbox.DownloadFile(restoreFilePath, remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Uncompress File
	compressedFile, err := compress.UncompressFile(downloadedFilePath, restoreFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Unarchive Files
	_, err = archive.Unarchive(compressedFile, restoreFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup tmp Downloaded and Uncompressed files
	os.Remove(downloadedFilePath)
	os.Remove(compressedFile)
}

func BackupDirectory(backupFilePath string, remoteFilePath string) {

	// Archive Files
	archivedFilePath, err := archive.Archive(backupFilePath, backupFilePath + ".tar")
	if err != nil {
		log.Fatal(err)
	}

	// Compress File
	compressedFilePath, err := compress.CompressFile(archivedFilePath, archivedFilePath + ".gz")
	if err != nil {
		log.Fatal(err)
	}

	// Upload File
	_, err = dropbox.UploadFile(compressedFilePath, remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup tmp Archived and Compressed files
	os.Remove(archivedFilePath)
	os.Remove(compressedFilePath)
}

func AgeRemoteFile(outputFilePath string, remoteFilePath string) {

	// Perform age lookup
	age, err := dropbox.AgeFile(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create and open the output file for writing
	outputFile, err := os.Create(outputFilePath)
	defer outputFile.Close()
	if err != nil {
		log.Fatal(err)
	}

	// Write out age value to file
	ageBytes := []byte(strconv.Itoa(age))
	if err := ioutil.WriteFile(outputFilePath, ageBytes, 0666); err != nil {
		log.Fatal(err)
	}

}

func DeleteRemoteFile(remoteFilePath string) {

	_, err := dropbox.DeleteFile(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

}
