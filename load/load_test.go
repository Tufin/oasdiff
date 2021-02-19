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

func TestLoadPath_NoError(t *testing.T) {
	_, err := NewSwaggerMockLoader().FromPath("openapi-test1.yaml")
	require.NoError(t, err)
}

func TestLoadPath_Error(t *testing.T) {
	_, err := NewSwaggerMockLoader().FromPath("null")
	require.Error(t, err)
}

func TestLoadURI_NoError(t *testing.T) {
	path, err := url.ParseRequestURI("http://localhost/openapi-test1.yaml")
	require.NoError(t, err)

	_, err = NewSwaggerMockLoader().FromURI(path)
	require.NoError(t, err)
}

func TestLoadURI_Error(t *testing.T) {
	path, err := url.ParseRequestURI("http://localhost/null")
	require.NoError(t, err)

	_, err = NewSwaggerMockLoader().FromURI(path)
	require.Error(t, err)
}

func TestLoad_NoError(t *testing.T) {
	_, err := NewSwaggerMockLoader().From("openapi-test1.yaml")
	require.NoError(t, err)
}
