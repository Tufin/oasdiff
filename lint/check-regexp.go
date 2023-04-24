package lint

import (
	"regexp"

	"github.com/tufin/oasdiff/load"
)

// *** THIS IS A TEMPORARY IMPLEMENTATION ***
// SHOULD USE ECMA 262, SEE: https://swagger.io/docs/specification/data-models/data-types/#pattern

func RegexCheck(source string, s *load.OpenAPISpecInfo) []Error {

	result := make([]Error, 0)

	if s == nil || s.Spec == nil {
		return result
	}

	for _, path := range s.Spec.Paths {
		for _, parameter := range path.Parameters {
			if parameter.Value == nil || parameter.Value.Schema == nil {
				continue
			}
			pattern := parameter.Value.Schema.Value.Pattern
			if pattern != "" {
				_, err := regexp.Compile(pattern)
				if err != nil {
					result = append(result, Error{
						Id:     "invalid-regex-pattern",
						Level:  LEVEL_ERROR,
						Text:   err.Error(),
						Source: source,
					})
				}
			}
		}

	}

	return result
}
