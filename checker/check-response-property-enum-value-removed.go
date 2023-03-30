package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const responsePropertyEnumValueRemovedId = "response-property-enum-value-removed"

func ResponseParameterEnumValueRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
			source := (*operationsSources)[operationItem.Revision]
			for _, responseItems := range operationItem.ResponsesDiff.Modified {
				for _, mediaTypeItem := range responseItems.ContentDiff.MediaTypeModified {
					if mediaTypeItem.SchemaDiff == nil {
						continue
					}

					if mediaTypeItem.SchemaDiff.PropertiesDiff == nil {
						continue
					}

					for property, propertyItem := range mediaTypeItem.SchemaDiff.PropertiesDiff.Modified {
						if propertyItem.EnumDiff == nil {
							continue
						}

						for _, enumVal := range propertyItem.EnumDiff.Deleted {
							result = append(result, BackwardCompatibilityError{
								Id:        responsePropertyEnumValueRemovedId,
								Level:     ERR,
								Text:      fmt.Sprintf(config.i18n("response-property-enum-value-removed"), property, ColorizedValue(enumVal)),
								Operation: operation,
								Path:      path,
								Source:    source,
							})
						}
					}
				}

			}
		}
	}
	return result
}
