package dropbox

import "testing"

type Test struct {
	restoreFilePath          string
	remoteFilePath           string
	expectedDownloadFilePath string
}

var tests = []Test{
	{"/foo/bar", "/bar.tar.gz", "/foo/bar"},
}

// FIXME write these tests properly
func TestDownloadFile(t *testing.T) {

	backupProviderToken := "fooBarToken"
	sp := DropboxStorageProvider{accessToken: backupProviderToken}

	t.Skip("skipping test.")

	for _, test := range tests {

		_, err := sp.DownloadFile(test.remoteFilePath, test.restoreFilePath)

		if err.Error() != "foo" {
			t.Fail()
		}
	}

}
