package diff

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
)

// SpecDiff describes the changes between two OpenAPI specifications: https://swagger.io/specification/#schema
type SpecDiff struct {
	ExtensionsDiff   *ExtensionsDiff           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	OpenAPIDiff      *ValueDiff                `json:"openAPI,omitempty" yaml:"openAPI,omitempty"`
	InfoDiff         *InfoDiff                 `json:"info,omitempty" yaml:"info,omitempty"`
	PathsDiff        *PathsDiff                `json:"paths,omitempty" yaml:"paths,omitempty"`
	SecurityDiff     *SecurityRequirementsDiff `json:"security,omitempty" yaml:"security,omitempty"`
	ServersDiff      *ServersDiff              `json:"servers,omitempty" yaml:"servers,omitempty"`
	TagsDiff         *TagsDiff                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocsDiff *ExternalDocsDiff         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	ComponentsDiff `json:"components,omitempty" yaml:"components,omitempty"`
}

func newSpecDiff() *SpecDiff {
	return &SpecDiff{}
}

// Empty indicates whether a change was found in this element
func (specDiff *SpecDiff) Empty() bool {
	return specDiff == nil || *specDiff == SpecDiff{}
}

func getDiff(config *Config, s1, s2 *openapi3.Swagger) (*SpecDiff, error) {

	if s1 == nil || s2 == nil {
		return nil, errors.New("spec is nil")
	}

	diff, err := getDiffInternal(config, s1, s2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getDiffInternal(config *Config, s1, s2 *openapi3.Swagger) (*SpecDiff, error) {

	result := newSpecDiff()
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, s1.ExtensionProps, s2.ExtensionProps)
	result.OpenAPIDiff = getValueDiff(s1.OpenAPI, s2.OpenAPI)

	result.InfoDiff, err = getInfoDiff(config, s1.Info, s2.Info)
	if err != nil {
		return nil, err
	}

	result.PathsDiff, err = getPathsDiff(config, s1.Paths, s2.Paths)
	if err != nil {
		return nil, err
	}
	result.SecurityDiff = getSecurityRequirementsDiff(config, &s1.Security, &s2.Security)
	result.ServersDiff = getServersDiff(config, &s1.Servers, &s2.Servers)
	result.TagsDiff = getTagsDiff(s1.Tags, s2.Tags)
	result.ExternalDocsDiff = getExternalDocsDiff(config, s1.ExternalDocs, s2.ExternalDocs)

	result.ComponentsDiff, err = getComponentsDiff(config, s1.Components, s2.Components)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (specDiff *SpecDiff) getSummary() *Summary {

	summary := newSummary()

	if specDiff.Empty() {
		return summary
	}

	summary.Diff = true
	summary.add(specDiff.PathsDiff, PathsDetail)
	summary.add(specDiff.SecurityDiff, SecurityDetail)
	summary.add(specDiff.ServersDiff, ServersDetail)
	summary.add(specDiff.TagsDiff, TagsDetail)
	summary.add(specDiff.SchemasDiff, SchemasDetail)
	summary.add(specDiff.ParametersDiff, ParametersDetail)
	summary.add(specDiff.HeadersDiff, HeadersDetail)
	summary.add(specDiff.RequestBodiesDiff, RequestBodiesDetail)
	summary.add(specDiff.ResponsesDiff, ResponsesDetail)
	summary.add(specDiff.SecuritySchemesDiff, SecuritySchemesDetail)
	summary.add(specDiff.CallbacksDiff, CallbacksDetail)

	return summary
}

// Apply applies the diff
func (specDiff *SpecDiff) Patch(s *openapi3.Swagger) error {

	if specDiff.Empty() {
		return nil
	}

	err := specDiff.PathsDiff.Patch(s.Paths)
	if err != nil {
		return err
	}

	return nil
}
