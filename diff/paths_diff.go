package diff

import (
	"fmt"
	"regexp"

	"github.com/getkin/kin-openapi/openapi3"
)

// PathsDiff describes the changes between a pair of Paths objects: https://swagger.io/specification/#paths-object
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

func (pathsDiff *PathsDiff) removeNonBreaking() {

	if pathsDiff.Empty() {
		return
	}

	pathsDiff.Added = nil
}

func newPathsDiff() *PathsDiff {
	return &PathsDiff{
		Added:    []string{},
		Deleted:  []string{},
		Modified: ModifiedPaths{},
	}
}

func getPathsDiff(config *Config, state *state, paths1, paths2 openapi3.Paths) (*PathsDiff, error) {

	if err := filterPaths(config.PathFilter, config.FilterExtension, paths1, paths2); err != nil {
		return nil, err
	}

	diff, err := getPathsDiffInternal(config, state, paths1, paths2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getPathsDiffInternal(config *Config, state *state, paths1, paths2 openapi3.Paths) (*PathsDiff, error) {

	result := newPathsDiff()

	addedPaths, deletedPaths, otherPaths := getPathItemsDiff(paths1, paths2, config.PathPrefix)

	for endpoint := range addedPaths {
		result.addAddedPath(endpoint)
	}

	for endpoint := range deletedPaths {
		result.addDeletedPath(endpoint)
	}

	for endpoint, pathItemPair := range otherPaths {
		err := result.addModifiedPath(config, state, endpoint, pathItemPair.PathItem1, pathItemPair.PathItem2)
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

func (pathsDiff *PathsDiff) addModifiedPath(config *Config, state *state, path1 string, pathItem1, pathItem2 *openapi3.PathItem) error {
	return pathsDiff.Modified.addPathDiff(config, state, path1, pathItem1, pathItem2)
}

func filterPaths(filter, filterExtension string, paths1, paths2 openapi3.Paths) error {

	if err := filterPathsByName(filter, paths1, paths2); err != nil {
		return err
	}

	if err := filterPathsByExtensions(filterExtension, paths1, paths2); err != nil {
		return err
	}

	return nil
}

func filterPathsByName(filter string, paths1, paths2 openapi3.Paths) error {
	if filter == "" {
		return nil
	}

	r, err := regexp.Compile(filter)
	if err != nil {
		return fmt.Errorf("failed to compile filter regex %q with %w", filter, err)
	}

	filterPathsInternal(paths1, r)
	filterPathsInternal(paths2, r)

	return nil
}

func filterPathsInternal(paths openapi3.Paths, r *regexp.Regexp) {
	for path := range paths {
		if !r.MatchString(path) {
			delete(paths, path)
		}
	}
}

func filterPathsByExtensions(filterExtension string, paths1, paths2 openapi3.Paths) error {
	if filterExtension == "" {
		return nil
	}

	r, err := regexp.Compile(filterExtension)
	if err != nil {
		return fmt.Errorf("failed to compile extension filter regex %q with %w", filterExtension, err)
	}

	filterPathsByExtensionInternal(paths1, r)
	filterPathsByExtensionInternal(paths2, r)

	return nil
}

func filterPathsByExtensionInternal(paths openapi3.Paths, r *regexp.Regexp) {
	for path, pathItem := range paths {
		for extension := range pathItem.Extensions {
			if r.MatchString(extension) {
				delete(paths, path)
				break
			}
		}
	}
}

// Patch applies the patch to paths
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
