package diff

import (
	"errors"

	"github.com/getkin/kin-openapi/openapi3"
)

// Diff describes the changes between a pair of OpenAPI specifications: https://swagger.io/specification/#schema
type Diff struct {
	ExtensionsDiff   *ExtensionsDiff           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	OpenAPIDiff      *ValueDiff                `json:"openAPI,omitempty" yaml:"openAPI,omitempty"`
	InfoDiff         *InfoDiff                 `json:"info,omitempty" yaml:"info,omitempty"`
	PathsDiff        *PathsDiff                `json:"paths,omitempty" yaml:"paths,omitempty"`
	EndpointsDiff    *EndpointsDiff            `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	SecurityDiff     *SecurityRequirementsDiff `json:"security,omitempty" yaml:"security,omitempty"`
	ServersDiff      *ServersDiff              `json:"servers,omitempty" yaml:"servers,omitempty"`
	TagsDiff         *TagsDiff                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocsDiff *ExternalDocsDiff         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	ComponentsDiff `json:"components,omitempty" yaml:"components,omitempty"`
}

func newDiff() *Diff {
	return &Diff{}
}

// Empty indicates whether a change was found in this element
func (diff *Diff) Empty() bool {
	return diff == nil || *diff == Diff{}
}

func (diff *Diff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.ExtensionsDiff = nil
	diff.OpenAPIDiff = nil
	diff.InfoDiff = nil
	diff.TagsDiff = nil
	diff.ExternalDocsDiff = nil

	diff.ComponentsDiff.removeNonBreaking()
}

/*
Get calculates the diff between a pair of OpenAPI specifications.

Note that Get expects OpenAPI References (https://swagger.io/docs/specification/using-ref/) to be resolved.
References are normally resolved automatically when you load the spec.
In other cases you can resolve refs using https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#SwaggerLoader.ResolveRefsIn.
*/
func Get(config *Config, s1, s2 *openapi3.T) (*Diff, error) {
	diff, err := getDiff(config, s1, s2)
	if err != nil {
		return nil, err
	}

	return diff, nil
}

func getDiff(config *Config, s1, s2 *openapi3.T) (*Diff, error) {

	if s1 == nil || s2 == nil {
		return nil, errors.New("spec is nil")
	}

	diff, err := getDiffInternal(config, s1, s2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getDiffInternal(config *Config, s1, s2 *openapi3.T) (*Diff, error) {

	result := newDiff()
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

	result.EndpointsDiff, err = getEndpointsDiff(config, s1.Paths, s2.Paths)
	if err != nil {
		return nil, err
	}

	result.SecurityDiff = getSecurityRequirementsDiff(config, &s1.Security, &s2.Security)
	result.ServersDiff = getServersDiff(config, &s1.Servers, &s2.Servers)
	result.TagsDiff = getTagsDiff(config, s1.Tags, s2.Tags)
	result.ExternalDocsDiff = getExternalDocsDiff(config, s1.ExternalDocs, s2.ExternalDocs)

	result.ComponentsDiff, err = getComponentsDiff(config, s1.Components, s2.Components)
	if err != nil {
		return nil, err
	}

	return result, nil
}

// GetSummary returns a summary of the changes in the diff
func (diff *Diff) GetSummary() *Summary {

	summary := newSummary()

	if diff.Empty() {
		return summary
	}

	summary.Diff = true

	// swagger
	summary.add(diff.PathsDiff, PathsDetail)
	summary.add(diff.SecurityDiff, SecurityDetail)
	summary.add(diff.ServersDiff, ServersDetail)
	summary.add(diff.TagsDiff, TagsDetail)

	// components
	summary.add(diff.SchemasDiff, SchemasDetail)
	summary.add(diff.ParametersDiff, ParametersDetail)
	summary.add(diff.HeadersDiff, HeadersDetail)
	summary.add(diff.RequestBodiesDiff, RequestBodiesDetail)
	summary.add(diff.ResponsesDiff, ResponsesDetail)
	summary.add(diff.SecuritySchemesDiff, SecuritySchemesDetail)
	summary.add(diff.ExamplesDiff, ExamplesDetail)
	summary.add(diff.LinksDiff, LinksDetail)
	summary.add(diff.CallbacksDiff, CallbacksDetail)

	// special
	summary.add(diff.EndpointsDiff, EndpointsDetail)
	return summary
}

// Patch applies the patch to a spec
func (diff *Diff) Patch(s *openapi3.T) error {

	if diff.Empty() {
		return nil
	}

	err := diff.PathsDiff.Patch(s.Paths)
	if err != nil {
		return err
	}

	return nil
}
