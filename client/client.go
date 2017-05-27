package client

import (
	"log"
	"os"
	"github.com/opinari/freighter/archive"
	"github.com/opinari/freighter/compress"
	"github.com/opinari/freighter/storage"
	"io/ioutil"
	"strconv"
	"fmt"
	"errors"
)

// StorageClient abstracts the basic IO operations from a StorageProvider
//
// The storage client performs operations like compression / archival etc which are operations applicable
// regardless of the underlying StorageProvider.
type StorageClient interface {
	RestoreFile(restoreFilePath string, remoteFilePath string) error
	BackupDirectory(backupFilePath string, remoteFilePath string) error
	AgeRemoteFile(outputFilePath string, remoteFilePath string) error
	DeleteRemoteFile(remoteFilePath string) error
}

// Basic Filesystem implementation of the StorageClient
//
// Standard permissions are assumed for writing files / directories.
type FilesystemStorageClient struct {
	storageProvider storage.StorageProvider
}

func (sc *FilesystemStorageClient) RestoreFile(restoreFilePath string, remoteFilePath string) error {

	if restoreFilePath[:1] != "/" && restoreFilePath[:1] != "~" {
		return errors.New("restore path must be populated and refer to an absolute filesystem path")
	}

	// Create Restore File Path
	if err := os.MkdirAll(restoreFilePath, 0755); err != nil {
		return fmt.Errorf("error preparing restore path: %v", err)
	}

	//  Download File
	downloadedFilePath, err := sc.storageProvider.DownloadFile(restoreFilePath, remoteFilePath)
	if err != nil {
		return err
	}

	// Uncompress File
	compressedFile, err := compress.UncompressFile(downloadedFilePath, restoreFilePath)
	if err != nil {
		return err
	}

	// Unarchive Files
	_, err = archive.Unarchive(compressedFile, restoreFilePath)
	if err != nil {
		return err
	}

	// Cleanup tmp Downloaded and Uncompressed files
	if err := os.Remove(downloadedFilePath); err != nil {
		return fmt.Errorf("error cleaning up original download: %v", err)
	}

	if err := os.Remove(compressedFile); err != nil {
		return fmt.Errorf("error cleaning up temporary archive: %v", err)
	}

	return nil
}

func (sc *FilesystemStorageClient) BackupDirectory(backupFilePath string, remoteFilePath string) error {

	if backupFilePath[:1] != "/" && backupFilePath[:1] != "~" {
		return errors.New("backup path must be populated and refer to an absolute filesystem path")
	}

	// Archive Files
	archivedFilePath, err := archive.Archive(backupFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Compress File
	compressedFilePath, err := compress.CompressFile(archivedFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Upload File
	_, err = sc.storageProvider.UploadFile(compressedFilePath, remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Cleanup tmp Archived and Compressed files
	if err := os.Remove(archivedFilePath); err != nil {
		return fmt.Errorf("error cleaning up archive file: %v", err)
	}

	if err := os.Remove(compressedFilePath); err != nil {
		return fmt.Errorf("error cleaning up compressed file: %v", err)
	}

	return nil
}

func (sc *FilesystemStorageClient) AgeRemoteFile(outputFilePath string, remoteFilePath string) error {

	if outputFilePath[:1] != "/" && outputFilePath[:1] != "~" {
		return errors.New("age output path must be populated and refer to an absolute filesystem path")
	}

	// Perform age lookup
	age, err := sc.storageProvider.AgeFile(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	// Create and open the output file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		log.Fatal(err)
	}
	defer outputFile.Close()

	// Write out age value to file
	ageBytes := []byte(strconv.Itoa(age))
	if err := ioutil.WriteFile(outputFilePath, ageBytes, 0666); err != nil {
		log.Fatal(err)
	}

	return nil
}

func (sc *FilesystemStorageClient) DeleteRemoteFile(remoteFilePath string) error {

	_, err := sc.storageProvider.DeleteFile(remoteFilePath)
	if err != nil {
		log.Fatal(err)
	}

	return nil
}

func NewStorageClient(sp storage.StorageProvider) StorageClient {
	return &FilesystemStorageClient{storageProvider: sp}
}
