package load_test

import (
	"net/url"
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

type MockLoader struct{}

func TestLoad_File(t *testing.T) {
	_, err := load.From(MockLoader{}, "openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoad_URI(t *testing.T) {
	_, err := load.From(MockLoader{}, "http://localhost/openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoad_URIError(t *testing.T) {
	_, err := load.From(MockLoader{}, "http://localhost/null")
	require.Error(t, err)
}
