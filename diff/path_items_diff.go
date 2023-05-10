package diff

import (
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
)

type pathItemPair struct {
	PathItem1     *openapi3.PathItem
	PathItem2     *openapi3.PathItem
	PathParamsMap PathParamsMap
}

type pathItemPairs map[string]*pathItemPair

func getPathItemsDiff(paths1, paths2 openapi3.Paths) (openapi3.Paths, openapi3.Paths, pathItemPairs) {

	added := openapi3.Paths{}
	deleted := openapi3.Paths{}
	other := pathItemPairs{}

	for endpoint1, pathItem1 := range paths1 {
		if pathItem2, pathParamsMap, ok := findEndpoint(endpoint1, paths2); ok {
			other[endpoint1] = &pathItemPair{
				PathItem1:     pathItem1,
				PathItem2:     pathItem2,
				PathParamsMap: pathParamsMap,
			}
		} else {
			deleted[endpoint1] = pathItem1
		}
	}

	for endpoint2, pathItem2 := range paths2 {
		if _, _, ok := findEndpoint(endpoint2, paths1); !ok {
			added[endpoint2] = pathItem2
		}
	}

	return added, deleted, other
}

func rewritePrefix(paths openapi3.Paths, strip, prepend string) openapi3.Paths {
	result := make(openapi3.Paths, len(paths))
	for path, pathItem := range paths {
		result[prepend+strings.TrimPrefix(path, strip)] = pathItem
	}
	return result
}

func findEndpoint(endpoint string, paths openapi3.Paths) (*openapi3.PathItem, PathParamsMap, bool) {
	if pathItem, ok := paths[endpoint]; ok {
		return pathItem, PathParamsMap{}, true
	}

	return findNormalizedEndpoint(endpoint, paths)
}

/*
findNormalizedEndpoint finds a corresponding path ignoring differences in template variable names

This implementation is based on Paths.Find in openapi3
*/
func findNormalizedEndpoint(key string, paths openapi3.Paths) (*openapi3.PathItem, PathParamsMap, bool) {
	normalizedPath, expected, pathParams1 := normalizeTemplatedPath(key)
	for path, pathItem := range paths {
		pathNormalized, got, pathParams2 := normalizeTemplatedPath(path)
		if got == expected && pathNormalized == normalizedPath {
			if pathParamsMap, ok := NewPathParamsMap(pathParams1, pathParams2); ok {
				return pathItem, pathParamsMap, true
			}
		}
	}
	return nil, nil, false
}

/*
normalizeTemplatedPath converts a path to its normalized form, without parameter names

For example:
/person/{personName} -> /person/{}

Return values:
1. The normalized path
2. Number of params
3. List of param names

This implementation is based on Paths.normalizeTemplatedPath in openapi3
*/
func normalizeTemplatedPath(path string) (string, uint, []string) {
	if strings.IndexByte(path, '{') < 0 {
		return path, 0, nil
	}

	var buffTpl strings.Builder
	buffTpl.Grow(len(path))

	var (
		cc         rune
		count      uint
		isVariable bool
		vars       = []string{}
		buffVar    strings.Builder
	)
	for i, c := range path {
		if isVariable {
			if c == '}' {
				// End path variable
				isVariable = false

				vars = append(vars, buffVar.String())
				buffVar = strings.Builder{}

				// First append possible '*' before this character
				// The character '}' will be appended
				if i > 0 && cc == '*' {
					buffTpl.WriteRune(cc)
				}
			} else {
				buffVar.WriteRune(c)
				continue
			}

		} else if c == '{' {
			// Begin path variable
			isVariable = true

			// The character '{' will be appended
			count++
		}

		// Append the character
		buffTpl.WriteRune(c)
		cc = c
	}
	return buffTpl.String(), count, vars
}
