package diff

import "github.com/getkin/kin-openapi/openapi3"

// OperationsDiff is a diff between the operation objects (https://swagger.io/specification/#operation-object) of two path item objects
type OperationsDiff struct {
	Added    StringList         `json:"added,omitempty"`
	Deleted  StringList         `json:"deleted,omitempty"`
	Modified ModifiedOperations `json:"modified,omitempty"`
}

func (operationsDiff *OperationsDiff) empty() bool {
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

func getOperationsDiff(config *Config, pathItem1, pathItem2 *openapi3.PathItem) *OperationsDiff {

	result := newOperationsDiff()

	result.diffOperation(config, pathItem1.Connect, pathItem2.Connect, "CONNECT")
	result.diffOperation(config, pathItem1.Delete, pathItem2.Delete, "DELETE")
	result.diffOperation(config, pathItem1.Get, pathItem2.Get, "GET")
	result.diffOperation(config, pathItem1.Head, pathItem2.Head, "HEAD")
	result.diffOperation(config, pathItem1.Options, pathItem2.Options, "OPTIONS")
	result.diffOperation(config, pathItem1.Patch, pathItem2.Patch, "PATCH")
	result.diffOperation(config, pathItem1.Post, pathItem2.Post, "POST")
	result.diffOperation(config, pathItem1.Put, pathItem2.Put, "PUT")
	result.diffOperation(config, pathItem1.Trace, pathItem2.Trace, "TRACE")

	if result.empty() {
		return nil
	}

	return result
}

func (operationsDiff *OperationsDiff) diffOperation(config *Config, operation1, operation2 *openapi3.Operation, method string) {
	if operation1 == nil && operation2 == nil {
		return
	}

	if operation1 == nil && operation2 != nil {
		operationsDiff.Added = append(operationsDiff.Added, method)
		return
	}

	if operation1 != nil && operation2 == nil {
		operationsDiff.Deleted = append(operationsDiff.Added, method)
		return
	}

	if diff := getMethodDiff(config, operation1, operation2); !diff.empty() {
		operationsDiff.Modified[method] = diff
	}
}
