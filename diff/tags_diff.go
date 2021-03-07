package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// TagsDiff is a diff between two lists of tag objects: https://swagger.io/specification/#tag-object
type TagsDiff struct {
	Added    StringList   `json:"added,omitempty"`
	Deleted  StringList   `json:"deleted,omitempty"`
	Modified ModifiedTags `json:"modified,omitempty"`
}

func newTagsDiff() *TagsDiff {
	return &TagsDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedTags{},
	}
}

// ModifiedTags is map of tag names to their respective diffs
type ModifiedTags map[string]TagDiff

func (tagsDiff *TagsDiff) empty() bool {
	if tagsDiff == nil {
		return true
	}

	return len(tagsDiff.Added) == 0 &&
		len(tagsDiff.Deleted) == 0 &&
		len(tagsDiff.Modified) == 0
}

func getTagsDiff(tags1, tags2 openapi3.Tags) *TagsDiff {

	result := newTagsDiff()

	for _, tag1 := range tags1 {
		if tag2 := tags2.Get(tag1.Name); tag2 != nil {
			if diff := getTagDiff(tag1, tag2); !diff.empty() {
				result.Modified[tag1.Name] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, tag1.Name)
		}
	}

	for _, tag2 := range tags2 {
		if tag1 := tags1.Get(tag2.Name); tag1 == nil {
			result.Added = append(result.Added, tag2.Name)
		}
	}

	if result.empty() {
		return nil
	}

	return result
}

func (tagsDiff *TagsDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(tagsDiff.Added),
		Deleted:  len(tagsDiff.Deleted),
		Modified: len(tagsDiff.Modified),
	}
}
