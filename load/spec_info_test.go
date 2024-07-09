package load_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

func TestSpecInfo_File(t *testing.T) {
	_, err := load.NewSpecInfo(MockLoader{}, load.NewSource("../data/openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestLoadInfo_URI(t *testing.T) {
	_, err := load.NewSpecInfo(MockLoader{}, load.NewSource("https://localhost/data/openapi-test1.yaml"))
	require.NoError(t, err)
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
	require.NoError(t, err)

	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.Write(content)
	require.NoError(t, err)

	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	_, err = load.NewSpecInfo(MockLoader{}, load.NewSource("-"))
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
	require.NoError(t, err)

	defer os.Remove(tmpfile.Name()) // clean up

	_, err = tmpfile.Write(content)
	require.NoError(t, err)

	_, err = tmpfile.Seek(0, 0)
	require.NoError(t, err)

	oldStdin := os.Stdin
	defer func() { os.Stdin = oldStdin }() // Restore original Stdin

	os.Stdin = tmpfile
	specInfo, err := load.NewSpecInfo(MockLoader{}, load.NewSource("-"))
	require.NoError(t, err)
	require.Empty(t, specInfo.Version)
}

func TestSpecInfo_GlobOK(t *testing.T) {
	_, err := load.NewSpecInfoFromGlob(MockLoader{}, "../data/*.yaml")
	require.NoError(t, err)
}

func TestSpecInfo_InvalidSpec(t *testing.T) {
	_, err := load.NewSpecInfoFromGlob(MockLoader{}, "../data/ignore-err-example.txt")
	require.EqualError(t, err, "failed to load \"../data/ignore-err-example.txt\": failed to unmarshal data: json error: invalid character 'G' looking for beginning of value, yaml error: error unmarshaling JSON: while decoding JSON: json: cannot unmarshal string into Go value of type openapi3.TBis")
}

func TestSpecInfo_InvalidGlob(t *testing.T) {
	_, err := load.NewSpecInfoFromGlob(MockLoader{}, "[*")
	require.EqualError(t, err, "syntax error in pattern")
}

func TestSpecInfo_URL(t *testing.T) {
	_, err := load.NewSpecInfoFromGlob(MockLoader{}, "http://localhost/openapi-test1.yaml")
	require.EqualError(t, err, "no matching files (should be a glob, not a URL)")
}

func TestSpecInfo_GlobNoFiles(t *testing.T) {
	_, err := load.NewSpecInfoFromGlob(MockLoader{}, "../data/*.xxx")
	require.EqualError(t, err, "no matching files")
}

func TestSpecInfo_Options(t *testing.T) {
	_, err := load.NewSpecInfo(MockLoader{}, load.NewSource("../data/openapi-test1.yaml"), load.GetOption(load.WithFlattenAllOf(), false), load.GetOption(load.WithFlattenAllOf(), true), load.WithFlattenParams())
	require.NoError(t, err)
}

func TestSpecInfo_GlobOptions(t *testing.T) {
	_, err := load.NewSpecInfoFromGlob(MockLoader{}, "../data/*.yaml", load.WithIdentity(), load.WithFlattenAllOf(), load.WithFlattenParams())
	require.NoError(t, err)
}
