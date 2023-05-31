package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	apiSchemasRemovedCheckId = "api-schema-removed"
)

func APIComponentsSchemaRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.ComponentsDiff.SchemasDiff == nil {
		return result
	}

	for _, pathItem := range diffReport.ComponentsDiff.SchemasDiff.Deleted {
		result = append(result, BackwardCompatibilityError{
			Id:        apiSchemasRemovedCheckId,
			Level:     ERR,
			Text:      config.i18n(apiSchemasRemovedCheckId),
			Operation: "",
			Path:      pathItem,
			Source:    "",
		})
	}
	return result
}
