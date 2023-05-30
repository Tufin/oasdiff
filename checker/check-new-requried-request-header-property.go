package checker

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func NewRequiredRequestHeaderPropertyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
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

			if operationItem.ParametersDiff == nil {
				continue
			}

			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {

				if paramLocation != "header" {
					continue
				}

				for paramName, paramDiff := range paramDiffs {
					CheckAddedPropertiesDiff(
						paramDiff.SchemaDiff,
						func(propertyPath string, newPropertyName string, newProperty *openapi3.Schema, parent *diff.SchemaDiff) {
							if newProperty.ReadOnly {
								return
							}
							if !slices.Contains(parent.Revision.Value.Required, newPropertyName) {
								return
							}

							result = append(result, BackwardCompatibilityError{
								Id:          "new-required-request-header-property",
								Level:       ERR,
								Text:        fmt.Sprintf(config.i18n("new-required-request-header-property"), ColorizedValue(paramName), ColorizedValue(propertyFullName(propertyPath, newPropertyName))),
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
