package diff

import "github.com/getkin/kin-openapi/openapi3"

// HeaderDiff is a diff between header objects: https://swagger.io/specification/#header-object
type HeaderDiff struct {
	ExtensionProps  *ExtensionsDiff `json:"extensions,omitempty"`
	DescriptionDiff *ValueDiff      `json:"description,omitempty"`
	DeprecatedDiff  *ValueDiff      `json:"deprecated,omitempty"`
	RequiredDiff    *ValueDiff      `json:"required,omitempty"`
	ExampleDiff     *ValueDiff      `json:"example,omitempty"`
	// Examples
	SchemaDiff  *SchemaDiff  `json:"schema,omitempty"`
	ContentDiff *ContentDiff `json:"content,omitempty"`
}

func (headerDiff HeaderDiff) empty() bool {
	return headerDiff == HeaderDiff{}
}

func diffHeaderValues(config *Config, header1, header2 *openapi3.Header) HeaderDiff {
	result := HeaderDiff{}

	if diff := getExtensionsDiff(config, header1.ExtensionProps, header2.ExtensionProps); !diff.empty() {
		result.ExtensionProps = diff
	}

	result.DescriptionDiff = getValueDiff(header1.Description, header2.Description)
	result.DeprecatedDiff = getValueDiff(header1.Deprecated, header2.Deprecated)
	result.RequiredDiff = getValueDiff(header1.Required, header2.Required)

	if schemaDiff := getSchemaDiff(config, header1.Schema, header2.Schema); !schemaDiff.empty() {
		result.SchemaDiff = &schemaDiff
	}

	if config.IncludeExamples {
		result.ExampleDiff = getValueDiff(header1.Example, header2.Example)
	}

	if contentDiff := getContentDiff(config, header1.Content, header2.Content); !contentDiff.empty() {
		result.ContentDiff = &contentDiff
	}

	return result
}
