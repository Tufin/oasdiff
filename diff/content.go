package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ContentDiff is the diff between two OAS contents
type ContentDiff struct {
	MediaTypeAdded   bool `json:"mediaTypeAdded,omitempty"`
	MediaTypeDeleted bool `json:"mediaTypeDeleted,omitempty"`
	MediaTypeDiff    bool `json:"mediaTypeDiff,omitempty"`

	// TODO: ExtensionProps

	SchemaDiff *SchemaDiff `json:"schemaDiff,omitempty"`
	//Encoding   map[string]*Encoding `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

func newContentDiff() ContentDiff {
	return ContentDiff{}
}

func (contentDiff ContentDiff) empty() bool {
	return contentDiff == newContentDiff()
}

func diffContent(content1 openapi3.Content, content2 openapi3.Content) ContentDiff {

	result := newContentDiff()

	if len(content1) == 0 && len(content2) == 0 {
		return result
	}

	if len(content1) == 0 && len(content2) != 0 {
		result.MediaTypeAdded = true
		return result
	}

	if len(content1) != 0 && len(content2) == 0 {
		result.MediaTypeDeleted = true
		return result
	}

	mediaType1, mediaTypeValue1, err := getSingleEntry(content1)
	if err != nil {
		return result
	}

	mediaType2, mediaTypeValue2, err := getSingleEntry(content2)
	if err != nil {
		return result
	}

	if mediaType1 != mediaType2 {
		result.MediaTypeDiff = true
		return result
	}

	if schemaDiff := diffSchema(mediaTypeValue1.Schema, mediaTypeValue2.Schema); !schemaDiff.empty() {
		result.SchemaDiff = &schemaDiff
	}

	return result
}

// getSingleEntry returns the single entry in the Content map
func getSingleEntry(content openapi3.Content) (string, *openapi3.MediaType, error) {

	var mediaType string
	var mediaTypeValue *openapi3.MediaType

	if len(content) != 1 {
		return "", nil, fmt.Errorf("content map has more than one value - this shouldn't happen. %+v", content)
	}

	for k, v := range content {
		mediaType = k
		mediaTypeValue = v
	}

	return mediaType, mediaTypeValue, nil
}
