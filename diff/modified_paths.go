package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ModifiedPaths is a map of paths to their respective diffs
type ModifiedPaths map[string]*PathDiff

func (modifiedPaths ModifiedPaths) addPathDiff(path1 string, pathItem1, pathItem2 *openapi3.PathItem) {

	if diff := getPathDiff(pathItem1, pathItem2); !diff.empty() {
		modifiedPaths[path1] = diff
	}
}
