package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedPaths maps paths to their diff
type ModifiedPaths map[string]*PathDiff

func (modifiedPaths ModifiedPaths) addPathDiff(path1 string, pathItem1, pathItem2 *openapi3.PathItem) {

	if diff := diffEndpoint(pathItem1, pathItem2); !diff.empty() {
		modifiedPaths[path1] = diff
	}
}
