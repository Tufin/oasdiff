package diff

import "github.com/getkin/kin-openapi/openapi3"

// SpecDiff describes the changes between two OpenAPI specifications: https://swagger.io/specification/#schema
type SpecDiff struct {
	ExtensionsDiff *ExtensionsDiff `json:"extensions,omitempty"`
	OpenAPIDiff    *ValueDiff      `json:"openAPI,omitempty"`
	PathsDiff      *PathsDiff      `json:"paths,omitempty"`
	ServersDiff    *ServersDiff    `json:"servers,omitempty"`
	TagsDiff       *TagsDiff       `json:"tags,omitempty"`

	ComponentsDiff
}

func newSpecDiff() *SpecDiff {
	return &SpecDiff{}
}

func (specDiff SpecDiff) empty() bool {
	return specDiff == SpecDiff{}
}

func getDiff(config *Config, s1, s2 *openapi3.Swagger) *SpecDiff {

	result := newSpecDiff()

	result.ExtensionsDiff = getExtensionsDiff(config, s1.ExtensionProps, s2.ExtensionProps)
	result.OpenAPIDiff = getValueDiff(s1.OpenAPI, s2.OpenAPI)
	// Info
	result.setPathsDiff(getPathsDiff(config, s1.Paths, s2.Paths))
	// Security
	result.setServersDiff(getServersDiff(config, &s1.Servers, &s2.Servers))
	result.setTagsDiff(getTagsDiff(s1.Tags, s2.Tags))
	// ExternalDocs

	result.ComponentsDiff = getComponentsDiff(config, s1.Components, s2.Components)

	return result
}

func (specDiff *SpecDiff) setPathsDiff(diff *PathsDiff) {
	specDiff.PathsDiff = nil

	if !diff.empty() {
		specDiff.PathsDiff = diff
	}
}

func (specDiff *SpecDiff) setServersDiff(diff *ServersDiff) {
	specDiff.ServersDiff = nil

	if !diff.empty() {
		specDiff.ServersDiff = diff
	}
}

func (specDiff *SpecDiff) setTagsDiff(diff *TagsDiff) {
	specDiff.TagsDiff = nil

	if !diff.empty() {
		specDiff.TagsDiff = diff
	}
}

func (specDiff *SpecDiff) getSummary() *Summary {

	summary := newSummary()

	summary.Diff = !specDiff.empty()
	summary.add(specDiff.PathsDiff, PathsComponent)
	summary.add(specDiff.ServersDiff, ServersComponent)
	summary.add(specDiff.TagsDiff, TagsComponent)
	summary.add(specDiff.SchemasDiff, SchemasComponent)
	summary.add(specDiff.ParametersDiff, ParametersComponent)
	summary.add(specDiff.HeadersDiff, HeadersComponent)
	summary.add(specDiff.RequestBodiesDiff, RequestBodiesComponent)
	summary.add(specDiff.ResponsesDiff, ResponsesComponent)
	summary.add(specDiff.SecuritySchemesDiff, SecuritySchemesComponent)
	summary.add(specDiff.CallbacksDiff, CallbacksComponent)

	return summary
}

func (specDiff *SpecDiff) filterByRegex(filter string) {
	if specDiff.PathsDiff != nil {

		specDiff.setPathsDiff(specDiff.PathsDiff.filterByRegex(filter))
	}
}
