package diff

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

// PathsDiff is a diff between two sets of paths
type PathsDiff struct {
	Added    PathList      `json:"added,omitempty"`
	Deleted  PathList      `json:"deleted,omitempty"`
	Modified ModifiedPaths `json:"modified,omitempty"`
}

// PathList is a list of paths
type PathList []string

func (pathsDiff *PathsDiff) empty() bool {
	return len(pathsDiff.Added) == 0 &&
		len(pathsDiff.Deleted) == 0 &&
		len(pathsDiff.Modified) == 0
}

func newPathsDiff() *PathsDiff {
	return &PathsDiff{
		Added:    []string{},
		Deleted:  []string{},
		Modified: ModifiedPaths{},
	}
}

func getPathsDiff(paths1, paths2 openapi3.Paths, prefix string) *PathsDiff {

	result := newPathsDiff()

	addedEndpoints, deletedEndpoints, otherEndpoints := getEndpointsDiff(paths1, paths2, prefix)

	for endpoint := range addedEndpoints {
		result.addAddedPath(endpoint)
	}

	for endpoint := range deletedEndpoints {
		result.addDeletedPath(endpoint)
	}

	for endpoint, pathItemPair := range otherEndpoints {
		result.addModifiedPath(endpoint, pathItemPair.PathItem1, pathItemPair.PathItem2)
	}

	return result
}

func (pathsDiff *PathsDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(pathsDiff.Added),
		Deleted:  len(pathsDiff.Deleted),
		Modified: len(pathsDiff.Modified),
	}
}

func (pathsDiff *PathsDiff) addAddedPath(path string) {
	pathsDiff.Added = append(pathsDiff.Added, path)
}

func (pathsDiff *PathsDiff) addDeletedPath(path string) {
	pathsDiff.Deleted = append(pathsDiff.Deleted, path)
}

func (pathsDiff *PathsDiff) addModifiedPath(path1 string, pathItem1, pathItem2 *openapi3.PathItem) {
	pathsDiff.Modified.addPathDiff(path1, pathItem1, pathItem2)
}

func (pathsDiff PathsDiff) filterByRegex(filter string) *PathsDiff {

	r, err := regexp.Compile(filter)
	if err != nil {
		log.Errorf("Failed to compile filter regex '%s' with '%v'", filter, err)
		return &pathsDiff
	}

	pathsDiff.Added = filterPaths(pathsDiff.Added, r)
	pathsDiff.Deleted = filterPaths(pathsDiff.Deleted, r)
	pathsDiff.Modified = filterModifiedPaths(pathsDiff.Modified, r)

	return &pathsDiff
}

func filterPaths(paths PathList, r *regexp.Regexp) []string {
	result := []string{}
	for _, path := range paths {
		if r.MatchString(path) {
			result = append(result, path)
		}
	}

	return result
}

func filterModifiedPaths(modifiedPaths ModifiedPaths, r *regexp.Regexp) ModifiedPaths {
	result := ModifiedPaths{}

	for endpoint, endpointDiff := range modifiedPaths {
		if r.MatchString(endpoint) {
			result[endpoint] = endpointDiff
		}
	}

	return result
}
