package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ResponsesDiff describes the changes between a pair of sets of response objects: https://swagger.io/specification/#responses-object
type ResponsesDiff struct {
	Added    StringList        `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList        `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedResponses `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (responsesDiff *ResponsesDiff) Empty() bool {
	if responsesDiff == nil {
		return true
	}

	return len(responsesDiff.Added) == 0 &&
		len(responsesDiff.Deleted) == 0 &&
		len(responsesDiff.Modified) == 0
}

func (responsesDiff *ResponsesDiff) removeNonBreaking() {

	if responsesDiff.Empty() {
		return
	}

	responsesDiff.Added = nil
}

// ModifiedResponses is map of response values to their respective diffs
type ModifiedResponses map[string]*ResponseDiff

func newResponsesDiff() *ResponsesDiff {
	return &ResponsesDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedResponses{},
	}
}

func getResponsesDiff(config *Config, state *state, responses1, responses2 openapi3.Responses) (*ResponsesDiff, error) {
	diff, err := getResponsesDiffInternal(config, state, responses1, responses2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getResponsesDiffInternal(config *Config, state *state, responses1, responses2 openapi3.Responses) (*ResponsesDiff, error) {

	result := newResponsesDiff()

	for responseValue1, responseRef1 := range responses1 {
		if responseRef2, ok := responses2[responseValue1]; ok {
			value1, err := derefResponse(responseRef1)
			if err != nil {
				return nil, err
			}

			value2, err := derefResponse(responseRef2)
			if err != nil {
				return nil, err
			}

			diff, err := diffResponseValues(config, state, value1, value2)
			if err != nil {
				return nil, err
			}
			if !diff.Empty() {
				result.Modified[responseValue1] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, responseValue1)
		}
	}

	for responseValue2 := range responses2 {
		if _, ok := responses1[responseValue2]; !ok {
			result.Added = append(result.Added, responseValue2)
		}
	}

	return result, nil
}

func derefResponse(ref *openapi3.ResponseRef) (*openapi3.Response, error) {

	if ref == nil || ref.Value == nil {
		return nil, fmt.Errorf("response reference is nil")
	}

	return ref.Value, nil
}

func (responsesDiff *ResponsesDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(responsesDiff.Added),
		Deleted:  len(responsesDiff.Deleted),
		Modified: len(responsesDiff.Modified),
	}
}
