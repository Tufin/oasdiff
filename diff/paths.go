package diff

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

type PathDiff struct {
	Added    EndpointList      `json:"added,omitempty"`
	Deleted  EndpointList      `json:"deleted,omitempty"`
	Modified ModifiedEndpoints `json:"modified,omitempty"`
}

type EndpointList []string

func (pathDiff *PathDiff) empty() bool {
	return len(pathDiff.Added) == 0 &&
		len(pathDiff.Deleted) == 0 &&
		len(pathDiff.Modified) == 0
}

func newPathDiff() *PathDiff {
	return &PathDiff{
		Added:    []string{},
		Deleted:  []string{},
		Modified: ModifiedEndpoints{},
	}
}

func diffPaths(paths1, paths2 openapi3.Paths, prefix string) *PathDiff {

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
		Added:    len(pathDiff.Added),
		Deleted:  len(pathDiff.Deleted),
		Modified: len(pathDiff.Modified),
	}
}

func (pathDiff *PathDiff) addAddedEndpoint(endpoint string) {
	pathDiff.Added = append(pathDiff.Added, endpoint)
}

func (pathDiff *PathDiff) addDeletedEndpoint(endpoint string) {
	pathDiff.Deleted = append(pathDiff.Deleted, endpoint)
}

func (pathDiff *PathDiff) addModifiedEndpoint(entrypoint1 string, pathItem1, pathItem2 *openapi3.PathItem) {
	pathDiff.Modified.addEndpointDiff(entrypoint1, pathItem1, pathItem2)
}

func (pathDiff PathDiff) filterByRegex(filter string) *PathDiff {

	r, err := regexp.Compile(filter)
	if err != nil {
		log.Errorf("Failed to compile filter regex '%s' with '%v'", filter, err)
		return &pathDiff
	}

	pathDiff.Added = filterEndpoints(pathDiff.Added, r)
	pathDiff.Deleted = filterEndpoints(pathDiff.Deleted, r)
	pathDiff.Modified = filterModifiedEndpoints(pathDiff.Modified, r)

	return &pathDiff
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
