package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	apiSchemasRemovedCheckId = "api-schema-removed"
)

func APIComponentsSchemaRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) IBackwardCompatibilityErrors {
	result := make(IBackwardCompatibilityErrors, 0)
	if diffReport.ComponentsDiff.SchemasDiff == nil {
		return result
	}

	for _, deletedSchema := range diffReport.ComponentsDiff.SchemasDiff.Deleted {
		result = append(result, BackwardCompatibilityError{
			Id:        apiSchemasRemovedCheckId,
			Level:     config.getLogLevel(apiSchemasRemovedCheckId, INFO),
			Text:      fmt.Sprintf(config.i18n(apiSchemasRemovedCheckId), ColorizedValue(deletedSchema)),
			Operation: "N/A",
			Path:      "",
			Source:    "components.schemas." + deletedSchema, // TODO: get the file name
		})
	}
	return result
}
