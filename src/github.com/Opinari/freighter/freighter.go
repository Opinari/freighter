package main

import (
	"flag"
	"fmt"
	"github.com/Opinari/freighter/archive"
	"github.com/Opinari/freighter/compress"
	"github.com/Opinari/freighter/dropbox"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func main() {

	showBanner()
	runCLI();
}

// TODO Make the CLI an instance and execute the args / flags upon it
func runCLI() {

	// Sub cmd parsing
	if len(os.Args) < 2 {
		log.Println("freighter: Provide at least one argument for operation")
		log.Fatalln("See 'freighter --help' for more info \n")
	}
	operation := os.Args[1]


	// Flag parsing
	// Declare placeholder vars for opts
	var remoteFilePath, restoreFilePath string

	subCmdArgs := os.Args[2:]
	subCmdFlagSet := flag.NewFlagSet("operationFlagset", flag.ErrorHandling(flag.ExitOnError))
	subCmdFlagSet.StringVar(&remoteFilePath, "remotePath", "", "The remote file path to use within the operation")
	subCmdFlagSet.StringVar(&restoreFilePath, "restoreFilePath", "", "The directory location of where to restore the file(s)")
	subCmdFlagSet.Parse(subCmdArgs)

	switch operation {
	case "restore":
		log.Println("Performing Restore")
		if (restoreFilePath != "" && remoteFilePath != "") {
			restoreFile(restoreFilePath, remoteFilePath)
		} else {
			fmt.Println("Required options for restore operation:")
			subCmdFlagSet.PrintDefaults();
		}
	case "backup":
		log.Println("Performing Backup")
	case "age":
		log.Println("Performing Age Check")
	case "delete":
		log.Println("Performing Delete")
	default :
		log.Printf("freighter: '%s' is not a valid freighter argument \n", operation)
		log.Fatalln("See 'freighter --help' for more info \n")
	}

}

// TODO The below should be moved into a different file, behind an interface if possible
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

func showBanner() {

	fmt.Println("   ____        _      __   __         ")
	fmt.Println("  / __/______ (_)__ _/ /  / /____ ____")
	fmt.Println(" / _// __/ -_) / _ `/ _ \\/ __/ -_) __/")
	fmt.Println("/_/ /_/  \\__/_/\\_, /_//_/\\__/\\__/_/   ")
	fmt.Println("              /___/                   ")

}