package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"golang.org/x/exp/slices"
)

const (
	NewRequiredRequestHeaderPropertyId = "new-required-request-header-property"
)

func NewRequiredRequestHeaderPropertyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
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
					CheckAddedPropertiesDiff(
						paramDiff.SchemaDiff,
						func(propertyPath string, newPropertyName string, newProperty *openapi3.Schema, parent *diff.SchemaDiff) {
							if newProperty.ReadOnly {
								return
							}
							if !slices.Contains(parent.Revision.Required, newPropertyName) {
								return
							}

							result = append(result, ApiChange{
								Id:          NewRequiredRequestHeaderPropertyId,
								Level:       ERR,
								Args:        []any{paramName, propertyFullName(propertyPath, newPropertyName)},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      load.NewSource(source),
							})
						})
				}
			}
		}
	}
	return result
}
