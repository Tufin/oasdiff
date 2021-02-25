package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// RequestBodiesDiff is a diff between two sets of request body objects: https://swagger.io/specification/#request-body-object
type RequestBodiesDiff struct {
	Added    StringList            `json:"added,omitempty"`
	Deleted  StringList            `json:"deleted,omitempty"`
	Modified ModifiedRequestBodies `json:"modified,omitempty"`
}

func (requestBodiesDiff *RequestBodiesDiff) empty() bool {
	return len(requestBodiesDiff.Added) == 0 &&
		len(requestBodiesDiff.Deleted) == 0 &&
		len(requestBodiesDiff.Modified) == 0
}

// ModifiedRequestBodies is map of requestBody names to their respective diffs
type ModifiedRequestBodies map[string]*RequestBodyDiff

func newRequestBodiesDiff() *RequestBodiesDiff {
	return &RequestBodiesDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedRequestBodies{},
	}
}

func getRequestBodiesDiff(requestBodies1, requestBodies2 openapi3.RequestBodies) *RequestBodiesDiff {

	result := newRequestBodiesDiff()

	for requestBodyValue1, requestBodyRef1 := range requestBodies1 {
		if requestBodyRef1 != nil && requestBodyRef1.Value != nil {
			if requestBodyValue2, ok := requestBodies2[requestBodyValue1]; ok {
				if diff := getRequestBodyDiff(requestBodyRef1, requestBodyValue2); !diff.empty() {
					result.Modified[requestBodyValue1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, requestBodyValue1)
			}
		}
	}

	for requestBodyValue2, requestBodyRef2 := range requestBodies2 {
		if requestBodyRef2 != nil && requestBodyRef2.Value != nil {
			if _, ok := requestBodies1[requestBodyValue2]; !ok {
				result.Added = append(result.Added, requestBodyValue2)
			}
		}
	}

	return result
}

func (requestBodiesDiff *RequestBodiesDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(requestBodiesDiff.Added),
		Deleted:  len(requestBodiesDiff.Deleted),
		Modified: len(requestBodiesDiff.Modified),
	}
}
