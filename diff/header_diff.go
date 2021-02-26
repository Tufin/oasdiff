package diff

import "github.com/getkin/kin-openapi/openapi3"

// HeaderDiff is a diff between header objects: https://swagger.io/specification/#header-object
type HeaderDiff struct {

	// TODO: diff ExtensionProps
	DescriptionDiff *ValueDiff `json:"description,omitempty"`
	DeprecatedDiff  *ValueDiff `json:"deprecated,omitempty"`
	RequiredDiff    *ValueDiff `json:"required,omitempty"`
	ExampleDiff     *ValueDiff `json:"example,omitempty"`
	// TODO: diff Examples
	SchemaDiff  *SchemaDiff  `json:"schema,omitempty"`
	ContentDiff *ContentDiff `json:"content,omitempty"`
}

func (headerDiff HeaderDiff) empty() bool {
	return headerDiff == HeaderDiff{}
}

func diffHeaderValues(header1, header2 *openapi3.Header) HeaderDiff {
	result := HeaderDiff{}

	result.DescriptionDiff = getValueDiff(header1.Description, header2.Description)
	result.DeprecatedDiff = getValueDiff(header1.Deprecated, header2.Deprecated)
	result.RequiredDiff = getValueDiff(header1.Required, header2.Required)

	if schemaDiff := getSchemaDiff(header1.Schema, header2.Schema); !schemaDiff.empty() {
		result.SchemaDiff = &schemaDiff
	}

	result.ExampleDiff = getValueDiff(header1.Example, header2.Example)

	if contentDiff := getContentDiff(header1.Content, header2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
