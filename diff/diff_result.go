package diff

import (
	"regexp"

	"github.com/apex/log"
	"github.com/getkin/kin-openapi/openapi3"
)

type DiffResult struct {
	MissingEndpoints  []string          `json:"missingEndpoints,omitempty"`
	ModifiedEndpoints ModifiedEndpoints `json:"modifiedEndpoints,omitempty"`
}

func (diffResult *DiffResult) empty() bool {
	return len(diffResult.MissingEndpoints) == 0 && len(diffResult.ModifiedEndpoints) == 0
}

func newDiffResult() *DiffResult {
	return &DiffResult{
		MissingEndpoints:  []string{},
		ModifiedEndpoints: ModifiedEndpoints{},
	}
}

func (diffResult *DiffResult) addMissingEndpoint(endpoint string) {
	diffResult.MissingEndpoints = append(diffResult.MissingEndpoints, endpoint)
}

func (diffResult *DiffResult) addModifiedEndpoint(entrypoint1 string, pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) {

	diff := diffEndpoints(pathItem1, pathItem2)
	if !diff.empty() {
		diffResult.ModifiedEndpoints.addEndpointDiff(entrypoint1, pathItem1, pathItem2)
	}
}

func (diffResult *DiffResult) FilterByRegex(filter string) {
	r, err := regexp.Compile(filter)
	if err != nil {
		log.Errorf("Failed to compile filter regex '%s' with '%v'", filter, err)
		return
	}

	diffResult.filterMissingEndpoints(r)
	diffResult.filterModifiedEndpoints(r)
}

func (diffResult *DiffResult) filterMissingEndpoints(r *regexp.Regexp) {
	result := []string{}
	for _, endpoint := range diffResult.MissingEndpoints {
		if r.MatchString(endpoint) {
			result = append(result, endpoint)
		}
	}

	diffResult.MissingEndpoints = result
}

func (diffResult *DiffResult) filterModifiedEndpoints(r *regexp.Regexp) {
	result := ModifiedEndpoints{}
	for endpoint, endpointDiff := range diffResult.ModifiedEndpoints {
		if r.MatchString(endpoint) {
			result[endpoint] = endpointDiff
		}
	}

	diffResult.ModifiedEndpoints = result
}

func (diffResult *DiffResult) GetSummary() *DiffSummary {
	return getDiffSummary(diffResult)
}
