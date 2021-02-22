package load

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

type Loader interface {
	LoadSwaggerFromURI(*url.URL) (*openapi3.Swagger, error)
	LoadSwaggerFromFile(string) (*openapi3.Swagger, error)
}

type OASLoader struct {
	Loader Loader
}

// NewOASLoader returns a loader object that can be used to load OpenAPI specs
func NewOASLoader() *OASLoader {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	return &OASLoader{
		Loader: loader,
	}
}

// From is a convenience function that opens an OpenAPI spec from a URL or a local path based on the format of the path parameter
func (oasLoader *OASLoader) From(path string) (*openapi3.Swagger, error) {

	uri, err := url.ParseRequestURI(path)
	if err == nil {
		oas, err := oasLoader.FromURI(uri)
		if err != nil {
			return nil, err
		}
		return oas, nil
	}

	oas, err := oasLoader.FromPath(path)
	if err != nil {
		return nil, err
	}

	return oas, nil
}

// FromPath opens an OpenAPI spec from a local path
func (oasLoader *OASLoader) FromPath(path string) (*openapi3.Swagger, error) {

	oas, err := oasLoader.Loader.LoadSwaggerFromFile(path)
	if err != nil {
		log.Errorf("failed to open spec from '%s' with '%v'", path, err)
		return nil, err
	}

	return oas, nil
}

// FromURI opens an OpenAPI spec from a URL
func (oasLoader *OASLoader) FromURI(path *url.URL) (*openapi3.Swagger, error) {

	oas, err := oasLoader.Loader.LoadSwaggerFromURI(path)
	if err != nil {
		log.Errorf("failed to open spec from '%s' with '%v'", path, err)
		return nil, err
	}

	return oas, nil
}
