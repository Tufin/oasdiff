package diff

import "github.com/getkin/kin-openapi/openapi3"

type DiffResponse struct {
	DiffResult  *DiffResult  `json:"diffResult,omitempty"`
	DiffSummary *DiffSummary `json:"diffSummary,omitempty"`
}

// GetDiffResponse returns the diff between two OAS (swagger) specs including a diff summary
func GetDiffResponse(s1 *openapi3.Swagger, s2 *openapi3.Swagger, prefix string, filter string) DiffResponse {
	diffResult := Diff(s1, s2, prefix)
	diffResult.FilterByRegex(filter)

	return DiffResponse{
		DiffResult:  diffResult,
		DiffSummary: diffResult.getSummary(),
	}
}
