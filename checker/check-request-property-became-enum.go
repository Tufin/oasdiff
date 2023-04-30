package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const requestPropertyBecameEnumId = "request-property-became-enum"
const requestBodyBecameEnumId = "request-body-became-enum"

func RequestPropertyBecameEnumCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]

			if operationItem.RequestBodyDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff == nil ||
				operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified == nil {
				continue
			}
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				schemaDiff := mediaTypeDiff.SchemaDiff

				if !schemaDiff.TypeDiff.Empty() {
					// type changed will trigger another kind of breaking-change
				} else {
					if schemaDiff.EnumDiff != nil && schemaDiff.EnumDiff.EnumAdded {
						result = append(result, BackwardCompatibilityError{
							Id:        requestBodyBecameEnumId,
							Level:     ERR,
							Text:      config.i18n(requestBodyBecameEnumId),
							Operation: operation,
							Path:      path,
							Source:    source,
						})
					}
				}

				CheckModifiedPropertiesDiff(
					schemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {

						if !propertyDiff.TypeDiff.Empty() {
							// type changed will trigger another kind of breaking-change
							return
						}

						if enumDiff := propertyDiff.EnumDiff; enumDiff == nil || !enumDiff.EnumAdded {
							return
						}

						result = append(result, BackwardCompatibilityError{
							Id:        requestPropertyBecameEnumId,
							Level:     ERR,
							Text:      fmt.Sprintf(config.i18n(requestPropertyBecameEnumId), ColorizedValue(propertyFullName(propertyPath))),
							Operation: operation,
							Path:      path,
							Source:    source,
						})
					})
			}
		}
	}
	return result
}
