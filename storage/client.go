package storage

import (
	"fmt"
	"github.com/opinari/freighter/archive"
	"github.com/opinari/freighter/compress"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"errors"
)

// StorageClient abstracts the basic IO orchestration operations from a StorageProvider
//
// The storage client orchestrates operations like compression / archival etc which may be applicable
// regardless of the underlying StorageProvider.
type StorageClient interface {
	RestoreFile(remoteFilePath string, restoreFilePath string) error
	BackupDirectory(backupFilePath string, remoteFilePath string) error
	AgeRemoteFile(remoteFilePath string, outputFilePath string) error
	DeleteRemoteFile(remoteFilePath string) error
}

// Basic Filesystem implementation of the StorageClient
//
// Standard permissions are assumed for writing files / directories.
type FilesystemStorageClient struct {
	storageProvider StorageProvider
	compressors     []compress.Compressor
	archivers       []archive.Archiver
}

func (sc *FilesystemStorageClient) RestoreFile(remoteFilePath string, restoreFilePath string) error {

	if restoreFilePath[:1] != "/" && restoreFilePath[:1] != "~" {
		return errors.New("Restore path must refer to an absolute filesystem path")
	}

	if filepath.Ext(remoteFilePath) == "" {
		return errors.New("Only files with extensions, not folders, can be restored from storage providers")
	}

	// If no file name given within the restore filepath then assume the filename of the remote file is desired
	restoreFilePath = normalizeRestoreFilePath(restoreFilePath, remoteFilePath)

	// Make sure the restore base dir exists and create it if not so
	createIntermediateRestoreDirectories(restoreFilePath)

	//  Download File
	log.Printf("Downloading file from: '%s' to: '%s'", remoteFilePath, restoreFilePath)
	downloadedFilePath, err := sc.storageProvider.DownloadFile(remoteFilePath, restoreFilePath)
	if err != nil {
		return err
	}
	log.Printf("Downloaded file successfully to: %s", downloadedFilePath)

	// Resolve compression strategy, if no known match, skip and proceed
	compressor := sc.resolveCompressionStrategy(downloadedFilePath)
	if compressor != nil {

		defer cleanupIntermediateFile(downloadedFilePath)

		// We assume we want the contents of the uncompressed file to go in the same base dir as the download
		downloadedFilePath, err = compressor.UncompressFile(downloadedFilePath, filepath.Dir(restoreFilePath))
		if err != nil {
			return err
		}
	}

	// Resolve archival strategy, if no known match, skip and proceed
	archiver := sc.resolveArchiveStrategy(downloadedFilePath)
	if archiver != nil {

		defer cleanupIntermediateFile(downloadedFilePath)

		// We assume we want the contents of the unarchived files to go in the same base dir as the download
		_, err = archiver.Unarchive(downloadedFilePath, filepath.Dir(restoreFilePath))
		if err != nil {
			return err
		}
	}

	return nil
}

func (sc *FilesystemStorageClient) BackupDirectory(backupFilePath string, remoteFilePath string) error {

	if backupFilePath[:1] != "/" && backupFilePath[:1] != "~" {
		return errors.New("backup path must be populated and refer to an absolute filesystem path")
	}

	// Check backup directory actually exists
	if _, err := os.Stat(backupFilePath); os.IsNotExist(err) {
		return errors.New("backup file path must exist and have required read permissions")
	}

	// Archive Files
	// TODO Assume always Tar archiving when uploading for now
	archiver := archive.NewTarArchiver()
	archivedFilePath, err := archiver.Archive(backupFilePath)
	if err != nil {
		return err
	}
	defer cleanupIntermediateFile(archivedFilePath)

	// Compress File
	// TODO Assume always Gzip compression when uploading for now
	compressor := compress.NewGzipCompressor()
	compressedFilePath, err := compressor.CompressFile(archivedFilePath)
	if err != nil {
		return err
	}
	defer cleanupIntermediateFile(compressedFilePath)

	// Upload File
	_, err = sc.storageProvider.UploadFile(compressedFilePath, remoteFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (sc *FilesystemStorageClient) AgeRemoteFile(remoteFilePath string, outputFilePath string) error {

	// Perform age lookup
	age, err := sc.storageProvider.AgeFile(remoteFilePath)
	if err != nil {
		return err
	}

	if outputFilePath != "" {

		if outputFilePath[:1] != "/" && outputFilePath[:1] != "~" {
			return fmt.Errorf("age output path must be populated and refer to an absolute filesystem path")
		}

		// Create and open the output file for writing
		outputFile, err := os.Create(outputFilePath)
		if err != nil {
			return err
		}
		defer outputFile.Close()

		// Write out age value to file
		ageBytes := []byte(strconv.Itoa(age))
		if err := ioutil.WriteFile(outputFilePath, ageBytes, 0666); err != nil {
			return err
		}
	}

	return nil
}

func (sc *FilesystemStorageClient) DeleteRemoteFile(remoteFilePath string) error {

	_, err := sc.storageProvider.DeleteFile(remoteFilePath)
	if err != nil {
		return err
	}

	return nil
}

func (sc *FilesystemStorageClient) resolveCompressionStrategy(remoteFilePath string) compress.Compressor {

	ext := filepath.Ext(remoteFilePath)
	for _, compressor := range sc.compressors {
		if compressor.IsSupported(ext) {
			return compressor
		}
	}

	return nil
}

func (sc *FilesystemStorageClient) resolveArchiveStrategy(remoteFilePath string) archive.Archiver {

	ext := filepath.Ext(remoteFilePath)
	for _, archiver := range sc.archivers {
		if archiver.IsSupported(ext) {
			return archiver
		}
	}

	return nil
}

func normalizeRestoreFilePath(restoreFilePath, remoteFilePath string) string {

	isDir := filepath.Ext(restoreFilePath) == ""
	if isDir {
		if restoreFilePath[len(restoreFilePath)-1:] != "/" {
			restoreFilePath = restoreFilePath + "/"
		}
		_, remoteFileName := filepath.Split(remoteFilePath)
		restoreFilePath = restoreFilePath + remoteFileName
	}

	return restoreFilePath
}

func createIntermediateRestoreDirectories(restoreFilePath string) error {

	base := filepath.Dir(restoreFilePath)

	if err := os.MkdirAll(base, 0755); err != nil {
		return fmt.Errorf("error preparing restore path: %v", err)
	}

	return nil
}

func cleanupIntermediateFile(path string) error {

	if err := os.Remove(path); err != nil {
		return fmt.Errorf("error cleaning up file: %s %v", path, err)
	}

	return nil
}

func NewStorageClient(sp StorageProvider) StorageClient {

	// Instantiating the available compressors / archivers here for now, wouldn't DI be nice!
	compressors := []compress.Compressor{compress.NewGzipCompressor()}
	archivers := []archive.Archiver{archive.NewTarArchiver()}

	return &FilesystemStorageClient{storageProvider: sp, compressors: compressors, archivers: archivers}
}
