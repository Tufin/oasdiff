package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ResponsesDiff is a diff between two sets of responses
type ResponsesDiff struct {
	Added    ResponseList      `json:"added,omitempty"`
	Deleted  ResponseList      `json:"deleted,omitempty"`
	Modified ModifiedResponses `json:"modified,omitempty"`
}

func (responsesDiff *ResponsesDiff) empty() bool {
	return len(responsesDiff.Added) == 0 &&
		len(responsesDiff.Deleted) == 0 &&
		len(responsesDiff.Modified) == 0
}

// ResponseList is a list of response values
type ResponseList []string

// ModifiedResponses is map of response value to their respective diffs
type ModifiedResponses map[string]ResponseDiff

func newResponsesDiff() *ResponsesDiff {
	return &ResponsesDiff{}
}

func getResponseDiff(responses1, responses2 openapi3.Responses) *ResponsesDiff {

	result := newResponsesDiff()

	for responseValue1, responseRef1 := range responses1 {
		if responseRef1 != nil && responseRef1.Value != nil {
			if responseValue2, ok := responses2[responseValue1]; ok {
				if diff := diffResponseValues(responseRef1.Value, responseValue2.Value); !diff.empty() {
					// result.addModifiedResponse(responseRef1.Value, diff)
				}
			} else {
				result.Deleted = append(result.Deleted, responseValue1)
			}
		}
	}

	for responseValue2, responseRef2 := range responses2 {
		if responseRef2 != nil && responseRef2.Value != nil {
			if _, ok := responses1[responseValue2]; !ok {
				result.Added = append(result.Added, responseValue2)
			}
		}
	}

	return result

}
