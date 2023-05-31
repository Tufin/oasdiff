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

	for _, deletedSchema := range diffReport.ComponentsDiff.SchemasDiff.Deleted {
		result = append(result, BackwardCompatibilityError{
			Id:        apiSchemasRemovedCheckId,
			Level:     ERR,
			Text:      config.i18n(apiSchemasRemovedCheckId),
			Operation: "Unknown",
			Path:      deletedSchema,
			Source:    "components.schemas." + deletedSchema,
		})
	}
	return result
}
