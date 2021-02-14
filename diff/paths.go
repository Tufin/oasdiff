package diff

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

type PathsDiff struct {
	AddedEndpoints    []string          `json:"addedEndpoints,omitempty"`
	DeletedEndpoints  []string          `json:"deletedEndpoints,omitempty"`
	ModifiedEndpoints ModifiedEndpoints `json:"modifiedEndpoints,omitempty"`
}

func (pathsDiff *PathsDiff) empty() bool {
	return len(pathsDiff.AddedEndpoints) == 0 &&
		len(pathsDiff.DeletedEndpoints) == 0 &&
		len(pathsDiff.ModifiedEndpoints) == 0
}

func newPathsDiff() *PathsDiff {
	return &PathsDiff{
		AddedEndpoints:    []string{},
		DeletedEndpoints:  []string{},
		ModifiedEndpoints: ModifiedEndpoints{},
	}
}

func diffPaths(paths1 openapi3.Paths, paths2 openapi3.Paths, prefix string) *PathsDiff {

	result := newPathsDiff()

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

func (pathsDiff *PathsDiff) getSummary() PathsDiffSummary {
	return PathsDiffSummary{
		AddedEndpoints:    len(pathsDiff.AddedEndpoints),
		DeletedEndpoints:  len(pathsDiff.DeletedEndpoints),
		ModifiedEndpoints: len(pathsDiff.ModifiedEndpoints),
	}
}
func (pathsDiff *PathsDiff) addAddedEndpoint(endpoint string) {
	pathsDiff.AddedEndpoints = append(pathsDiff.AddedEndpoints, endpoint)
}

func (pathsDiff *PathsDiff) addDeletedEndpoint(endpoint string) {
	pathsDiff.DeletedEndpoints = append(pathsDiff.DeletedEndpoints, endpoint)
}

func (pathsDiff *PathsDiff) addModifiedEndpoint(entrypoint1 string, pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) {
	pathsDiff.ModifiedEndpoints.addEndpointDiff(entrypoint1, pathItem1, pathItem2)
}

func (pathsDiff *PathsDiff) filterByRegex(filter string) {
	r, err := regexp.Compile(filter)
	if err != nil {
		log.Errorf("Failed to compile filter regex '%s' with '%v'", filter, err)
		return
	}

	pathsDiff.AddedEndpoints = filterEndpoints(pathsDiff.AddedEndpoints, r)
	pathsDiff.DeletedEndpoints = filterEndpoints(pathsDiff.DeletedEndpoints, r)
	pathsDiff.ModifiedEndpoints = filterModifiedEndpoints(pathsDiff.ModifiedEndpoints, r)
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
