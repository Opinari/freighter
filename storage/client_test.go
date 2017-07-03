package storage

import (
	"testing"
	"os"
	"github.com/opinari/freighter/dropbox"
	"github.com/opinari/freighter/github"
)

const arbitrary_restore_file_path = "/foo/bar"
const arbitrary_remote_file_path = "/moo/doo"

const real_dropbox_remote_single_file_path = "/android.txt"
const real_dropbox_remote_archive_file_path = "/nginx.tar.gz"
const real_dropbox_restore_file_path = "/Users/Dave/tmp/dropbox/test/foo.txt"
const real_dropbox_restore_dir_path = "/Users/Dave/tmp/dropbox/test/"
const dropbox_acccess_token = "secret"

const real_github_remote_single_file_path = "/opinari/onthefence-data/contents/configuration/nginx/nginx.conf"
const real_github_remote_single_archive_path = "/opinari/onthefence-data/contents/configuration/nginx/nginx.conf"
const real_github_restore_file_path = "/Users/Dave/tmp/github/test/bar.txt"
const real_github_restore_dir_path = "/Users/Dave/tmp/github/test"
const github_oauth_token = "secret"


func TestRestoreFile(t *testing.T) {

	// TODO Need to put archival and compression behind interface to fake out to unit test this properly
	t.SkipNow()

	sp := &fakeStorageProvider{}
	storageClient := NewStorageClient(sp)

	restoreFilePath := os.TempDir() + arbitrary_restore_file_path
	err := storageClient.RestoreFile(arbitrary_remote_file_path, restoreFilePath)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRestoreFile_SingleFromDropbox(t *testing.T) {

	// TODO Uses real tokens and hits real APIs, just using for sanity
	t.SkipNow()

	sp := dropbox.NewDropboxStorageClient(dropbox_acccess_token)
	storageClient := NewStorageClient(sp)

	err := storageClient.RestoreFile(real_dropbox_remote_single_file_path, real_dropbox_restore_dir_path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRestoreFile_ArchiveFromDropbox(t *testing.T) {

	// TODO Uses real tokens and hits real APIs, just using for sanity
	t.SkipNow()

	sp := dropbox.NewDropboxStorageClient(dropbox_acccess_token)
	storageClient := NewStorageClient(sp)

	err := storageClient.RestoreFile(real_dropbox_remote_archive_file_path, real_dropbox_restore_dir_path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRestoreFile_SingleFromGithub(t *testing.T) {

	// TODO Uses real tokens and hits real APIs, just using for sanity
	t.SkipNow()

	sp := github.NewGithubStorageProvider(github_oauth_token)
	storageClient := NewStorageClient(sp)

	err := storageClient.RestoreFile(real_github_remote_single_file_path, real_github_restore_file_path)
	if err != nil {
		t.Fatal(err)
	}
}

func TestRestoreFile_ArchiveFromGithub(t *testing.T) {

	// TODO Uses real tokens and hits real APIs, just using for sanity
	t.SkipNow()

	sp := dropbox.NewDropboxStorageClient(dropbox_acccess_token)
	storageClient := NewStorageClient(sp)

	err := storageClient.RestoreFile(real_dropbox_remote_archive_file_path, real_dropbox_restore_dir_path)
	if err != nil {
		t.Fatal(err)
	}
}


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
