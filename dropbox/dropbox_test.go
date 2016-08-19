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
func IgnoreTestDownloadFile(t *testing.T) {

	for _, test := range tests {
		_, err := DownloadFile(test.restoreFilePath, test.remoteFilePath)

		if err.Error() != "foo" {
			t.Fail()
		}
	}
}
