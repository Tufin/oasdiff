package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	APISchemasRemovedId = "api-schema-removed"
)

func APIComponentsSchemaRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.ComponentsDiff.SchemasDiff == nil {
		return result
	}

	for _, deletedSchema := range diffReport.ComponentsDiff.SchemasDiff.Deleted {
		result = append(result, ComponentChange{
			Id:     APISchemasRemovedId,
			Level:  config.getLogLevel(APISchemasRemovedId, INFO),
			Text:   config.Localize(APISchemasRemovedId, ColorizedValue(deletedSchema)),
			Source: "Components",
		})
	}
	return result
}
