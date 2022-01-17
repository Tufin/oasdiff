package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

/*
EndpointsDiff is an alternate, simplified view of PathsDiff.
It describes the changes between Endpoints which are a flattened combination of OpenAPI Paths and Operations.

For example, if there's a new path "/test" with method POST then EndpointsDiff will describe this as a new endpoint: POST /test.

Or, if path "/test" was modified to include a new methdod, PUT, then EndpointsDiff will describe this as a new endpoint: PUT /test.
*/
type EndpointsDiff struct {
	Added    Endpoints         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  Endpoints         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedEndpoints `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Endpoint is a combination of an HTTP method and a Path
type Endpoint struct {
	Method string `json:"method,omitempty" method:"added,omitempty"`
	Path   string `json:"path,omitempty" yaml:"path,omitempty"`
}

// Empty indicates whether a change was found in this element
func (diff *EndpointsDiff) Empty() bool {
	if diff == nil {
		return true
	}

	return len(diff.Added) == 0 &&
		len(diff.Deleted) == 0 &&
		len(diff.Modified) == 0
}

func (diff *EndpointsDiff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.Added = nil
}

func newEndpointsDiff() *EndpointsDiff {
	return &EndpointsDiff{
		Added:    Endpoints{},
		Deleted:  Endpoints{},
		Modified: ModifiedEndpoints{},
	}
}

func getEndpointsDiff(config *Config, paths1, paths2 openapi3.Paths) (*EndpointsDiff, error) {

	err := filterPaths2(config.PathFilter, paths1, paths2)
	if err != nil {
		return nil, err
	}

	diff, err := getEndpointsDiffInternal(config, paths1, paths2)
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

func getEndpointsDiffInternal(config *Config, paths1, paths2 openapi3.Paths) (*EndpointsDiff, error) {

	result := newEndpointsDiff()

	addedPaths, deletedPaths, otherPaths := getPathItemsDiff(paths1, paths2, config.PathPrefix)

	for path, pathItem := range addedPaths {
		for method := range pathItem.Operations() {
			result.addAddedPath(path, method)
		}
	}

	for path, pathItem := range deletedPaths {
		for method := range pathItem.Operations() {
			result.addDeletedPath(path, method)
		}
	}

	for path, pathItemPair := range otherPaths {
		result.addModifiedPaths(config, path, pathItemPair)
	}

	return result, nil
}

func (diff *EndpointsDiff) addAddedPath(path string, method string) {
	diff.Added = append(diff.Added, Endpoint{
		Method: method,
		Path:   path,
	})
}

func (diff *EndpointsDiff) addDeletedPath(path string, method string) {
	diff.Deleted = append(diff.Deleted, Endpoint{
		Method: method,
		Path:   path,
	})
}

func (diff *EndpointsDiff) addModifiedPaths(config *Config, path string, pathItemPair *pathItemPair) error {

	pathDiff, err := getPathDiff(config, pathItemPair.PathItem1, pathItemPair.PathItem2)
	if err != nil {
		return err
	}

	if pathDiff.Empty() || pathDiff.OperationsDiff.Empty() {
		return nil
	}

	for _, method := range pathDiff.OperationsDiff.Added {
		diff.addAddedPath(path, method)
	}

	for _, method := range pathDiff.OperationsDiff.Deleted {
		diff.addDeletedPath(path, method)
	}

	for method, methodDiff := range pathDiff.OperationsDiff.Modified {
		diff.Modified[Endpoint{
			Method: method,
			Path:   path,
		}] = methodDiff
	}

	return nil
}

func (diff *EndpointsDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(diff.Added),
		Deleted:  len(diff.Deleted),
		Modified: len(diff.Modified),
	}
}
