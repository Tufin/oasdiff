package load

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

// Load is a convenience function that opens a swagger spec from a URL or a local path based on the format of the path parameter
func Load(path string) (*openapi3.Swagger, error) {

	uri, err := url.ParseRequestURI(path)
	if err == nil {
		swagger, err := LoadURI(uri)
		if err != nil {
			return nil, err
		}
		return swagger, nil
	}

	swagger, err := LoadPath(path)
	if err != nil {
		return nil, err
	}

	return swagger, nil
}

// LoadPath opens a swagger spec from a local path
func LoadPath(path string) (*openapi3.Swagger, error) {

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true

	swagger, err := loader.LoadSwaggerFromFile(path)
	if err != nil {
		log.Errorf("failed to open swagger from '%s' with '%v'", path, err)
		return nil, err
	}

	return swagger, nil
}

// LoadPath opens a swagger spec from a URL
func LoadURI(path *url.URL) (*openapi3.Swagger, error) {

	loader := openapi3.NewSwaggerLoader()
	loader.IsExternalRefsAllowed = true

	swagger, err := loader.LoadSwaggerFromURI(path)
	if err != nil {
		log.Errorf("failed to open swagger from '%s' with '%v'", path, err)
		return nil, err
	}

	return swagger, nil
}
