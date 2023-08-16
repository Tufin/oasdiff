package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyMediaTypeAdded   = "request-body-media-type-added"
	RequestBodyMediaTypeRemoved = "request-body-media-type-removed"
)

func RequestBodyMediaTypeChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			source := (*operationsSources)[operationItem.Revision]

			addedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeAdded
			for _, mediaType := range addedMediaTypes {
				result = append(result, ApiChange{
					Id:          RequestBodyMediaTypeAdded,
					Level:       INFO,
					Text:        config.Localize(RequestBodyMediaTypeAdded, mediaType),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}

			removedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeDeleted
			for _, mediaType := range removedMediaTypes {
				result = append(result, ApiChange{
					Id:          RequestBodyMediaTypeRemoved,
					Level:       ERR,
					Text:        config.Localize(RequestBodyMediaTypeRemoved, mediaType),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}

		}
	}
	return result
}
