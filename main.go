package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"github.com/Opinari/freighter/client"
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
			client.RestoreFile(restoreFilePath, remoteFilePath)
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

func showBanner() {

	fmt.Println("   ____        _      __   __         ")
	fmt.Println("  / __/______ (_)__ _/ /  / /____ ____")
	fmt.Println(" / _// __/ -_) / _ `/ _ \\/ __/ -_) __/")
	fmt.Println("/_/ /_/  \\__/_/\\_, /_//_/\\__/\\__/_/   ")
	fmt.Println("              /___/                   ")

}
