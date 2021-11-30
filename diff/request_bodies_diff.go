package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// RequestBodiesDiff describes the changes between a pair of sets of request body objects: https://swagger.io/specification/#request-body-object
type RequestBodiesDiff struct {
	Added    StringList            `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList            `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedRequestBodies `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (requestBodiesDiff *RequestBodiesDiff) Empty() bool {
	if requestBodiesDiff == nil {
		return true
	}

	return len(requestBodiesDiff.Added) == 0 &&
		len(requestBodiesDiff.Deleted) == 0 &&
		len(requestBodiesDiff.Modified) == 0
}

// Breaking indicates whether this element includes a breaking change
func (diff *RequestBodiesDiff) Breaking() bool {
	return len(diff.Deleted) > 0 ||
		diff.Modified.Breaking()
}

// ModifiedRequestBodies is map of requestBody names to their respective diffs
type ModifiedRequestBodies map[string]*RequestBodyDiff

// Breaking indicates whether this element includes a breaking change
func (diff ModifiedRequestBodies) Breaking() bool {
	for _, modifiedRequestBody := range diff {
		if modifiedRequestBody.Breaking() {
			return true
		}
	}
	return false
}

func newRequestBodiesDiff() *RequestBodiesDiff {
	return &RequestBodiesDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedRequestBodies{},
	}
}

func getRequestBodiesDiff(config *Config, requestBodies1, requestBodies2 openapi3.RequestBodies) (*RequestBodiesDiff, error) {
	diff, err := getRequestBodiesDiffInternal(config, requestBodies1, requestBodies2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	if config.BreakingOnly && !diff.Breaking() {
		return nil, nil
	}

	return diff, nil
}

func getRequestBodiesDiffInternal(config *Config, requestBodies1, requestBodies2 openapi3.RequestBodies) (*RequestBodiesDiff, error) {

	result := newRequestBodiesDiff()

	for requestBodyValue1, requestBodyRef1 := range requestBodies1 {
		if requestBodyValue2, ok := requestBodies2[requestBodyValue1]; ok {
			diff, err := getRequestBodyDiff(config, requestBodyRef1, requestBodyValue2)
			if err != nil {
				return nil, err
			}
			if !diff.Empty() {
				result.Modified[requestBodyValue1] = diff
			}
		} else {
			result.Deleted = append(result.Deleted, requestBodyValue1)
		}
	}

	for requestBodyValue2 := range requestBodies2 {
		if _, ok := requestBodies1[requestBodyValue2]; !ok {
			result.Added = append(result.Added, requestBodyValue2)
		}
	}

	return result, nil
}

func (requestBodiesDiff *RequestBodiesDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(requestBodiesDiff.Added),
		Deleted:  len(requestBodiesDiff.Deleted),
		Modified: len(requestBodiesDiff.Modified),
	}
}
