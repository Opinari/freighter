package github

import (
	"testing"
	"os"
	"io/ioutil"
	"fmt"
)

const exampleRealPublicFilePath = "/opinari/freighter/contents/storage/provider.go"
const exampleRealPrivateFilePath = "/opinari/onthefence-data/contents/configuration/prometheus/prometheus.yml"

const exampleRealPublicFolderPath = "/opinari/freighter/contents/client"
const exampleRealPrivateFolderPath = "/opinari/onthefence-data/contents/configuration/"

const outputFilename = "provider.go"
const emptyOauthToken = ""

// This is more of an integration test than a unit test, but gives some confidence nonetheless
func TestDownloadFile_ForIndividualFile(t *testing.T) {

	// TODO Can't write to any tmp space on travis...
	t.SkipNow()

	sp := NewGithubStorageProvider(emptyOauthToken)

	tmpRestorePath := os.TempDir() + outputFilename

	outputPath, err := sp.DownloadFile(exampleRealPublicFilePath, tmpRestorePath)
	if err != nil {
		t.Fatal(err)
	}

	if outputPath != tmpRestorePath {
		t.Fatalf("expected path %s does not match actual path %s", tmpRestorePath, outputPath)
	}

	// Check the file exists and is not empty
	fileInfo, err := os.Stat(outputPath)
	if os.IsNotExist(err) {
		t.Fatal(err)
	}

	if fileInfo.Size() < 1 {
		t.Fatal("file is empty")
	}

	// Finally, print it for sanity's sake
	printFileContents(outputPath, t)
}

// Check that we fail gracefully for any provided directory
func TestDownloadFile_ForDirectory(t *testing.T) {

	pathsTreatedAsDirs := []string{
		"/foo/bar",
		"foo/bar",
		"/bar",
		"bar",
		" ",
		"a ",
		"/a ",
	}

	sp := NewGithubStorageProvider(emptyOauthToken)

	for _, path := range pathsTreatedAsDirs {

		t.Run(fmt.Sprintf("path:%s", path), func(t *testing.T) {

			tmpRestorePath := os.TempDir() + "test.file"

			_, err := sp.DownloadFile(path, tmpRestorePath)
			if err == nil {
				t.Fatal("unsupported error for directories expected")
			}

			expectedErr := "Directories are not yet supported for this storage provder: " + path
			if err.Error() != expectedErr {
				t.Fatalf("expected error: %s but got error: %s", expectedErr, err)
			}
		})
	}
}

func printFileContents(outputPath string, t *testing.T) {

	t.Log("output filepath is: " + outputPath)
	bytes, _ := ioutil.ReadFile(outputPath)
	t.Log(string(bytes))
}
