package client

import (
	"testing"
	"os"
)

const arbitrary_restore_file_path string = "/foo/bar"
const arbitrary_remote_file_path string = "/moo/doo"

// Create a fake storage provider to decouple the client from a concrete provider
type fakeStorageProvider struct {
}

func (sp *fakeStorageProvider) DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string, err error) {
	return remoteFilePath, nil
}

func (sp *fakeStorageProvider) UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string, err error) {
	return remoteFilePath, nil
}

func (sp *fakeStorageProvider) DeleteFile(remoteFilePath string) (deleteFilePath string, err error) {
	return remoteFilePath, nil
}

func (sp *fakeStorageProvider) AgeFile(remoteFilePath string) (ageOfFile int, err error) {
	return 1, nil
}

func TestRestoreFile(t *testing.T) {

	// TODO Need to put archival and compression behind interface to fake out to unit test this properly
	t.SkipNow()

	sp := &fakeStorageProvider{}
	storageClient := NewStorageClient(sp)

	restoreFilePath := os.TempDir() + arbitrary_restore_file_path
	err := storageClient.RestoreFile(restoreFilePath, arbitrary_remote_file_path)
	if err != nil {
		t.Fatal(err)
	}
}
