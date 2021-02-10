package diff

import (
	"regexp"

	"github.com/apex/log"
	"github.com/getkin/kin-openapi/openapi3"
)

type DiffResult struct {
	AddedEndpoints    []string          `json:"addedEndpoints,omitempty"`
	DeletedEndpoints  []string          `json:"deletedEndpoints,omitempty"`
	ModifiedEndpoints ModifiedEndpoints `json:"modifiedEndpoints,omitempty"`
}

func (diffResult *DiffResult) empty() bool {
	return len(diffResult.AddedEndpoints) == 0 &&
		len(diffResult.DeletedEndpoints) == 0 &&
		len(diffResult.ModifiedEndpoints) == 0
}

func newDiffResult() *DiffResult {
	return &DiffResult{
		AddedEndpoints:    []string{},
		DeletedEndpoints:  []string{},
		ModifiedEndpoints: ModifiedEndpoints{},
	}
}

func (diffResult *DiffResult) addAddedEndpoint(endpoint string) {
	diffResult.AddedEndpoints = append(diffResult.AddedEndpoints, endpoint)
}

func (diffResult *DiffResult) addDeletedEndpoint(endpoint string) {
	diffResult.DeletedEndpoints = append(diffResult.DeletedEndpoints, endpoint)
}

func (diffResult *DiffResult) addModifiedEndpoint(entrypoint1 string, pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) {

	diff := diffEndpoint(pathItem1, pathItem2)
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

	diffResult.AddedEndpoints = filterEndpoints(diffResult.AddedEndpoints, r)
	diffResult.DeletedEndpoints = filterEndpoints(diffResult.DeletedEndpoints, r)
	diffResult.filterModifiedEndpoints(r)
}

func filterEndpoints(endpoints []string, r *regexp.Regexp) []string {
	result := []string{}
	for _, endpoint := range endpoints {
		if r.MatchString(endpoint) {
			result = append(result, endpoint)
		}
	}

	return result
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
