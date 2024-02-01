package load

import (
	"errors"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/yargevad/filepathx"
)

// SpecInfo contains information about an OpenAPI spec and its metadata
type SpecInfo struct {
	Url     string
	Spec    *openapi3.T
	Version string
}

func (specInfo *SpecInfo) GetVersion() string {
	if specInfo == nil || specInfo.Version == "" {
		return "n/a"
	}
	return specInfo.Version
}

func newSpecInfo(spec *openapi3.T, path string) *SpecInfo {
	return &SpecInfo{
		Spec:    spec,
		Url:     path,
		Version: getVersion(spec),
	}
}

func getVersion(spec *openapi3.T) string {
	if spec == nil || spec.Info == nil {
		return ""
	}

	return spec.Info.Version
}

// LoadSpecInfo creates a SpecInfo from a local file path, a URL, or stdin
func LoadSpecInfo(loader Loader, source *Source) (*SpecInfo, error) {
	s, err := from(loader, source)
	if err != nil {
		return nil, err
	}
	return newSpecInfo(s, source.Path), nil
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
