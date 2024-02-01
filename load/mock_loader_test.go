package load_test

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

func (mockLoader MockLoader) LoadFromFile(path string) (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromFile(path)
}

func (mockLoader MockLoader) LoadFromURI(location *url.URL) (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromFile(".." + location.Path)
}

func (mockLoader MockLoader) LoadFromStdin() (*openapi3.T, error) {
	return openapi3.NewLoader().LoadFromStdin()
}

type MockLoader struct{}
