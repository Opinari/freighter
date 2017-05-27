package dropbox

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"encoding/json"
	"io"
	"gopkg.in/cheggaaa/pb.v1"
	"bufio"
	"time"
	"bytes"
	"io/ioutil"
	"github.com/opinari/freighter/storage"
)

const dropboxAPI = "https://api.dropboxapi.com/2"
const dropboxContentAPI = "https://content.dropboxapi.com/2"
const dropboxAPIArg = "Dropbox-API-Arg"

const downloadPath = "/files/download"
const tmpDownloadedFileName = "/tmp.tar.gz"

const uploadPath = "/files/upload"
const movePath = "/files/move"
const metadataPath = "/files/get_metadata"
const deletePath = "/files/delete"

type DropboxStorageProvider struct {
	accessToken string
}

type DropboxOptions struct {
	RemoteFilePath string `json:"path"`
}

type DropboxMoveOptions struct {
	RemoteFileFromPath string `json:"from_path"`
	RemoteFileToPath   string `json:"to_path"`
}

func (sp DropboxStorageProvider) DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string, err error) {

	log.Printf("Downloading file from: '%s' to: '%s'", remoteFilePath, restoreFilePath)

	// Construct API File Path in required json format
	apiPath := DropboxOptions{RemoteFilePath: remoteFilePath}
	apiPathJsonBytes, err := json.Marshal(apiPath)
	if err != nil {
		return "", fmt.Errorf("Error occured building json request payload: %s", err.Error())
	}
	apiPathJsonString := string(apiPathJsonBytes)

	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxContentAPI+downloadPath, nil)
	if err != nil {
		return "", fmt.Errorf("Error occured building http request: %s", err.Error())
	}
	request.Header.Add(dropboxAPIArg, apiPathJsonString)
	request.Header.Add("Authorization", "Bearer "+sp.accessToken)

	// Create File
	downloadedFilePath := restoreFilePath + tmpDownloadedFileName
	outputFile, err := os.Create(downloadedFilePath)
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
		return "", fmt.Errorf("Error downloading file, statusCode: %d, status: %s", response.StatusCode, response.Status)
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
	downloadFilePath = downloadedFilePath
	return
}

func (sp DropboxStorageProvider) UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string, err error) {

	tmpRemotePath := remoteFilePath + ".tmp"

	//
	log.Printf("Uploading new tmp file from: '%s' to: '%s' ", backupFilePath, tmpRemotePath)
	uploadFile(sp.accessToken, backupFilePath, tmpRemotePath)
	log.Printf("File was succesfully uploaded to: '%s'", tmpRemotePath)

	instant := time.Now().UTC().Format(time.RFC3339)
	archiveRemotePath := remoteFilePath + "-" + instant

	//
	log.Printf("Archiving old backup from: '%s' to: '%s' ", remoteFilePath, archiveRemotePath)
	_, err = moveFile(sp.accessToken, remoteFilePath, archiveRemotePath)
	if err != nil {
		log.Printf("Assuming new backup, Error occured moving file: %s", err.Error())
	} else {
		log.Printf("File was succesfully moved to: '%s'", archiveRemotePath)
	}

	//
	log.Printf("Setting primary from: '%s' to: '%s' ", tmpRemotePath, remoteFilePath)
	_, err = moveFile(sp.accessToken, tmpRemotePath, remoteFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured moving file: %s", err.Error())
	}
	log.Printf("Primary Backup File was succesfully set at: '%s'", remoteFilePath)

	uploadFilePath = remoteFilePath
	return
}

func (sp DropboxStorageProvider) DeleteFile(remoteFilePath string) (deleteFilePath string, err error) {

	// Construct API File Path in required json format
	apiPath := DropboxOptions{RemoteFilePath: remoteFilePath}
	apiPathJsonBytes, _ := json.Marshal(apiPath)
	if err != nil {
		return "", fmt.Errorf("Error occured building json request payload: %s", err.Error())
	}
	bodyReader := bytes.NewReader(apiPathJsonBytes)

	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxAPI+deletePath, bodyReader)
	if err != nil {
		return "", fmt.Errorf("Error occured building http request: %s", err.Error())
	}
	request.Header.Add("Authorization", "Bearer "+sp.accessToken)
	request.Header.Add("Content-Type", "application/json")

	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("Error executing delete request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Error deleting file, statusCode: %d, status: %s", response.StatusCode, response.Status)
	}

	log.Printf("File: '%s' was succesfully deleted", remoteFilePath)

	return
}

func (sp DropboxStorageProvider) AgeFile(remoteFilePath string) (ageOfFile int, err error) {

	// Construct API File Path in required json format
	apiPath := DropboxOptions{RemoteFilePath: remoteFilePath}
	apiPathJsonBytes, _ := json.Marshal(apiPath)
	if err != nil {
		return 0, fmt.Errorf("Error occured building json request payload: %s", err.Error())
	}
	bodyReader := bytes.NewReader(apiPathJsonBytes)

	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxAPI+metadataPath, bodyReader)
	if err != nil {
		return 0, fmt.Errorf("Error occured building http request: %s", err.Error())
	}
	request.Header.Add("Authorization", "Bearer "+sp.accessToken)
	request.Header.Add("Content-Type", "application/json")

	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return 0, fmt.Errorf("Error executing age request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return 0, fmt.Errorf("Error ageing file, statusCode: %d, status: %s", response.StatusCode, response.Status)
	}

	rawBody, err := ioutil.ReadAll(response.Body)

	var body map[string]interface{}
	if err := json.Unmarshal(rawBody, &body); err != nil {
		panic(err)
	}

	lastModifiedDateString := body["server_modified"].(string)
	log.Printf("Last modified date of file was: '%s' ", lastModifiedDateString)

	// parse date
	lastModifiedTime, err := time.Parse(time.RFC3339, lastModifiedDateString)
	timeNow := time.Now().UTC()
	timeDiffInDays := int(timeNow.Sub(lastModifiedTime.UTC()).Hours() / 24)
	log.Printf("Age of file was: '%d' day(s) ", timeDiffInDays)

	ageOfFile = timeDiffInDays
	return
}

func uploadFile(accessToken, backupFilePath, remoteFilePath string) (outputFilePath string, err error) {

	// Open File
	uploadFile, err := os.Open(backupFilePath)
	if err != nil {
		return "", fmt.Errorf("Error occured during opening of download file: %s", err.Error())
	}
	defer uploadFile.Close()

	fileInfo, err := uploadFile.Stat()
	if err != nil {
		return "", fmt.Errorf("Error occured reading file metadata: %s", err.Error())
	}

	reader := bufio.NewReader(uploadFile)

	// Create Proxy Reader for upload status
	progressBar := pb.New(int(fileInfo.Size())).SetUnits(pb.U_BYTES)
	progressBar.SetMaxWidth(100)
	progressBar.Start()
	proxyReader := progressBar.NewProxyReader(reader)

	// Construct API File Path in required json format
	apiPath := DropboxOptions{RemoteFilePath: remoteFilePath}
	apiPathJsonBytes, _ := json.Marshal(apiPath)
	if err != nil {
		return "", fmt.Errorf("Error occured building json request payload: %s", err.Error())
	}
	apiPathJsonString := string(apiPathJsonBytes)

	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxContentAPI+uploadPath, proxyReader)
	if err != nil {
		return "", fmt.Errorf("Error occured building http request: %s", err.Error())
	}
	request.Header.Add(dropboxAPIArg, apiPathJsonString)
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/octet-stream")

	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	progressBar.Finish()
	if err != nil {
		return "", fmt.Errorf("Error executing upload request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Error uploading file, statusCode: %d, status: %s", response.StatusCode, response.Status)
	}

	return
}

func moveFile(accessToken, fromPath, toPath string) (outputFilePath string, err error) {

	// Construct API File Path in required json format
	apiPath := DropboxMoveOptions{RemoteFileFromPath: fromPath, RemoteFileToPath: toPath}
	apiPathJsonBytes, _ := json.Marshal(apiPath)
	if err != nil {
		return "", fmt.Errorf("Error occured building json request payload: %s", err.Error())
	}
	bodyReader := bytes.NewReader(apiPathJsonBytes)

	// Build http Request
	request, err := http.NewRequest(http.MethodPost, dropboxAPI+movePath, bodyReader)
	if err != nil {
		return "", fmt.Errorf("Error occured building http request: %s", err.Error())
	}
	request.Header.Add("Authorization", "Bearer "+accessToken)
	request.Header.Add("Content-Type", "application/json")

	// Execute the http request
	response, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", fmt.Errorf("Error executing move request: %s", err.Error())
	}
	defer response.Body.Close()

	if response.StatusCode != 200 {
		return "", fmt.Errorf("Error moving file, statusCode: %d, status: %s", response.StatusCode, response.Status)
	}

	return
}

func NewDropboxStorageClient(backupProviderToken string) storage.StorageProvider {
	return &DropboxStorageProvider{accessToken: backupProviderToken}
}
