package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
	"golang.org/x/exp/slices"
)

func NewRequiredRequestPropertyCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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
			modifiedMediaTypes := operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified
			for _, mediaTypeDiff := range modifiedMediaTypes {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}
				if mediaTypeDiff.SchemaDiff.PropertiesDiff == nil {
					continue
				}

				for _, topPropertyName := range mediaTypeDiff.SchemaDiff.PropertiesDiff.Added {
					propertyName := topPropertyName
					propertyItem := mediaTypeDiff.SchemaDiff.Revision.Value.Properties[topPropertyName].Value
					parent := mediaTypeDiff.SchemaDiff.Revision.Value
					if !propertyItem.ReadOnly &&
						slices.Contains(parent.Required, propertyName) {
						source := (*operationsSources)[operationItem.Revision]
						result = append(result, BackwardCompatibilityError{
							Id:        "new-required-request-property",
							Level:     ERR,
							Text:      fmt.Sprintf("added new required request property %s", ColorizedValue(propertyName)),
							Operation: operation,
							Path:      path,
							Source:    source,
							ToDo:      "Add to exceptions-list.md",
						})
					}
				}

				for topPropertyName, topPropertyDiff  := range mediaTypeDiff.SchemaDiff.PropertiesDiff.Modified {
					processModifiedPropertiesDiff(
						"",
						topPropertyName,
						topPropertyDiff,
						nil,
						func(propertyPath string, propertyName string, propertyItem *diff.SchemaDiff, parent *diff.SchemaDiff) {
							if propertyItem.PropertiesDiff == nil {
								return
							}
							if propertyItem.PropertiesDiff.Added == nil {
								return
							}
			
							for _, newPropertyName := range propertyItem.PropertiesDiff.Added {
								newPropertyItem := propertyItem.Revision.Value.Properties[newPropertyName].Value
								newParent := propertyItem.Revision.Value
								if !newPropertyItem.ReadOnly &&
									slices.Contains(newParent.Required, newPropertyName) {
									source := (*operationsSources)[operationItem.Revision]
									result = append(result, BackwardCompatibilityError{
										Id:        "new-required-request-property",
										Level:     ERR,
										Text:      fmt.Sprintf("added new required request property %s", ColorizedValue(propertyFullName(propertyPath, propertyName, newPropertyName))),
										Operation: operation,
										Path:      path,
										Source:    source,
										ToDo:      "Add to exceptions-list.md",
									})
								}
							}
						})
				}
			}
		}
	}
	return result
}
