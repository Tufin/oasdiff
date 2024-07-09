package load_test

import (
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestLoadInfo_FileWindows(t *testing.T) {
	_, err := load.NewSpecInfo(MockLoader{}, load.NewSource(`C:\dev\OpenApi\spec2.yaml`))
	require.EqualError(t, err, "open C:\\dev\\OpenApi\\spec2.yaml: The system cannot find the path specified.")
}

func TestLoadInfo_UriInvalid(t *testing.T) {
	_, err := load.NewSpecInfo(MockLoader{}, load.NewSource("http://localhost/null"))
	require.EqualError(t, err, "open ..\\null: The system cannot find the file specified.")
}

func TestLoadInfo_UriBadScheme(t *testing.T) {
	_, err := load.NewSpecInfo(MockLoader{}, load.NewSource("ftp://localhost/null"))
	require.EqualError(t, err, "open ftp:\\\\localhost\\null: The filename, directory name, or volume label syntax is incorrect.")
}
