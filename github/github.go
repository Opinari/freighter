package github

import (
	"errors"
	"github.com/opinari/freighter/storage"
	"fmt"
	"net/http"
	"encoding/json"
	"io/ioutil"
	"os"
	"encoding/base64"
	"path/filepath"
)

const githubAPI = "https://api.github.com"
const acceptHeader = "application/vnd.github.v3+json"

const githubRepoPath = "/repos"

// A Github implementation of a storage provider which uses the Github API directly.
//
// Currently only supports the downloading of individual files from either a private or public repo.
// A valid OAuth API token must be specified for private repos and can be specified for public repos but if it
// is provided and invalid then the operation will fail regardless of its public visibility.
type GithubStorageProvider struct {
	oauthToken string
}

type GithubContentsResponse struct {
	FileType string `json:"type"`
	Encoding string
	Content  string
}

func (sp GithubStorageProvider) DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string, err error) {

	// Determine file vs directory
	isDir := filepath.Ext(remoteFilePath) == ""
	if isDir {
		return "", errors.New("Directories are not yet supported for this storage provder: " + remoteFilePath)
	}

	// Build http Request
	request, err := http.NewRequest(http.MethodGet, githubAPI+githubRepoPath+remoteFilePath, nil)
	if err != nil {
		return "", fmt.Errorf("Error occured building http request: %s", err.Error())
	}
	request.Header.Add("Accept", acceptHeader)
	if sp.oauthToken != "" {
		request.Header.Add("Authorization", "token "+sp.oauthToken)
	}

	// Create File
	outputFile, err := os.Create(restoreFilePath)
	if err != nil {
		return "", fmt.Errorf("Error opening download file: %s", err.Error())
	}
	defer outputFile.Close()

	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("Error executing download request: %s", err.Error())
	}
	defer response.Body.Close()
	if response.StatusCode != 200 {
		return "", errors.New("Error during request for content: " + response.Status)
	}

	rawBody, err := ioutil.ReadAll(response.Body)

	var contentsResp = GithubContentsResponse{}
	if err := json.Unmarshal(rawBody, &contentsResp); err != nil {
		return "", fmt.Errorf("Error unmarhsalling github reponse: %s", err.Error())
	}
	bytes, err := base64.StdEncoding.DecodeString(contentsResp.Content)
	outputFile.Write(bytes)

	return restoreFilePath, nil
}

// https://developer.github.com/v3/repos/contents/#update-a-file
// TODO NO-OP This will most likely just be a read only StorageProvider as committing, merging, pushing branches will be fraught with difficulty, probably...
func (sp GithubStorageProvider) UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string, err error) {
	return "", errors.New("This is a read only provider, this operation is not supported")
}

// https://developer.github.com/v3/repos/contents/#delete-a-file
// TODO NO-OP This will most likely just be a read only StorageProvider as committing, merging, pushing branches will be fraught with difficulty, probably...
func (sp GithubStorageProvider) DeleteFile(remoteFilePath string) (deleteFilePath string, err error) {
	return "", errors.New("This is a read only provider, this operation is not supported")
}

// TODO NO-OP Haven't stumbled upon individual file metadata within the contents API section, yet...
func (sp GithubStorageProvider) AgeFile(remoteFilePath string) (ageOfFile int, err error) {
	return 0, errors.New("This operation is not yet supported")
}

func NewGithubStorageProvider(backupProviderToken string) storage.StorageProvider {

	// Overrides the default backup provider token if the specific Git token is provided
	gitAPIToken := os.Getenv("GIT_API_TOKEN")
	if gitAPIToken != "" {
		backupProviderToken = gitAPIToken
	}

	return &GithubStorageProvider{oauthToken: backupProviderToken}
}
