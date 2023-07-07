package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	apiSchemasRemovedCheckId = "api-schema-removed"
)

func APIComponentsSchemaRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.ComponentsDiff.SchemasDiff == nil {
		return result
	}

	for _, deletedSchema := range diffReport.ComponentsDiff.SchemasDiff.Deleted {
		result = append(result, ComponentChange{
			Id:     apiSchemasRemovedCheckId,
			Level:  config.getLogLevel(apiSchemasRemovedCheckId, INFO),
			Text:   fmt.Sprintf(config.i18n(apiSchemasRemovedCheckId), ColorizedValue(deletedSchema)),
			Source: "components.schemas." + deletedSchema, // TODO: get the file name
		})
	}
	return result
}
