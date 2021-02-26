package diff

import "github.com/getkin/kin-openapi/openapi3"

// SpecDiff describes the changes between two OpenAPI specifications: https://swagger.io/specification/#specification
type SpecDiff struct {
	PathsDiff   *PathsDiff   `json:"paths,omitempty"`
	ServersDiff *ServersDiff `json:"servers,omitempty"`
	TagsDiff    *TagsDiff    `json:"tags,omitempty"`

	// Components
	SchemasDiff       *SchemasDiff       `json:"schemas,omitempty"`
	ParametersDiff    *ParametersDiff    `json:"parameters,omitempty"`
	HeadersDiff       *HeadersDiff       `json:"headers,omitempty"`
	RequestBodiesDiff *RequestBodiesDiff `json:"requestBodies,omitempty"`
	ResponsesDiff     *ResponsesDiff     `json:"responses,omitempty"`
	CallbacksDiff     *CallbacksDiff     `json:"callbacks,omitempty"`
}

func newSpecDiff() *SpecDiff {
	return &SpecDiff{}
}

func (specDiff SpecDiff) empty() bool {
	return specDiff == SpecDiff{}
}

func getDiff(s1, s2 *openapi3.Swagger, prefix string) *SpecDiff {

	diff := newSpecDiff()

	diff.setPathsDiff(getPathsDiff(s1.Paths, s2.Paths, prefix))
	diff.setServersDiff(getServersDiff(&s1.Servers, &s2.Servers))
	diff.setTagsDiff(getTagsDiff(s1.Tags, s2.Tags))

	// components
	diff.setSchemasDiff(getSchemasDiff(s1.Components.Schemas, s2.Components.Schemas))
	diff.setParametersDiff(getParametersDiff(toParameters(s1.Components.Parameters), toParameters(s2.Components.Parameters)))
	diff.setHeadersDiff(getHeadersDiff(s1.Components.Headers, s2.Components.Headers))
	diff.setRequestBodiesDiff(getRequestBodiesDiff(s1.Components.RequestBodies, s2.Components.RequestBodies))
	diff.setResponsesDiff(getResponsesDiff(s1.Components.Responses, s2.Components.Responses))
	diff.setCallbacksDiff(getCallbacksDiff(s1.Components.Callbacks, s2.Components.Callbacks))

	return diff
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

func (specDiff *SpecDiff) setSchemasDiff(diff *SchemasDiff) {
	specDiff.SchemasDiff = nil

	if !diff.empty() {
		specDiff.SchemasDiff = diff
	}
}

func (specDiff *SpecDiff) setParametersDiff(diff *ParametersDiff) {
	specDiff.ParametersDiff = nil

	if !diff.empty() {
		specDiff.ParametersDiff = diff
	}
}

func (specDiff *SpecDiff) setHeadersDiff(diff *HeadersDiff) {
	specDiff.HeadersDiff = nil

	if !diff.empty() {
		specDiff.HeadersDiff = diff
	}
}

func (specDiff *SpecDiff) setRequestBodiesDiff(diff *RequestBodiesDiff) {
	specDiff.RequestBodiesDiff = nil

	if !diff.empty() {
		specDiff.RequestBodiesDiff = diff
	}
}

func (specDiff *SpecDiff) setResponsesDiff(diff *ResponsesDiff) {
	specDiff.ResponsesDiff = nil

	if !diff.empty() {
		specDiff.ResponsesDiff = diff
	}
}

func (specDiff *SpecDiff) setCallbacksDiff(diff *CallbacksDiff) {
	specDiff.CallbacksDiff = nil

	if !diff.empty() {
		specDiff.CallbacksDiff = diff
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
	summary.add(specDiff.CallbacksDiff, CallbacksComponent)

	return summary
}

func (specDiff *SpecDiff) filterByRegex(filter string) {
	if specDiff.PathsDiff != nil {

		specDiff.setPathsDiff(specDiff.PathsDiff.filterByRegex(filter))
	}
}
