package diff

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

type PathDiff struct {
	AddedEndpoints    []string          `json:"added,omitempty"`
	DeletedEndpoints  []string          `json:"deleted,omitempty"`
	ModifiedEndpoints ModifiedEndpoints `json:"modified,omitempty"`
}

func (pathDiff *PathDiff) empty() bool {
	return len(pathDiff.AddedEndpoints) == 0 &&
		len(pathDiff.DeletedEndpoints) == 0 &&
		len(pathDiff.ModifiedEndpoints) == 0
}

func newPathDiff() *PathDiff {
	return &PathDiff{
		AddedEndpoints:    []string{},
		DeletedEndpoints:  []string{},
		ModifiedEndpoints: ModifiedEndpoints{},
	}
}

func diffPaths(paths1 openapi3.Paths, paths2 openapi3.Paths, prefix string) *PathDiff {

	result := newPathDiff()

	addedEndpoints, deletedEndpoints, otherEndpoints := diffEndpoints(paths1, paths2, prefix)

	for endpoint := range addedEndpoints {
		result.addAddedEndpoint(endpoint)
	}

	for endpoint := range deletedEndpoints {
		result.addDeletedEndpoint(endpoint)
	}

	for endpoint, pathItemPair := range otherEndpoints {
		result.addModifiedEndpoint(endpoint, pathItemPair.PathItem1, pathItemPair.PathItem2)
	}

	return result
}

func (pathDiff *PathDiff) getSummary() *PathSummary {
	return &PathSummary{
		Added:    len(pathDiff.AddedEndpoints),
		Deleted:  len(pathDiff.DeletedEndpoints),
		Modified: len(pathDiff.ModifiedEndpoints),
	}
}

func (pathDiff *PathDiff) addAddedEndpoint(endpoint string) {
	pathDiff.AddedEndpoints = append(pathDiff.AddedEndpoints, endpoint)
}

func (pathDiff *PathDiff) addDeletedEndpoint(endpoint string) {
	pathDiff.DeletedEndpoints = append(pathDiff.DeletedEndpoints, endpoint)
}

func (pathDiff *PathDiff) addModifiedEndpoint(entrypoint1 string, pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) {
	pathDiff.ModifiedEndpoints.addEndpointDiff(entrypoint1, pathItem1, pathItem2)
}

func (pathDiff *PathDiff) filterByRegex(filter string) {
	r, err := regexp.Compile(filter)
	if err != nil {
		log.Errorf("Failed to compile filter regex '%s' with '%v'", filter, err)
		return
	}

	pathDiff.AddedEndpoints = filterEndpoints(pathDiff.AddedEndpoints, r)
	pathDiff.DeletedEndpoints = filterEndpoints(pathDiff.DeletedEndpoints, r)
	pathDiff.ModifiedEndpoints = filterModifiedEndpoints(pathDiff.ModifiedEndpoints, r)
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

func filterModifiedEndpoints(modifiedEndpoints ModifiedEndpoints, r *regexp.Regexp) ModifiedEndpoints {
	result := ModifiedEndpoints{}

	for endpoint, endpointDiff := range modifiedEndpoints {
		if r.MatchString(endpoint) {
			result[endpoint] = endpointDiff
		}
	}

	return result
}
