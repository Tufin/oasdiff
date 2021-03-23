package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ExtensionsDiff describes the changes between a pair of sets of specification extensions: https://swagger.io/specification/#specification-extensions
type ExtensionsDiff struct {
	Added    StringList         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedExtensions `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// ModifiedExtensions is map of extensions names to their respective diffs
type ModifiedExtensions map[string]*ValueDiff

// Empty indicates whether a change was found in this element
func (diff *ExtensionsDiff) Empty() bool {
	if diff == nil {
		return true
	}

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

func getExtensionsDiff(config *Config, extensions1, extensions2 openapi3.ExtensionProps) *ExtensionsDiff {
	diff := getExtensionsDiffInternal(config, extensions1, extensions2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getExtensionsDiffInternal(config *Config, extensions1, extensions2 openapi3.ExtensionProps) *ExtensionsDiff {

	result := newExtensionsDiff()

	for name1, extension1 := range extensions1.Extensions {
		if _, ok := config.IncludeExtensions[name1]; ok {
			if extension2, ok := extensions2.Extensions[name1]; ok {
				if diff := getValueDiff(extension1, extension2); !diff.Empty() {
					result.Modified[name1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, name1)
			}
		}
	}

	for name2 := range extensions2.Extensions {
		if _, ok := config.IncludeExtensions[name2]; ok {
			if _, ok := extensions1.Extensions[name2]; !ok {
				result.Added = append(result.Deleted, name2)
			}
		}
	}

	return result
}
