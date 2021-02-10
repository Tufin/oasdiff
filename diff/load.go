package diff

import (
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

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
