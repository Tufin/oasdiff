package load_test

import (
	"log"
	"net/url"
	"os"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

const RelativeDataPath = "../data/"

func (mockLoader MockLoader) LoadFromFile(path string) (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromFile(RelativeDataPath + path)
}

func (mockLoader MockLoader) LoadFromURI(location *url.URL) (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromFile(RelativeDataPath + location.Path)
}

func (mockLoader MockLoader) LoadFromStdin() (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromStdin()
}

type MockLoader struct{}

func TestLoad_File(t *testing.T) {
	_, err := load.From(MockLoader{}, load.GetSource("openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestLoad_URI(t *testing.T) {
	_, err := load.From(MockLoader{}, load.GetSource("http://localhost/openapi-test1.yaml"))
	require.NoError(t, err)
}

func TestLoad_URIError(t *testing.T) {
	_, err := load.From(MockLoader{}, load.GetSource("http://localhost/null"))
	require.Error(t, err)
}

func TestLoad_Stdin(t *testing.T) {
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
	_, err = load.From(MockLoader{}, load.GetSource("-"))
	require.NoError(t, err)
}
