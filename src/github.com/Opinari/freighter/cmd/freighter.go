package main

import (
	"fmt"
	"github.com/Opinari/freighter/archive"
	"github.com/Opinari/freighter/compress"
	"github.com/Opinari/freighter/dropbox"
)


// TODO This will parse the command line args and decide what to operation to invoke
func main() {

	var input = "hello, world"
	fmt.Println(input)
}


// TODO The below should be moved into a different file, behind an interface if possible
func restoreFile(restoreFilePath string, remoteFilePath string) {

	//  Download File
	var downloadedFile = dropbox.DownloadFile(restoreFilePath, remoteFilePath);


	// Uncompress File
	var compressedFile = compress.UncompressFile(downloadedFile, restoreFilePath)
	fmt.Printf("Uncompressed file: %s", compressedFile)


	// Unarchive File
	var unarchiveDirPath = archive.UnarchiveFile(compressedFile, restoreFilePath)
	fmt.Printf("Unarchived archived was: " + unarchiveDirPath)

}

func backupFile(archiveDir string, backupFilePath string, remoteFilePath string) {

	// Archive File
	var archivedFile = archive.ArchiveFile(archiveDir, backupFilePath + ".tar")
	fmt.Printf("Files were archived to: %s", archivedFile)


	// Compress File
	var compressedFile = compress.CompressFile(archivedFile, archivedFile + ".gz")
	fmt.Printf("File was compressed to: %s", compressedFile)


	// Upload File
	var uploadedFile = dropbox.UploadFile(backupFilePath, remoteFilePath);
	fmt.Printf("File was uploaded to: %s", uploadedFile)

}

func ageRemoteFile(outputDir string, remoteFilePath string) {

}

func deleteRemoteFile(remoteFilePath string) {

}

