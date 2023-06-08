package load

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yargevad/filepathx"
)

type SpecInfo struct {
	Url  string
	Spec *openapi3.T
}

// LoadSpecInfoFromFile creates a SpecInfo from a local file path
func LoadSpecInfoFromFile(loader Loader, location string) (*SpecInfo, error) {
	s, err := loader.LoadFromFile(location)
	return &SpecInfo{Spec: s, Url: location}, err
}

// LoadSpecInfo creates a SpecInfo from a local file path or a URL
func LoadSpecInfo(loader Loader, location string) (*SpecInfo, error) {
	s, err := From(loader, location)
	return &SpecInfo{Spec: s, Url: location}, err
}

// FromGlob creates SpecInfo specs from local files matching the specified glob parameter
func FromGlob(loader Loader, glob string) ([]SpecInfo, error) {
	files, err := filepathx.Glob(glob)
	if err != nil {
		return nil, err
	}
	result := make([]SpecInfo, 0)
	for _, file := range files {
		spec, err := loader.LoadFromFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, SpecInfo{Url: file, Spec: spec})
	}

	return result, nil
}
