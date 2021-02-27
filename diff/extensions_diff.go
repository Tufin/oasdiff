package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExtensionsDiff is a diff between two sets of specification extensions: https://swagger.io/specification/#specification-extensions
type ExtensionsDiff struct {
	Added   StringList `json:"added,omitempty"`
	Deleted StringList `json:"deleted,omitempty"`
	// Modified ModifiedPaths `json:"modified,omitempty"`
}

func (diff *ExtensionsDiff) empty() bool {
	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 // && len(diff.Modified) == 0
}

func newExtensionsDiff() *ExtensionsDiff {
	return &ExtensionsDiff{
		Added:   StringList{},
		Deleted: StringList{},
		// Modified: ModifiedServers{},
	}
}

func getExtensionsDiff(extensions1, extensions2 openapi3.ExtensionProps) *ExtensionsDiff {

	diff := newExtensionsDiff()

	for name1 := range extensions1.Extensions {
		if _, ok := extensions2.Extensions[name1]; ok {
		} else {
			diff.Deleted = append(diff.Deleted, name1)
		}
	}

	for name2 := range extensions2.Extensions {
		if _, ok := extensions1.Extensions[name2]; ok {
		} else {
			diff.Added = append(diff.Deleted, name2)
		}
	}

	return diff
}
