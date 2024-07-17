package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	ResponseMediaTypeRemovedId = "response-media-type-removed"
	ResponseMediaTypeAddedId   = "response-media-type-added"
)

func ResponseMediaTypeUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil {
				continue
			}
			if operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			for responseStatus, responsesDiff := range operationItem.ResponsesDiff.Modified {
				if responsesDiff.ContentDiff == nil {
					continue
				}
				if responsesDiff.ContentDiff.MediaTypeDeleted == nil {
					continue
				}
				for _, mediaType := range responsesDiff.ContentDiff.MediaTypeDeleted {
					result = append(result, NewApiChange(
						ResponseMediaTypeRemovedId,
						config,
						[]any{mediaType, responseStatus},
						"",
						operationsSources,
						operationItem.Revision,
						operation,
						path,
					))
				}
				for _, mediaType := range responsesDiff.ContentDiff.MediaTypeAdded {
					result = append(result, NewApiChange(
						ResponseMediaTypeAddedId,
						config,
						[]any{mediaType, responseStatus},
						"",
						operationsSources,
						operationItem.Revision,
						operation,
						path,
					))
				}
			}
		}
	}
	return result
}
