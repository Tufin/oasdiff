package load_test

import (
	"net/url"
	"testing"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stretchr/testify/require"
	"github.com/tufin/oasdiff/load"
)

const RelativeDataPath = "../data/"

type MockLoader struct{}

func (swaggerMockLoader SwaggerMockLoader) LoadSwaggerFromFile(path string) (*openapi3.Swagger, error) {
	return openapi3.NewSwaggerLoader().LoadSwaggerFromFile(RelativeDataPath + path)
}

func (swaggerMockLoader SwaggerMockLoader) LoadSwaggerFromURI(location *url.URL) (*openapi3.Swagger, error) {
	return openapi3.NewSwaggerLoader().LoadSwaggerFromFile(RelativeDataPath + location.Path)
}

type SwaggerMockLoader struct{}

func NewSwaggerMockLoader() *load.SwaggerLoader {
	return &load.SwaggerLoader{
		Loader: SwaggerMockLoader{},
	}
}

func TestLoad_FileError(t *testing.T) {
	_, err := load.NewSwaggerLoader().From("null")
	require.Error(t, err)
}

func TestLoad_File(t *testing.T) {
	_, err := NewSwaggerMockLoader().From("openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoad_URI(t *testing.T) {
	_, err := NewSwaggerMockLoader().From("http://localhost/openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoad_URIError(t *testing.T) {
	_, err := NewSwaggerMockLoader().From("http://localhost/null")
	require.Error(t, err)
}
