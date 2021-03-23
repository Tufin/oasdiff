package diff

import (
	"net/http"

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

func newOperationsDiff() *OperationsDiff {
	return &OperationsDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedOperations{},
	}
}

// ModifiedOperations is a map of HTTP methods to their respective diffs
type ModifiedOperations map[string]*MethodDiff

func getOperationsDiff(config *Config, pathItem1, pathItem2 *openapi3.PathItem) (*OperationsDiff, error) {
	diff, err := getOperationsDiffInternal(config, pathItem1, pathItem2)
	if err != nil {
		return nil, err
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

func getOperationsDiffInternal(config *Config, pathItem1, pathItem2 *openapi3.PathItem) (*OperationsDiff, error) {

	result := newOperationsDiff()
	var err error

	for _, op := range operations {
		err = result.diffOperation(config, pathItem1.GetOperation(op), pathItem2.GetOperation(op), op)
		if err != nil {
			return nil, err
		}
	}

	return result, nil
}

func (operationsDiff *OperationsDiff) diffOperation(config *Config, operation1, operation2 *openapi3.Operation, method string) error {
	if operation1 == nil && operation2 == nil {
		return nil
	}

	if operation1 == nil && operation2 != nil {
		operationsDiff.Added = append(operationsDiff.Added, method)
		return nil
	}

	if operation1 != nil && operation2 == nil {
		operationsDiff.Deleted = append(operationsDiff.Added, method)
		return nil
	}

	diff, err := getMethodDiff(config, operation1, operation2)
	if err != nil {
		return err
	}

	if !diff.Empty() {
		operationsDiff.Modified[method] = diff
	}

	return nil
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
