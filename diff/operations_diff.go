package diff

import (
	"encoding/json"
	"fmt"
	"net/http"
	"regexp"
	"time"

	"github.com/getkin/kin-openapi/openapi3"
)

// OperationsDiff describes the changes between a pair of operation objects (https://swagger.io/specification/#operation-object) of two path item objects
type OperationsDiff struct {
	Added    StringList         `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList         `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedOperations `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (operationsDiff *OperationsDiff) Empty() bool {
	if operationsDiff == nil {
		return true
	}

	return len(operationsDiff.Added) == 0 &&
		len(operationsDiff.Deleted) == 0 &&
		len(operationsDiff.Modified) == 0
}

func properlyDeprecated(config *Config, state *state, pathItem1, pathItem2 *openapi3.PathItem, op string) bool {
	if pathItem1 == nil {
		return false
	}

	operation := pathItem1.GetOperation(op)

	if !operation.Deprecated {
		return false
	}

	sunsetJson, ok := operation.ExtensionProps.Extensions["x-sunset"].(json.RawMessage)
	if !ok {
		return false
	}

	var sunset string
	if err := json.Unmarshal(sunsetJson, &sunset); err != nil {
		return false
	}

	date, err := time.Parse("2006-01-02", sunset)
	if err != nil {
		return false
	}

	return time.Now().After(date)
}

func (operationsDiff *OperationsDiff) removeProperlyDeprecated(config *Config, state *state, pathItem1, pathItem2 *openapi3.PathItem) {
	deleted := []string{}
	for _, op := range operationsDiff.Deleted {
		if !properlyDeprecated(config, state, pathItem1, pathItem2, op) {
			deleted = append(deleted, op)
		}
	}
	operationsDiff.Deleted = deleted

}

func (operationsDiff *OperationsDiff) removeNonBreaking(config *Config, state *state, pathItem1, pathItem2 *openapi3.PathItem) {

	if operationsDiff.Empty() {
		return
	}

	operationsDiff.removeProperlyDeprecated(config, state, pathItem1, pathItem2)
	operationsDiff.Added = nil
}

func newOperationsDiff() *OperationsDiff {
	return &OperationsDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedOperations{},
	}
}

// ModifiedOperations is a map of HTTP methods to their respective diffs
type ModifiedOperations map[string]*MethodDiff

func getOperationsDiff(config *Config, state *state, pathItem1, pathItem2 *openapi3.PathItem) (*OperationsDiff, error) {
	if err := filterOperations(config.FilterExtension, pathItem1, pathItem2); err != nil {
		return nil, err
	}

	diff, err := getOperationsDiffInternal(config, state, pathItem1, pathItem2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking(config, state, pathItem1, pathItem2)
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

var operations = []string{
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func getOperationsDiffInternal(config *Config, state *state, pathItem1, pathItem2 *openapi3.PathItem) (*OperationsDiff, error) {

	result := newOperationsDiff()
	var err error

	for _, op := range operations {
		err = result.diffOperation(config, state, pathItem1.GetOperation(op), pathItem2.GetOperation(op), op)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (operationsDiff *OperationsDiff) diffOperation(config *Config, state *state, operation1, operation2 *openapi3.Operation, method string) error {
	if operation1 == nil && operation2 == nil {
		return nil
	}

	if operation1 == nil && operation2 != nil {
		operationsDiff.Added = append(operationsDiff.Added, method)
		return nil
	}

	if operation1 != nil && operation2 == nil {
		operationsDiff.Deleted = append(operationsDiff.Deleted, method)
		return nil
	}

	diff, err := getMethodDiff(config, state, operation1, operation2)
	if err != nil {
		return err
	}

	if !diff.Empty() {
		operationsDiff.Modified[method] = diff
	}

	return nil
}

func filterOperations(filterExtension string, pathItem1, pathItem2 *openapi3.PathItem) error {

	if err := filterOperationsByExtensions(filterExtension, pathItem1, pathItem2); err != nil {
		return err
	}

	return nil
}

func filterOperationsByExtensions(filterExtension string, pathItem1, pathItem2 *openapi3.PathItem) error {
	if filterExtension == "" {
		return nil
	}

	r, err := regexp.Compile(filterExtension)
	if err != nil {
		return fmt.Errorf("failed to compile extension filter regex %q with %w", filterExtension, err)
	}

	filterOperationsByExtensionInternal(pathItem1, r)
	filterOperationsByExtensionInternal(pathItem2, r)

	return nil
}

func filterOperationsByExtensionInternal(pathItem *openapi3.PathItem, r *regexp.Regexp) {
	for method, operation := range pathItem.Operations() {
		for extension := range operation.Extensions {
			if r.MatchString(extension) {
				pathItem.SetOperation(method, nil)
				break
			}
		}
	}
}

// Patch applies the patch to operations
func (operationsDiff *OperationsDiff) Patch(operations map[string]*openapi3.Operation) error {

	if operationsDiff.Empty() {
		return nil
	}

	for method, methodDiff := range operationsDiff.Modified {
		err := methodDiff.Patch(operations[method])
		if err != nil {
			return err
		}
	}

	return nil
}
