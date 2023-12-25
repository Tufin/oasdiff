package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	APITagRemovedId = "api-tag-removed"
	APITagAddedId   = "api-tag-added"
)

func APITagUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	changeGetter := newApiChangeGetter(config, operationsSources)
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			op := pathItem.Base.Operations()[operation]

			if operationItem.TagsDiff == nil {
				continue
			}

			for _, tag := range operationItem.TagsDiff.Deleted {
				result = append(result, changeGetter(
					APITagRemovedId,
					INFO,
					[]any{tag},
					"",
					operation,
					op,
					path,
					op,
				))
			}

			for _, tag := range operationItem.TagsDiff.Added {
				result = append(result, changeGetter(
					APITagAddedId,
					INFO,
					[]any{tag},
					"",
					operation,
					op,
					path,
					op,
				))
			}
		}
	}
	return result
}
