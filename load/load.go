// Package load provides a function to load an OpenAPI spec from a URL or a Path
package load

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

// Loader interface is used to make openapi3 loaders testable
type Loader interface {
	LoadSwaggerFromURI(*url.URL) (*openapi3.Swagger, error)
	LoadSwaggerFromFile(string) (*openapi3.Swagger, error)
}

// OASLoader implements Loader interface
type OASLoader struct {
	Loader
}

// NewOASLoader returns a loader object that can be used to load OpenAPI specs
func NewOASLoader(loader Loader) *OASLoader {
	return &OASLoader{
		Loader: loader,
	}
}

// From is a convenience function that opens an OpenAPI spec from a URL or a local path based on the format of the path parameter
func (oasLoader *OASLoader) From(path string) (*openapi3.Swagger, error) {

	uri, err := url.ParseRequestURI(path)
	if err == nil {
		return oasLoader.loadFromURI(uri)
	}

	return oasLoader.LoadSwaggerFromFile(path)
}

func (oasLoader *OASLoader) loadFromURI(uri *url.URL) (*openapi3.Swagger, error) {
	oas, err := oasLoader.LoadSwaggerFromURI(uri)
	if err != nil {
		return nil, err
	}
	return oas, nil
}
