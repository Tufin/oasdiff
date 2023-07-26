package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const requestHeaderPropertyBecameEnumId = "request-header-property-became-enum"

func RequestHeaderPropertyBecameEnumCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			source := (*operationsSources)[operationItem.Revision]

			if operationItem.ParametersDiff == nil {
				continue
			}

			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {

				if paramLocation != "header" {
					continue
				}

				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil {
						continue
					}

					if paramDiff.SchemaDiff.EnumDiff != nil && paramDiff.SchemaDiff.EnumDiff.EnumAdded {
						result = append(result, ApiChange{
							Id:          requestHeaderPropertyBecameEnumId,
							Level:       ERR,
							Text:        fmt.Sprintf(config.i18n(requestHeaderPropertyBecameEnumId), colorizedValue(paramName)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

					checkModifiedPropertiesDiff(
						paramDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {

							if enumDiff := propertyDiff.EnumDiff; enumDiff == nil || !enumDiff.EnumAdded {
								return
							}

							result = append(result, ApiChange{
								Id:          requestHeaderPropertyBecameEnumId,
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n(requestHeaderPropertyBecameEnumId), colorizedValue(paramName), colorizedValue(propertyFullName(propertyPath, propertyName))),
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						})
				}
			}
		}
	}
	return result
}
