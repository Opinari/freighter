package dropbox

import (
	"fmt"
	"net/http"
	"io/ioutil"
)

func DownloadFile(restoreFilePath string, remoteFilePath string) (downloadFilePath string) {

	fmt.Printf("Downloading File from: '%s' to: '%s'", remoteFilePath, restoreFilePath);

	response, err := http.Get("https://google.com")
	if err != nil {
		// FIXME don't panic, return errs instead
		panic("Something really went wrong Downloading")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// FIXME don't panic, return errs instead
		panic("Something really went wrong Reading the Body")
	}
	fmt.Printf("Body was: %s!", body)

	downloadFilePath = restoreFilePath
	return
}

func UploadFile(backupFilePath string, remoteFilePath string) (uploadFilePath string) {

	fmt.Printf("Uploading File from: '%s' to: '%s' ", backupFilePath, remoteFilePath);

	response, err := http.Get("https://google.com")
	if err != nil {
		// FIXME don't panic, return errs instead
		panic("Something really went wrong Uploading")
	}

	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		// FIXME don't panic, return errs instead
		panic("Something really went wrong")
	}
	fmt.Printf("Body was: %s!", body)

	uploadFilePath = remoteFilePath
	return
}
