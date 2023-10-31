package load

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

// Loader interface includes the OAS load functions
type Loader interface {
	LoadFromURI(*url.URL) (*openapi3.T, error)
	LoadFromFile(string) (*openapi3.T, error)
	LoadFromStdin() (*openapi3.T, error)
}

// From is a convenience function that opens an OpenAPI spec from a URL or a local path based on the format of the path parameter
func From(loader Loader, source Source) (*openapi3.T, error) {

	if source.Stdin {
		return loader.LoadFromStdin()
	}

	uri, err := url.ParseRequestURI(source.Path)
	if err == nil {
		return loadFromURI(loader, uri)
	}

	return loader.LoadFromFile(source.Path)
}

func loadFromURI(loader Loader, uri *url.URL) (*openapi3.T, error) {
	oas, err := loader.LoadFromURI(uri)
	if err != nil {
		return nil, err
	}
	return oas, nil
}
