package load

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yargevad/filepathx"
)

type OpenAPISpecInfo struct {
	Url  string
	Spec *openapi3.T
}

// LoadOpenAPISpecInfoFromFile loads a LoadOpenAPISpecInfoFromFile from a local file path
func LoadOpenAPISpecInfoFromFile(loader Loader, location string) (*OpenAPISpecInfo, error) {
	s, err := loader.LoadFromFile(location)
	return &OpenAPISpecInfo{Spec: s, Url: location}, err
}

func LoadOpenAPISpecInfo(loader Loader, location string) (*OpenAPISpecInfo, error) {
	s, err := From(loader, location)
	return &OpenAPISpecInfo{Spec: s, Url: location}, err
}

// FromGlob is a convenience function that opens OpenAPI specs from local files matching the specified glob parameter
func FromGlob(loader Loader, glob string) ([]OpenAPISpecInfo, error) {
	files, err := filepathx.Glob(glob)
	if err != nil {
		return nil, err
	}
	result := make([]OpenAPISpecInfo, 0)
	for _, file := range files {
		spec, err := loader.LoadFromFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, OpenAPISpecInfo{Url: file, Spec: spec})
	}

	return result, nil
}
