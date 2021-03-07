package diff

import "github.com/getkin/kin-openapi/openapi3"

// EncodingsDiff is a diff between two sets of encoding objects: https://swagger.io/specification/#encoding-object
type EncodingsDiff struct {
	Added    StringList        `json:"added,omitempty"`
	Deleted  StringList        `json:"deleted,omitempty"`
	Modified ModifiedEncodings `json:"modified,omitempty"`
}

// ModifiedEncodings is map of enconding names to their respective diffs
type ModifiedEncodings map[string]*EncodingDiff

func (diff *EncodingsDiff) empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

func newEncodingsDiff() *EncodingsDiff {
	return &EncodingsDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedEncodings{},
	}
}

func getEncodingsDiff(config *Config, encodings1, encodings2 map[string]*openapi3.Encoding) *EncodingsDiff {
	diff := getEncodingsDiffInternal(config, encodings1, encodings2)
	if diff.empty() {
		return nil
	}
	return diff
}

func getEncodingsDiffInternal(config *Config, encodings1, encodings2 map[string]*openapi3.Encoding) *EncodingsDiff {

	result := newEncodingsDiff()

	for name1, encoding1 := range encodings1 {
		if encoding2, ok := encodings2[name1]; ok {
			if diff := getEncodingDiff(config, encoding1, encoding2); !diff.empty() {
				result.Modified[name1] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, name1)
		}
	}

	for name2 := range encodings2 {
		if _, ok := encodings1[name2]; !ok {
			result.Added = append(result.Added, name2)
		}
	}

	return result
}
