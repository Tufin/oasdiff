package load_test

import (
	"net/url"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

const RelativeDataPath = "../data/"

func (mockLoader MockLoader) LoadSwaggerFromFile(path string) (*openapi3.Swagger, error) {
	return openapi3.NewSwaggerLoader().LoadSwaggerFromFile(RelativeDataPath + path)
}

func (mockLoader MockLoader) LoadSwaggerFromURI(location *url.URL) (*openapi3.Swagger, error) {
	return openapi3.NewSwaggerLoader().LoadSwaggerFromFile(RelativeDataPath + location.Path)
}

type MockLoader struct{}

func newMockLoader() *load.OASLoader {
	return &load.OASLoader{
		Loader: MockLoader{},
	}
}

func TestLoad_FileError(t *testing.T) {
	_, err := load.NewOASLoader().From("null")
	require.Error(t, err)
}

func TestLoad_File(t *testing.T) {
	_, err := newMockLoader().From("openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoad_URI(t *testing.T) {
	_, err := newMockLoader().From("http://localhost/openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoad_URIError(t *testing.T) {
	_, err := newMockLoader().From("http://localhost/null")
	require.Error(t, err)
}
