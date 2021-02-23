package diff

import "github.com/getkin/kin-openapi/openapi3"

// SpecDiff describes the changes between two OAS specs
type SpecDiff struct {
	PathsDiff      *PathsDiff      `json:"paths,omitempty"`      // deep diff of paths including their schemas, parameters, responses etc.
	TagsDiff       *TagsDiff       `json:"tags,omitempty"`       // diff of tags
	SchemasDiff    *SchemasDiff    `json:"schemas,omitempty"`    // diff of top-level schemas (under components)
	ParametersDiff *ParametersDiff `json:"parameters,omitempty"` // diff of top-level parameters (under components)
	HeadersDiff    *HeadersDiff    `json:"headers,omitempty"`    // diff of top-level headers (under components)
	ResponsesDiff  *ResponsesDiff  `json:"responses,omitempty"`  // diff of top-level responses (under components)
	CallbacksDiff  *CallbacksDiff  `json:"callbacks,omitempty"`  // diff of top-level callbacks (under components)
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
	diff.setTagsDiff(getTagsDiff(s1.Tags, s2.Tags))

	// components
	diff.setSchemasDiff(getSchemasDiff(s1.Components.Schemas, s2.Components.Schemas))
	diff.setParametersDiff(getParametersDiff(toParameters(s1.Components.Parameters), toParameters(s2.Components.Parameters)))
	diff.setHeadersDiff(getHeadersDiff(s1.Components.Headers, s2.Components.Headers))
	diff.setResponsesDiff(getResponsesDiff(s1.Components.Responses, s2.Components.Responses))
	diff.setCallbacksDiff(getCallbacksDiff(s1.Components.Callbacks, s2.Components.Callbacks))

	return diff
}

func (specDiff *SpecDiff) setPathsDiff(pathsDiff *PathsDiff) {
	specDiff.PathsDiff = nil

	if !pathsDiff.empty() {
		specDiff.PathsDiff = pathsDiff
	}
}

func (specDiff *SpecDiff) setTagsDiff(tagsDiff *TagsDiff) {
	specDiff.TagsDiff = nil

	if !tagsDiff.empty() {
		specDiff.TagsDiff = tagsDiff
	}
}

func (specDiff *SpecDiff) setSchemasDiff(schemasDiff *SchemasDiff) {
	specDiff.SchemasDiff = nil

	if !schemasDiff.empty() {
		specDiff.SchemasDiff = schemasDiff
	}
}

func (specDiff *SpecDiff) setParametersDiff(parametersDiff *ParametersDiff) {
	specDiff.ParametersDiff = nil

	if !parametersDiff.empty() {
		specDiff.ParametersDiff = parametersDiff
	}
}

func (specDiff *SpecDiff) setHeadersDiff(headersDiff *HeadersDiff) {
	specDiff.HeadersDiff = nil

	if !headersDiff.empty() {
		specDiff.HeadersDiff = headersDiff
	}
}

func (specDiff *SpecDiff) setResponsesDiff(responsesDiff *ResponsesDiff) {
	specDiff.ResponsesDiff = nil

	if !responsesDiff.empty() {
		specDiff.ResponsesDiff = responsesDiff
	}
}

func (specDiff *SpecDiff) setCallbacksDiff(callbacksDiff *CallbacksDiff) {
	specDiff.CallbacksDiff = nil

	if !callbacksDiff.empty() {
		specDiff.CallbacksDiff = callbacksDiff
	}
}

func (specDiff *SpecDiff) getSummary() *Summary {

	summary := newSummary()

	summary.Diff = !specDiff.empty()
	summary.add(specDiff.PathsDiff, "paths")
	summary.add(specDiff.TagsDiff, "tags")
	summary.add(specDiff.SchemasDiff, "schemas")
	summary.add(specDiff.ParametersDiff, "parameters")
	summary.add(specDiff.HeadersDiff, "headers")
	summary.add(specDiff.ResponsesDiff, "responses")
	summary.add(specDiff.CallbacksDiff, "callbacks")

	return summary
}

func (specDiff *SpecDiff) filterByRegex(filter string) {
	if specDiff.PathsDiff != nil {

		specDiff.setPathsDiff(specDiff.PathsDiff.filterByRegex(filter))
	}
}
