package diff

import "github.com/getkin/kin-openapi/openapi3"

// TagDiff is a diff between tag objects: https://swagger.io/specification/#tag-object
type TagDiff struct {
	NameDiff        *ValueDiff `json:"name,omitempty" yaml:"name,omitempty"`
	DescriptionDiff *ValueDiff `json:"description,omitempty" yaml:"description,omitempty"`
}

// Empty return true if there is no diff
func (diff *TagDiff) Empty() bool {
	return diff == nil || *diff == TagDiff{}
}

func getTagDiff(tag1, tag2 *openapi3.Tag) *TagDiff {
	diff := getTagDiffInternal(tag1, tag2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getTagDiffInternal(tag1, tag2 *openapi3.Tag) *TagDiff {
	result := TagDiff{}

	result.NameDiff = getValueDiff(tag1.Name, tag2.Name)
	result.DescriptionDiff = getValueDiff(tag1.Description, tag2.Description)

	return &result
}
