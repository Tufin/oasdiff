// Package load provides a function to load an OpenAPI spec from a URL or a Path
package load

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

// Loader interface includes the OAS load functions
type Loader interface {
	LoadSwaggerFromURI(*url.URL) (*openapi3.Swagger, error)
	LoadSwaggerFromFile(string) (*openapi3.Swagger, error)
}

// From is a convenience function that opens an OpenAPI spec from a URL or a local path based on the format of the path parameter
func From(loader Loader, path string) (*openapi3.Swagger, error) {

	uri, err := url.ParseRequestURI(path)
	if err == nil {
		return loadFromURI(loader, uri)
	}

	return loader.LoadSwaggerFromFile(path)
}

func loadFromURI(loader Loader, uri *url.URL) (*openapi3.Swagger, error) {
	oas, err := loader.LoadSwaggerFromURI(uri)
	if err != nil {
		return nil, err
	}
	return oas, nil
}
