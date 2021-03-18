package diff

import (
	"fmt"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
)

// PathsDiff is a diff between two sets of path item objects: https://swagger.io/specification/#path-item-object
type PathsDiff struct {
	Added    StringList    `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList    `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedPaths `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
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

func getPathsDiff(config *Config, paths1, paths2 openapi3.Paths) (*PathsDiff, error) {

	err := filterPaths2(config.Filter, paths1, paths2)
	if err != nil {
		return nil, err
	}

	diff, err := getPathsDiffInternal(config, paths1, paths2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getPathsDiffInternal(config *Config, paths1, paths2 openapi3.Paths) (*PathsDiff, error) {

	result := newPathsDiff()

	addedEndpoints, deletedEndpoints, otherEndpoints := getEndpointsDiff(paths1, paths2, config.Prefix)

	for endpoint := range addedEndpoints {
		result.addAddedPath(endpoint)
	}

	for endpoint := range deletedEndpoints {
		result.addDeletedPath(endpoint)
	}

	for endpoint, pathItemPair := range otherEndpoints {
		err := result.addModifiedPath(config, endpoint, pathItemPair.PathItem1, pathItemPair.PathItem2)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
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

func (pathsDiff *PathsDiff) addModifiedPath(config *Config, path1 string, pathItem1, pathItem2 *openapi3.PathItem) error {
	return pathsDiff.Modified.addPathDiff(config, path1, pathItem1, pathItem2)
}

func filterPaths2(filter string, paths1, paths2 openapi3.Paths) error {

	if filter == "" {
		return nil
	}

	r, err := regexp.Compile(filter)
	if err != nil {
		return fmt.Errorf("failed to compile filter regex '%s' with %w", filter, err)
	}

	filterPaths1(paths1, r)
	filterPaths1(paths2, r)

	return nil
}

func filterPaths1(paths openapi3.Paths, r *regexp.Regexp) {
	for path := range paths {
		if !r.MatchString(path) {
			delete(paths, path)
		}
	}
}

// Apply applies the patch
func (pathsDiff *PathsDiff) Patch(paths openapi3.Paths) error {

	if pathsDiff.Empty() {
		return nil
	}

	for path, pathDiff := range pathsDiff.Modified {
		err := pathDiff.Patch(paths.Find(path))
		if err != nil {
			return err
		}
	}

	return nil
}
