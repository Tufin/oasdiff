package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	APITagRemovedId = "api-tag-removed"
	APITagAddedId   = "api-tag-added"
)

func APITagUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			op := pathItem.Base.GetOperation(operation)

			if operationItem.TagsDiff == nil {
				continue
			}

			for _, tag := range operationItem.TagsDiff.Deleted {
				result = append(result, NewApiChange(
					APITagRemovedId,
					config,
					[]any{tag},
					"",
					operationsSources,
					op,
					operation,
					path,
				))
			}

			for _, tag := range operationItem.TagsDiff.Added {
				result = append(result, NewApiChange(
					APITagAddedId,
					config,
					[]any{tag},
					"",
					operationsSources,
					op,
					operation,
					path,
				))
			}
		}
	}
	return result
}
