package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// HeadersDiff is a diff between two sets of header objects: https://swagger.io/specification/#header-object
type HeadersDiff struct {
	Added    StringList      `json:"added,omitempty"`
	Deleted  StringList      `json:"deleted,omitempty"`
	Modified ModifiedHeaders `json:"modified,omitempty"`
}

func (headersDiff *HeadersDiff) empty() bool {
	return len(headersDiff.Added) == 0 &&
		len(headersDiff.Deleted) == 0 &&
		len(headersDiff.Modified) == 0
}

// ModifiedHeaders is map of header names to their respective diffs
type ModifiedHeaders map[string]HeaderDiff

func newHeadersDiff() *HeadersDiff {
	return &HeadersDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedHeaders{},
	}
}

func getHeadersDiff(headers1, headers2 openapi3.Headers) *HeadersDiff {

	result := newHeadersDiff()

	for headerValue1, headerRef1 := range headers1 {
		if headerRef1 != nil && headerRef1.Value != nil {
			if headerValue2, ok := headers2[headerValue1]; ok {
				if diff := diffHeaderValues(headerRef1.Value, headerValue2.Value); !diff.empty() {
					result.Modified[headerValue1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, headerValue1)
			}
		}
	}

	for headerValue2, headerRef2 := range headers2 {
		if headerRef2 != nil && headerRef2.Value != nil {
			if _, ok := headers1[headerValue2]; !ok {
				result.Added = append(result.Added, headerValue2)
			}
		}
	}

	return result
}

func (headersDiff *HeadersDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(headersDiff.Added),
		Deleted:  len(headersDiff.Deleted),
		Modified: len(headersDiff.Modified),
	}
}
