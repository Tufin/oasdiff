package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMediaTypeAddedId   = "request-body-media-type-added"
	RequestBodyMediaTypeRemovedId = "request-body-media-type-removed"
)

func RequestBodyMediaTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}

			addedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeAdded
			for _, mediaType := range addedMediaTypes {
				result = append(result, changeGetter(
					RequestBodyMediaTypeAddedId,
					INFO,
					[]any{mediaType},
					"",
					operation,
					operationItem.Revision,
					path,
					operationItem.Revision,
				))
			}

			removedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeDeleted
			for _, mediaType := range removedMediaTypes {
				result = append(result, changeGetter(
					RequestBodyMediaTypeRemovedId,
					ERR,
					[]any{mediaType},
					"",
					operation,
					operationItem.Revision,
					path,
					operationItem.Revision,
				))
			}
		}
	}
	return result
}
