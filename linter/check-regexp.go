package linter

import (
	"regexp"

	"github.com/tufin/oasdiff/load"
)

// *** THIS IS A TEMPORARY IMPLEMENTATION ***
// SHOULD USE ECMA 262, SEE: https://swagger.io/docs/specification/data-models/data-types/#pattern

func RegexCheck(spec *load.OpenAPISpecInfo) []Error {
	result := make([]Error, 0)

	for _, path := range spec.Spec.Paths {
		for _, parameter := range path.Parameters {
			pattern := parameter.Value.Schema.Value.Pattern
			if pattern != "" {
				_, err := regexp.Compile(pattern)
				if err != nil {
					result = append(result, Error{
						Id:    "invalid-regex-pattern",
						Level: LEVEL_ERROR,
						Text:  err.Error(),
						// Path:   path,
					})
				}
			}
		}

	}
	return result
}
