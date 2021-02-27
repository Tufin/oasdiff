package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExtensionsDiff is a diff between two sets of specification extensions: https://swagger.io/specification/#specification-extensions
type ExtensionsDiff struct {
	Added    StringList         `json:"added,omitempty"`
	Deleted  StringList         `json:"deleted,omitempty"`
	Modified ModifiedExtensions `json:"modified,omitempty"`
}

type ModifiedExtensions map[string]*ValueDiff

func (diff *ExtensionsDiff) empty() bool {
	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

func newExtensionsDiff() *ExtensionsDiff {
	return &ExtensionsDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedExtensions{},
	}
}

func getExtensionsDiff(extensions1, extensions2 openapi3.ExtensionProps) *ExtensionsDiff {

	result := newExtensionsDiff()

	for name1, extension1 := range extensions1.Extensions {
		if extension2, ok := extensions2.Extensions[name1]; ok {
			if diff := getValueDiff(extension1, extension2); diff != nil {
				result.Modified[name1] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, name1)
		}
	}

	for name2 := range extensions2.Extensions {
		if _, ok := extensions1.Extensions[name2]; !ok {
			result.Added = append(result.Deleted, name2)
		}
	}

	return result
}
