package freighter

import (
	"log"
	"os"
	"github.com/Opinari/freighter/archive"
	"github.com/Opinari/freighter/compress"
	"github.com/Opinari/freighter/dropbox"
)

func restoreFile(restoreFilePath string, remoteFilePath string) {

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
	_, err = archive.UnarchiveFile(compressedFile, restoreFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup tmp Compressed and Archived files
	os.Remove(downloadedFilePath)
	os.Remove(compressedFile)
}

// TODO WIP
func backupFile(archiveDir string, backupFilePath string, remoteFilePath string) {

	// Archive File
	archivedFile := archive.ArchiveFile(archiveDir, backupFilePath + ".tar")
	log.Printf("Files were archived to: '%s'", archivedFile)

	// Compress File
	compressedFile := compress.CompressFile(archivedFile, archivedFile + ".gz")
	log.Printf("File was compressed to: '%s'", compressedFile)

	// Upload File
	uploadedFile, err := dropbox.UploadFile(backupFilePath, remoteFilePath)
	if err != nil {
		panic("Done For!")
	}
	log.Printf("File was uploaded to: '%s'", uploadedFile)
}

func ageRemoteFile(outputDir string, remoteFilePath string) {

}

func deleteRemoteFile(remoteFilePath string) {

}
