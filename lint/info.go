package lint

import (
	"fmt"
	"net/url"

	"github.com/tufin/oasdiff/load"
)

// InfoCheck based on REQUIRED fields (Version and Info) from swagger docs,
// see: https://swagger.io/docs/specification/api-general-info/
func InfoCheck(source string, spec *load.SpecInfo) []*Error {

	result := make([]*Error, 0)

	if spec == nil || spec.Spec == nil {
		return result
	}

	if spec.Spec.Info == nil {
		result = append(result, &Error{
			Id:      "info-missing",
			Level:   LEVEL_ERROR,
			Text:    "info is missing",
			Comment: "It is a good practice to include general information about your API into the specification. Title and Version fields are required.",
			Source:  source,
		})
		return result
	}

	if spec.Spec.Info.Title == "" {
		result = append(result, &Error{
			Id:     "info-title-missing",
			Level:  LEVEL_ERROR,
			Text:   "the title of the API is missing",
			Source: source,
		})
	}
	if spec.Spec.Info.Version == "" {
		result = append(result, &Error{
			Id:     "info-version-missing",
			Level:  LEVEL_ERROR,
			Text:   "the version of the API is missing",
			Source: source,
		})
	}

	if tos := spec.Spec.Info.TermsOfService; tos != "" {
		if _, err := url.ParseRequestURI(tos); err != nil {
			result = append(result, &Error{
				Id:     "info-invalid-terms-of-service",
				Level:  LEVEL_ERROR,
				Text:   fmt.Sprintf("terms of service must be in the format of a URL: %s", tos),
				Source: source,
			})
		}
	}

	return result
}
