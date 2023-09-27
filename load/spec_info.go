package load

import (
	"errors"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yargevad/filepathx"
)

type SpecInfo struct {
	Url  string
	Spec *openapi3.T
}

// LoadSpecInfo creates a SpecInfo from a local file path or a URL
func LoadSpecInfo(loader Loader, location string) (*SpecInfo, error) {
	s, err := From(loader, location)
	return &SpecInfo{Spec: s, Url: location}, err
}

// FromGlob creates SpecInfo specs from local files matching the specified glob parameter
func FromGlob(loader Loader, glob string) ([]*SpecInfo, error) {
	files, err := filepathx.Glob(glob)
	if err != nil {
		return nil, err
	}
	result := make([]*SpecInfo, 0)
	for _, file := range files {
		spec, err := loader.LoadFromFile(file)
		if err != nil {
			return nil, err
		}
		result = append(result, &SpecInfo{Url: file, Spec: spec})
	}

	if len(result) > 0 {
		return result, nil
	}

	if isUrl(glob) {
		return nil, errors.New("no matching files (should be a glob, not a URL)")
	}

	return nil, errors.New("no matching files")

}

func isUrl(spec string) bool {
	_, err := url.ParseRequestURI(spec)
	return err == nil
}
