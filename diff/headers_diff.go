package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

// HeadersDiff describes the changes between a pair of sets of header objects: https://swagger.io/specification/#header-object
type HeadersDiff struct {
	Added    utils.StringList `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  utils.StringList `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedHeaders  `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (headersDiff *HeadersDiff) Empty() bool {
	if headersDiff == nil {
		return true
	}

	return len(headersDiff.Added) == 0 &&
		len(headersDiff.Deleted) == 0 &&
		len(headersDiff.Modified) == 0
}

// ModifiedHeaders is map of header names to their respective diffs
type ModifiedHeaders map[string]*HeaderDiff

func newHeadersDiff() *HeadersDiff {
	return &HeadersDiff{
		Added:    utils.StringList{},
		Deleted:  utils.StringList{},
		Modified: ModifiedHeaders{},
	}
}

func getHeadersDiff(config *Config, state *state, headers1, headers2 openapi3.Headers) (*HeadersDiff, error) {
	diff, err := getHeadersDiffInternal(config, state, headers1, headers2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getHeadersDiffInternal(config *Config, state *state, headers1, headers2 openapi3.Headers) (*HeadersDiff, error) {

	result := newHeadersDiff()

	for headerName1, headerRef1 := range headers1 {
		if headerRef2, ok := headers2[headerName1]; ok {
			value1, err := derefHeader(headerRef1)
			if err != nil {
				return nil, err
			}

			value2, err := derefHeader(headerRef2)
			if err != nil {
				return nil, err
			}

			diff, err := getHeaderDiff(config, state, value1, value2)
			if err != nil {
				return nil, err
			}

			if !diff.Empty() {
				result.Modified[headerName1] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, headerName1)
		}
	}

	for headerName2 := range headers2 {
		if _, ok := headers1[headerName2]; !ok {
			result.Added = append(result.Added, headerName2)
		}
	}

	return result, nil
}

func derefHeader(ref *openapi3.HeaderRef) (*openapi3.Header, error) {

	if ref == nil || ref.Value == nil {
		return nil, fmt.Errorf("header reference is nil")
	}

	return ref.Value, nil
}

func (headersDiff *HeadersDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(headersDiff.Added),
		Deleted:  len(headersDiff.Deleted),
		Modified: len(headersDiff.Modified),
	}
}
