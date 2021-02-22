package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// PathDiff is the diff between two OAS paths
type PathDiff struct {
	Operations `json:"operations,omitempty"`
}

// Operations is the diff between two sets of operations (methods)
type Operations struct {
	Added    OperationList      `json:"added,omitempty"`
	Deleted  OperationList      `json:"deleted,omitempty"`
	Modified ModifiedOperations `json:"modified,omitempty"`
}

// OperationList is a list of operations (methods)
type OperationList []string

func newPathDiff() *PathDiff {
	return &PathDiff{
		Operations: Operations{
			Added:    OperationList{},
			Deleted:  OperationList{},
			Modified: ModifiedOperations{},
		},
	}
}

func (pathDiff *PathDiff) empty() bool {
	return len(pathDiff.Added) == 0 &&
		len(pathDiff.Deleted) == 0 &&
		pathDiff.Modified.empty()
}

func (pathDiff *PathDiff) diffOperation(pathItem1, pathItem2 *openapi3.Operation, method string) {

	if pathItem1 == nil && pathItem2 == nil {
		return
	}

	if pathItem1 == nil && pathItem2 != nil {
		pathDiff.Added = append(pathDiff.Added, method)
		return
	}

	if pathItem1 != nil && pathItem2 == nil {
		pathDiff.Deleted = append(pathDiff.Added, method)
		return
	}

	if diff := getMethodDiff(pathItem1, pathItem2); !diff.empty() {
		pathDiff.Modified[method] = diff
	}
}

func getPathDiff(pathItem1, pathItem2 *openapi3.PathItem) *PathDiff {
	pathDiff := newPathDiff()

	pathDiff.diffOperation(pathItem1.Connect, pathItem2.Connect, "CONNECT")
	pathDiff.diffOperation(pathItem1.Delete, pathItem2.Delete, "DELETE")
	pathDiff.diffOperation(pathItem1.Get, pathItem2.Get, "GET")
	pathDiff.diffOperation(pathItem1.Head, pathItem2.Head, "HEAD")
	pathDiff.diffOperation(pathItem1.Options, pathItem2.Options, "OPTIONS")
	pathDiff.diffOperation(pathItem1.Patch, pathItem2.Patch, "PATCH")
	pathDiff.diffOperation(pathItem1.Post, pathItem2.Post, "POST")
	pathDiff.diffOperation(pathItem1.Put, pathItem2.Put, "PUT")
	pathDiff.diffOperation(pathItem1.Trace, pathItem2.Trace, "TRACE")

	return pathDiff
}
