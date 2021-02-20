package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type EndpointDiff struct {
	Operations `json:"operations,omitempty"`
}

type Operations struct {
	AddedOperations    OperationList      `json:"added,omitempty"`
	DeletedOperations  OperationList      `json:"deleted,omitempty"`
	ModifiedOperations ModifiedOperations `json:"modified,omitempty"`
}

type OperationList []string

func newEndpointDiff() *EndpointDiff {
	return &EndpointDiff{
		Operations: Operations{
			AddedOperations:    OperationList{},
			DeletedOperations:  OperationList{},
			ModifiedOperations: ModifiedOperations{},
		},
	}
}

func (endpointDiff *EndpointDiff) empty() bool {
	return len(endpointDiff.AddedOperations) == 0 &&
		len(endpointDiff.DeletedOperations) == 0 &&
		endpointDiff.ModifiedOperations.empty()
}

func (endpointDiff *EndpointDiff) diffOperation(pathItem1 *openapi3.Operation, pathItem2 *openapi3.Operation, method string) {

	if pathItem1 == nil && pathItem2 == nil {
		return
	}

	if pathItem1 == nil && pathItem2 != nil {
		endpointDiff.AddedOperations = append(endpointDiff.AddedOperations, method)
		return
	}

	if pathItem1 != nil && pathItem2 == nil {
		endpointDiff.DeletedOperations = append(endpointDiff.AddedOperations, method)
		return
	}

	if diff := diffParameters(pathItem1.Parameters, pathItem2.Parameters); !diff.empty() {
		endpointDiff.ModifiedOperations[method] = diff
	}
}

func diffEndpoint(pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) *EndpointDiff {
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
