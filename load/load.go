package load

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yargevad/filepathx"
)

// Loader interface includes the OAS load functions
type Loader interface {
	LoadFromURI(*url.URL) (*openapi3.T, error)
	LoadFromFile(string) (*openapi3.T, error)
}

// From is a convenience function that opens an OpenAPI spec from a URL or a local path based on the format of the path parameter
func From(loader Loader, path string) (*openapi3.T, error) {

	uri, err := url.ParseRequestURI(path)
	if err == nil {
		return loadFromURI(loader, uri)
	}

	// return loader.LoadFromURI(&url.URL{Path: filepath.ToSlash(path)})
	return loader.LoadFromFile(path)
}

// From is a convenience function that opens OpenAPI specs from a local iles matching the specified glob parameter
func FromGlob(loader Loader, glob string) (*[]*openapi3.T, error) {
	files, err := filepathx.Glob(glob)
	if err != nil {
		return nil, err
	}
	result := make([]*openapi3.T, 0)
	for _, file := range files {
		api, err := loader.LoadFromFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, api)
	}

	return &result, nil
}

func loadFromURI(loader Loader, uri *url.URL) (*openapi3.T, error) {
	oas, err := loader.LoadFromURI(uri)
	if err != nil {
		return nil, err
	}
	return oas, nil
}
