package main

import (
	"flag"
	"fmt"
	"github.com/opinari/freighter/client"
	"github.com/opinari/freighter/dropbox"
	"github.com/opinari/freighter/github"
	"github.com/opinari/freighter/storage"
	"log"
	"os"
)

const version = "0.4.0"

func init() {
	log.SetOutput(os.Stdout)
}

func main() {
	showBanner()
	runCLI()
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
	var storageProvider, remoteFilePath, restoreFilePath, backupFilePath, ageOutputFilePath string

	subCmdArgs := os.Args[2:]
	subCmdFlagSet := flag.NewFlagSet("operationFlagset", flag.ErrorHandling(flag.ExitOnError))

	// Defaulting to Dropbox for now due to backwards compatibility
	subCmdFlagSet.StringVar(&storageProvider, "storageProvider", "dropbox", "The storage backend provider to use")
	subCmdFlagSet.StringVar(&remoteFilePath, "remoteFilePath", "", "The remote file path to use within the operation")
	subCmdFlagSet.StringVar(&restoreFilePath, "restoreFilePath", "", "The directory location of where to restore the file(s)")
	subCmdFlagSet.StringVar(&backupFilePath, "backupFilePath", "", "The path to the directory of which to backup")
	subCmdFlagSet.StringVar(&ageOutputFilePath, "ageOutputFilePath", "", "The path to the directory of which to backup")
	subCmdFlagSet.Parse(subCmdArgs)

	backupProviderToken := resolveBackupProviderToken()
	sp := resolveStorageProvider(storageProvider, backupProviderToken)
	storageClient := client.NewStorageClient(sp)

	switch operation {
	case "restore":
		log.Println("Performing Restore")
		if restoreFilePath != "" && remoteFilePath != "" {
			if err := storageClient.RestoreFile(restoreFilePath, remoteFilePath); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Required options for restore operation:")
			subCmdFlagSet.PrintDefaults()
		}
	case "backup":
		log.Println("Performing Backup")
		if backupFilePath != "" && remoteFilePath != "" {
			if err := storageClient.BackupDirectory(backupFilePath, remoteFilePath); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Required options for backup operation:")
			subCmdFlagSet.PrintDefaults()
		}
	case "age":
		log.Println("Performing Age Check")
		if ageOutputFilePath != "" && remoteFilePath != "" {
			if err := storageClient.AgeRemoteFile(ageOutputFilePath, remoteFilePath); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Required options for age operation:")
			subCmdFlagSet.PrintDefaults()
		}
	case "delete":
		log.Println("Performing Delete")
		if remoteFilePath != "" {
			if err := storageClient.DeleteRemoteFile(remoteFilePath); err != nil {
				log.Fatal(err)
			}
		} else {
			fmt.Println("Required options for delete operation:")
			subCmdFlagSet.PrintDefaults()
		}
	default:
		log.Printf("freighter: '%s' is not a valid freighter argument \n", operation)
		log.Fatalln("See 'freighter --help' for more info \n")
	}

}

// TODO Provide alternative strategies to provide this other than just env var
func resolveBackupProviderToken() string {

	backupProviderToken := os.Getenv("BACKUP_PROVIDER_TOKEN")
	return backupProviderToken
}

func resolveStorageProvider(storageProvider, backupProviderToken string) storage.StorageProvider {

	switch storageProvider {

	case "dropbox":
		log.Println("Using Dropbox as a Storage Provider")
		return dropbox.NewDropboxStorageClient(backupProviderToken)

	case "github":
		log.Println("Using Github as a Storage Provider")
		return github.NewGithubStorageProvider(backupProviderToken)

	default:
		log.Fatalln("Error: Invalid Provider Type given: " + storageProvider)
		return nil
	}
}

func showBanner() {

	fmt.Println("   ____        _      __   __         ")
	fmt.Println("  / __/______ (_)__ _/ /  / /____ ____")
	fmt.Println(" / _// __/ -_) / _ `/ _ \\/ __/ -_) __/")
	fmt.Println("/_/ /_/  \\__/_/\\_, /_//_/\\__/\\__/_/   ")
	fmt.Println("              /___/                   v" + version)
}
