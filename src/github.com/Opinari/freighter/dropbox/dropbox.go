package dropbox

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"encoding/json"
)

const dropboxAPI = "https://content.dropboxapi.com/2"
const downloadPath = "/files/download"
const dropboxAPIArg = "Dropbox-API-Arg"

// TODO Fix this
const tmpDownloadedFileName = "/tmp.tar.gz"

var accessToken string

type DropboxFilePath struct {
	RemoteFilePath string `json:"path"`
}

// TODO Perhaps this should be passed via constructor / composite literal
func init() {
	if accessToken == "" {
		accessToken = os.Getenv("ACCESS_TOKEN")
	}
}

func DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string, err error) {

	log.Printf("Downloading file from: '%s' to: '%s'", remoteFilePath, restoreFilePath)

	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxAPI + downloadPath, nil)
	request.Header.Add("Authorization", "Bearer " + accessToken)

	// Construct API File Path in required json format
	apiPath := DropboxFilePath{RemoteFilePath: remoteFilePath}
	apiPathJsonBytes, err := json.Marshal(apiPath)
	apiPathJsonString := string(apiPathJsonBytes)
	request.Header.Add(dropboxAPIArg, apiPathJsonString)


	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("Error executing download request: %s", err.Error())
	}


	// Defer the response body close
	defer response.Body.Close()

	// Read the response
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", fmt.Errorf("Error occured whilst parsing body: %s", err.Error())
	}

	// TODO defer file close?
	// TODO Parse filename
	// Write response to file
	downloadedFilePath := restoreFilePath + tmpDownloadedFileName
	err = ioutil.WriteFile(downloadedFilePath, body, 0644)
	if err != nil {
		return "", fmt.Errorf("Error occured during saving of download: %s", err.Error())
	}

	log.Printf("Downloaded file succesfully to: %s", downloadedFilePath)

	// Return the path of the file of which was written
	downloadFilePath = downloadedFilePath;
	return
}

// TODO NO-OP
func UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string, err error) {

	log.Printf("Uploading File from: '%s' to: '%s' ", backupFilePath, remoteFilePath)
	return
}
