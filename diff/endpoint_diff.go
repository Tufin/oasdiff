package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type EndpointDiff struct {
	MissingOperations  MissingOperations  `json:"missingOperations,omitempty"`
	ModifiedOperations ModifiedOperations `json:"modifiedOperations,omitempty"`
}

type MissingOperations map[string]struct{}

func newEndpointDiff() *EndpointDiff {
	return &EndpointDiff{
		MissingOperations:  MissingOperations{},
		ModifiedOperations: ModifiedOperations{},
	}
}

func (endpointDiff *EndpointDiff) empty() bool {
	return len(endpointDiff.MissingOperations) == 0 && endpointDiff.ModifiedOperations.empty()
}

func (endpointDiff *EndpointDiff) diffOperation(pathItem1 *openapi3.Operation, pathItem2 *openapi3.Operation, method string) {
	if pathItem1 == nil {
		return
	}

	if pathItem2 == nil {
		endpointDiff.MissingOperations[method] = struct{}{}
		return
	}

	diff := diffParameters(pathItem1.Parameters, pathItem2.Parameters)
	if !diff.empty() {
		endpointDiff.ModifiedOperations[method] = diff
	}
}

func diffEndpoints(pathItem1 *openapi3.PathItem, pathItem2 *openapi3.PathItem) *EndpointDiff {
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
