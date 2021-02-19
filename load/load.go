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

type SwaggerLoader struct {
	Loader Loader
}

func NewSwaggerLoader() *SwaggerLoader {
	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true
	return &SwaggerLoader{
		Loader: loader,
	}
}

// From is a convenience function that opens a swagger spec from a URL or a local path based on the format of the path parameter
func (swaggerLoader *SwaggerLoader) From(path string) (*openapi3.Swagger, error) {

	uri, err := url.ParseRequestURI(path)
	if err == nil {
		swagger, err := swaggerLoader.FromURI(uri)
		if err != nil {
			return nil, err
		}
		return swagger, nil
	}

	swagger, err := swaggerLoader.FromPath(path)
	if err != nil {
		return nil, err
	}

	return swagger, nil
}

// FromPath opens a swagger spec from a local path
func (swaggerLoader *SwaggerLoader) FromPath(path string) (*openapi3.Swagger, error) {

	swagger, err := swaggerLoader.Loader.LoadSwaggerFromFile(path)
	if err != nil {
		log.Errorf("failed to open swagger from '%s' with '%v'", path, err)
		return nil, err
	}

	return swagger, nil
}

// FromURI opens a swagger spec from a URL
func (swaggerLoader *SwaggerLoader) FromURI(path *url.URL) (*openapi3.Swagger, error) {

	swagger, err := swaggerLoader.Loader.LoadSwaggerFromURI(path)
	if err != nil {
		log.Errorf("failed to open swagger from '%s' with '%v'", path, err)
		return nil, err
	}

	return swagger, nil
}
