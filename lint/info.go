package lint

import (
	"github.com/tufin/oasdiff/load"
)

func InfoCheck(source string, spec *load.OpenAPISpecInfo) []*Error {

	result := make([]*Error, 0)

	if spec == nil || spec.Spec == nil {
		return result
	}

	if spec.Spec.Info.Title == "" {
		result = append(result, &Error{
			Id:     "info-title-missing",
			Level:  LEVEL_WARN,
			Text:   "It is a good practice to include general information about your API into the specification: title is missing",
			Source: source,
		})
	}
	if spec.Spec.Info.Version == "" {
		result = append(result, &Error{
			Id:     "info-version-missing",
			Level:  LEVEL_WARN,
			Text:   "It is a good practice to include general information about your API into the specification: version number is missing",
			Source: source,
		})
	}
	if spec.Spec.Info.License.Name == "" {
		result = append(result, &Error{
			Id:     "info-version-missing",
			Level:  LEVEL_WARN,
			Text:   "It is a good practice to include general information about your API into the specification: license name is missing",
			Source: source,
		})
	}
	if spec.Spec.Info.License.URL == "" {
		result = append(result, &Error{
			Id:     "info-version-missing",
			Level:  LEVEL_WARN,
			Text:   "It is a good practice to include general information about your API into the specification: license URL is missing",
			Source: source,
		})
	}

	return result
}
