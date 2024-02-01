package load_test

import (
	"log"
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestSpecInfo_File(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.NewSource("../data/openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestLoadInfo_FileWindows(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.NewSource(`C:\dev\OpenApi\spec2.yaml`))
	require.Condition(t, func() (success bool) {
		return err.Error() == "open C:\\dev\\OpenApi\\spec2.yaml: no such file or directory" ||
			err.Error() == "open C:/dev/OpenApi/spec2.yaml: The system cannot find the path specified."
	})
}

func TestLoadInfo_URI(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.NewSource("https://localhost/data/openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestLoadInfo_UriInvalid(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.NewSource("http://localhost/null"))
	require.Condition(t, func() (success bool) {
		return err.Error() == "open ../null: no such file or directory" ||
			err.Error() == "open ../null: The system cannot find the file specified."
	})
}

func TestLoadInfo_UriBadScheme(t *testing.T) {
	_, err := load.LoadSpecInfo(MockLoader{}, load.NewSource("ftp://localhost/null"))
	require.Condition(t, func() (success bool) {
		return err.Error() == "open ftp://localhost/null: no such file or directory" ||
			err.Error() == "open ftp://localhost/null: The filename, directory name, or volume label syntax is incorrect."
	})
}

func TestLoadInfo_Stdin(t *testing.T) {
	content := []byte(`openapi: 3.0.1
info:
  title: Test API
  version: v1
paths:
  /partner-api/test/some-method:
    get:
     responses:
       "200":
         description: Success
`)

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	_, err = load.LoadSpecInfo(MockLoader{}, load.NewSource("-"))
	require.NoError(t, err)
}

func TestLoadInfo_NoVersion(t *testing.T) {
	content := []byte(`openapi: 3.0.1
paths:
  /partner-api/test/some-method:
    get:
     responses:
       "200":
         description: Success
`)

	tmpfile, err := os.CreateTemp("", "example")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove(tmpfile.Name()) // clean up

	if _, err := tmpfile.Write(content); err != nil {
		log.Fatal(err)
	}

	if _, err := tmpfile.Seek(0, 0); err != nil {
		log.Fatal(err)
	}

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	specInfo, err := load.LoadSpecInfo(MockLoader{}, load.NewSource("-"))
	require.NoError(t, err)
	require.Empty(t, specInfo.Version)
}

func TestSpecInfo_GlobOK(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "../data/*.yaml")
	require.NoError(t, err)
}

func TestSpecInfo_InvalidSpec(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "../data/ignore-err-example.txt")
	require.EqualError(t, err, "error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type openapi3.TBis")
}

func TestSpecInfo_InvalidGlob(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "[*")
	require.EqualError(t, err, "syntax error in pattern")
}

func TestSpecInfo_URL(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "http://localhost/openapi-test1.yaml")
	require.EqualError(t, err, "no matching files (should be a glob, not a URL)")
}

func TestSpecInfo_GlobNoFiles(t *testing.T) {
	_, err := load.FromGlob(MockLoader{}, "../data/*.xxx")
	require.EqualError(t, err, "no matching files")
}
