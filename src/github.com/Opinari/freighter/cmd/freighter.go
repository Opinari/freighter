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

var (
	remoteFilePath string
	restoreFilePath string
)


// TODO This will parse the command line args and decide what to operation to invoke
func main() {

	// TODO Do Some ASCII Magic here
	fmt.Println("Freighter 2016.")

	// Put in some init?
	log.SetOutput(os.Stdout)
	runCLI();
}

// TODO Make the CLI an instance and execute the args / flags upon it
func runCLI() {

	// Sub cmd arg parsing
	operation := os.Args[1]

	// Flag parsing
	subCmdArgs := os.Args[2:]
	subCmdFlagSet := flag.NewFlagSet("restore", flag.ErrorHandling(1))
	subCmdFlagSet.StringVar(&remoteFilePath, "remotePath", "", "The remote file path to use within the operation")
	subCmdFlagSet.StringVar(&restoreFilePath, "restoreFilePath", "", "The directory location of where to restore the file(s)")
	subCmdFlagSet.Parse(subCmdArgs)

	switch operation {
	case "restore":
		log.Println("Restore invoked")
		if (restoreFilePath != "" && remoteFilePath != "") {
			restoreFile(restoreFilePath, remoteFilePath)
		} else {
			flag.Usage();
		}
	case "backup":
		log.Println("Backup invoked")

	case "age":
		log.Println("Age invoked")
	case "delete":
		log.Println("Delete invoked")
	default :
		fmt.Fprintf(os.Stderr, "freighter: '%s' is not a valid freighter argument \n", operation)
		fmt.Fprint(os.Stderr, "See 'freighter --help' for more info \n")
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
