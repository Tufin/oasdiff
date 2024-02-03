package load

import (
	"errors"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/flatten/allof"
	"github.com/tufin/oasdiff/flatten/commonparams"
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

type Option func(Loader, []*SpecInfo) ([]*SpecInfo, error)

// WithIdentity returns the original SpecInfos
func WithIdentity() Option {
	return func(loader Loader, specInfos []*SpecInfo) ([]*SpecInfo, error) {
		return specInfos, nil
	}
}

func GetOption(option Option, enable bool) Option {
	if !enable {
		return WithIdentity()
	}
	return option
}

// WithFlattenAllOf returns SpecInfos with flattened allOf
func WithFlattenAllOf() Option {
	return func(loader Loader, specInfos []*SpecInfo) ([]*SpecInfo, error) {
		var err error
		for _, specInfo := range specInfos {
			if specInfo.Spec, err = allof.MergeSpec(specInfo.Spec); err != nil {
				return nil, err
			}
		}
		return specInfos, nil
	}
}

// WithFlattenParams returns SpecInfos with Common Parameters combined into operation parameters
// See here for Common Parameters definition: https://swagger.io/docs/specification/describing-parameters/
func WithFlattenParams() Option {
	return func(loader Loader, specInfos []*SpecInfo) ([]*SpecInfo, error) {
		for _, specInfo := range specInfos {
			commonparams.Move(specInfo.Spec)
		}
		return specInfos, nil
	}
}

// NewSpecInfo creates a SpecInfo from a local file path, a URL, or stdin
func NewSpecInfo(loader Loader, source *Source, options ...Option) (*SpecInfo, error) {
	specInfo, err := loadSpecInfo(loader, source)
	if err != nil {
		return nil, err
	}
	specInfos := []*SpecInfo{specInfo}

	for _, option := range options {
		if specInfos, err = option(loader, specInfos); err != nil {
			return nil, err
		}
	}
	return specInfos[0], nil
}

// NewSpecInfoFromGlob creates SpecInfos from local files matching the specified glob parameter
func NewSpecInfoFromGlob(loader Loader, glob string, options ...Option) ([]*SpecInfo, error) {
	specInfos, err := fromGlob(loader, glob)
	if err != nil {
		return nil, err
	}

	for _, option := range options {
		if specInfos, err = option(loader, specInfos); err != nil {
			return nil, err
		}
	}
	return specInfos, nil
}

func loadSpecInfo(loader Loader, source *Source) (*SpecInfo, error) {
	s, err := from(loader, source)
	if err != nil {
		return nil, err
	}
	return newSpecInfo(s, source.Path), nil
}

func fromGlob(loader Loader, glob string) ([]*SpecInfo, error) {
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
