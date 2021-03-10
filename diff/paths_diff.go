package diff

import (
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
)

// PathsDiff is a diff between two sets of path item objects: https://swagger.io/specification/#path-item-object
type PathsDiff struct {
	Added    StringList    `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList    `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedPaths `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty return true if there is no diff
func (pathsDiff *PathsDiff) Empty() bool {
	if pathsDiff == nil {
		return true
	}

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

func getPathsDiff(config *Config, paths1, paths2 openapi3.Paths) *PathsDiff {

	filterPaths2(config.Filter, paths1, paths2)

	diff := getPathsDiffInternal(config, paths1, paths2)

	if diff.Empty() {
		return nil
	}

	return diff
}

func getPathsDiffInternal(config *Config, paths1, paths2 openapi3.Paths) *PathsDiff {

	result := newPathsDiff()

	addedEndpoints, deletedEndpoints, otherEndpoints := getEndpointsDiff(paths1, paths2, config.Prefix)

	for endpoint := range addedEndpoints {
		result.addAddedPath(endpoint)
	}

	for endpoint := range deletedEndpoints {
		result.addDeletedPath(endpoint)
	}

	for endpoint, pathItemPair := range otherEndpoints {
		result.addModifiedPath(config, endpoint, pathItemPair.PathItem1, pathItemPair.PathItem2)
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

func (pathsDiff *PathsDiff) addModifiedPath(config *Config, path1 string, pathItem1, pathItem2 *openapi3.PathItem) {
	pathsDiff.Modified.addPathDiff(config, path1, pathItem1, pathItem2)
}

func filterPaths2(filter string, paths1, paths2 openapi3.Paths) {

	if filter == "" {
		return
	}

	r, err := regexp.Compile(filter)
	if err != nil {
		log.Errorf("Failed to compile filter regex '%s' with '%v'", filter, err)
		return
	}

	filterPaths1(paths1, r)
	filterPaths1(paths2, r)
}

func filterPaths1(paths openapi3.Paths, r *regexp.Regexp) {
	for path := range paths {
		if !r.MatchString(path) {
			delete(paths, path)
		}
	}
}
