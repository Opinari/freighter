package dropbox

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io"
	"gopkg.in/cheggaaa/pb.v1"
)

const dropboxAPI = "https://content.dropboxapi.com/2"
const downloadPath = "/files/download"
const dropboxAPIArg = "Dropbox-API-Arg"
const tmpDownloadedFileName = "/tmp.tar.gz"

var accessToken string

// TODO Perhaps this should be passed via constructor / composite literal
func init() {
	if accessToken == "" {
		accessToken = os.Getenv("ACCESS_TOKEN")
	}
	if accessToken == "" {
		log.Fatalln("Error: No Access token provided for storage provider dropbox")
	}
}

func DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string, err error) {

	log.Printf("Downloading file from: '%s' to: '%s'", remoteFilePath, restoreFilePath)


	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxAPI + downloadPath, nil)
	request.Header.Add("Authorization", "Bearer " + accessToken)


	// Construct API File Path in required json format
	type DropboxFilePath struct{ RemoteFilePath string `json:"path"` }
	apiPath := DropboxFilePath{RemoteFilePath: remoteFilePath}
	apiPathJsonBytes, err := json.Marshal(apiPath)
	apiPathJsonString := string(apiPathJsonBytes)
	request.Header.Add(dropboxAPIArg, apiPathJsonString)


	// Create File
	downloadedFilePath := restoreFilePath + tmpDownloadedFileName
	outputFile, err := os.Create(downloadedFilePath)
	defer outputFile.Close()
	if err != nil {
		return "", fmt.Errorf("Error occured during saving of download: %s", err.Error())
	}


	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	defer response.Body.Close()
	if err != nil {
		return "", fmt.Errorf("Error executing download request: %s", err.Error())
	}


	// Create Proxy Reader for download status
	progressBar := pb.New(int(response.ContentLength)).SetUnits(pb.U_BYTES)
	progressBar.SetMaxWidth(100)
	progressBar.Start()
	proxyReader := progressBar.NewProxyReader(response.Body)


	// Copy from http body reader to output file writer
	_, err = io.Copy(outputFile, proxyReader)
	if err != nil {
		return "", fmt.Errorf("Error occured whilst parsing body: %s", err.Error())
	}

	progressBar.Finish()
	log.Printf("Downloaded file successfully to: %s", downloadedFilePath)


	// Return the path of the file of which was written
	downloadFilePath = downloadedFilePath;
	return
}


// TODO NO-OP
func UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string, err error) {

	log.Printf("Uploading File from: '%s' to: '%s' ", backupFilePath, remoteFilePath)
	return
}
