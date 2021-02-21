package diff

import "github.com/getkin/kin-openapi/openapi3"

type Responses struct {
	Added    map[string]ResponseList      `json:"added,omitempty"`
	Deleted  map[string]ResponseList      `json:"deleted,omitempty"`
	Modified map[string]ModifiedResponses `json:"modified,omitempty"`
}

func (responses *Responses) empty() bool {
	return len(responses.Added) == 0 &&
		len(responses.Deleted) == 0 &&
		len(responses.Modified) == 0
}

// ResponseList is a set of response values
type ResponseList map[string]struct{}

// ModifiedResponses is map of response value to their respective diffs
type ModifiedResponses map[string]ResponseDiff

type ResponseDiff struct {
}

func newResponses() *Responses {
	return &Responses{}
}

func getResponseDiff(responses1, responses2 openapi3.Responses) *Responses {
	result := newResponses()

	return result
}
