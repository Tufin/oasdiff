package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type EndpointDiff struct {
	Operations `json:"operations,omitempty"`
}

type Operations struct {
	Added    OperationList      `json:"added,omitempty"`
	Deleted  OperationList      `json:"deleted,omitempty"`
	Modified ModifiedOperations `json:"modified,omitempty"`
}

type OperationList []string

func newEndpointDiff() *EndpointDiff {
	return &EndpointDiff{
		Operations: Operations{
			Added:    OperationList{},
			Deleted:  OperationList{},
			Modified: ModifiedOperations{},
		},
	}
}

func (endpointDiff *EndpointDiff) empty() bool {
	return len(endpointDiff.Added) == 0 &&
		len(endpointDiff.Deleted) == 0 &&
		endpointDiff.Modified.empty()
}

func (endpointDiff *EndpointDiff) diffOperation(pathItem1, pathItem2 *openapi3.Operation, method string) {

	if pathItem1 == nil && pathItem2 == nil {
		return
	}

	if pathItem1 == nil && pathItem2 != nil {
		endpointDiff.Added = append(endpointDiff.Added, method)
		return
	}

	if pathItem1 != nil && pathItem2 == nil {
		endpointDiff.Deleted = append(endpointDiff.Added, method)
		return
	}

	if diff := getMethodDiff(pathItem1, pathItem2); !diff.empty() {
		endpointDiff.Modified[method] = diff
	}
}

func diffEndpoint(pathItem1, pathItem2 *openapi3.PathItem) *EndpointDiff {
	endpointDiff := newEndpointDiff()

	endpointDiff.diffOperation(pathItem1.Connect, pathItem2.Connect, "CONNECT")
	endpointDiff.diffOperation(pathItem1.Delete, pathItem2.Delete, "DELETE")
	endpointDiff.diffOperation(pathItem1.Get, pathItem2.Get, "GET")
	endpointDiff.diffOperation(pathItem1.Head, pathItem2.Head, "HEAD")
	endpointDiff.diffOperation(pathItem1.Options, pathItem2.Options, "OPTIONS")
	endpointDiff.diffOperation(pathItem1.Patch, pathItem2.Patch, "PATCH")
	endpointDiff.diffOperation(pathItem1.Post, pathItem2.Post, "POST")
	endpointDiff.diffOperation(pathItem1.Put, pathItem2.Put, "PUT")
	endpointDiff.diffOperation(pathItem1.Trace, pathItem2.Trace, "TRACE")

	return endpointDiff
}
