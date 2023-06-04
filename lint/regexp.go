package lint

import (
	"regexp"
)

func checkRegex(pattern string, s *state) *Error {
	if pattern == "" {
		return nil
	}

	if err := s.validate(pattern); err != nil {
		return &Error{
			Id:     "invalid-regex-pattern",
			Level:  LEVEL_ERROR,
			Text:   err.Error(),
			Source: s.source,
		}
	}

	return nil
}

// *** THIS IS A TEMPORARY IMPLEMENTATION ***
// SHOULD USE ECMA 262, SEE: https://swagger.io/docs/specification/data-models/data-types/#pattern
func (s *state) validate(pattern string) error {
	if result, ok := s.cache[pattern]; ok {
		return result
	}
	_, err := regexp.Compile(pattern)
	s.cache[pattern] = err
	return err
}
