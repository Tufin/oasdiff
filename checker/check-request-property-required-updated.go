package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestPropertyBecameRequiredId = "request-property-became-required"
	RequestPropertyBecameOptionalId = "request-property-became-optional"
)

func RequestPropertyRequiredUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

				if mediaTypeDiff.SchemaDiff.RequiredDiff != nil {
					for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Added {
						if mediaTypeDiff.SchemaDiff.Base.Properties[changedRequiredPropertyName] == nil {
							// it is a new property, checked by the new-required-request-property check
							continue
						}
						if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName] == nil {
							// property was removed, checked by request-property-removed
							continue
						}
						if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName].Value.ReadOnly {
							continue
						}
						result = append(result, ApiChange{
							Id:          RequestPropertyBecameRequiredId,
							Level:       ERR,
							Text:        config.Localize(RequestPropertyBecameRequiredId, ColorizedValue(changedRequiredPropertyName)),
							Args:        []any{changedRequiredPropertyName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
					for _, changedRequiredPropertyName := range mediaTypeDiff.SchemaDiff.RequiredDiff.Deleted {
						if mediaTypeDiff.SchemaDiff.Base.Properties[changedRequiredPropertyName] == nil {
							// it is a new property, checked by the new-required-request-property check
							continue
						}
						if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName] == nil {
							// property was removed, checked by request-property-removed
							continue
						}
						if mediaTypeDiff.SchemaDiff.Revision.Properties[changedRequiredPropertyName].Value.ReadOnly {
							continue
						}
						result = append(result, ApiChange{
							Id:          RequestPropertyBecameOptionalId,
							Level:       INFO,
							Text:        config.Localize(RequestPropertyBecameOptionalId, ColorizedValue(changedRequiredPropertyName)),
							Args:        []any{changedRequiredPropertyName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}

				}

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						requiredDiff := propertyDiff.RequiredDiff
						if requiredDiff == nil {
							return
						}
						for _, changedRequiredPropertyName := range requiredDiff.Added {
							if propertyDiff.Revision.Properties[changedRequiredPropertyName] == nil {
								continue
							}
							if propertyDiff.Revision.Properties[changedRequiredPropertyName].Value.ReadOnly {
								continue
							}
							if propertyDiff.Base.Properties[changedRequiredPropertyName] == nil {
								// it is a new property, checked by the new-required-request-property check
								continue
							}

							propName := propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName))

							result = append(result, ApiChange{
								Id:          RequestPropertyBecameRequiredId,
								Level:       ERR,
								Text:        config.Localize(RequestPropertyBecameRequiredId, ColorizedValue(propName)),
								Args:        []any{propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}

						for _, changedRequiredPropertyName := range requiredDiff.Deleted {
							if propertyDiff.Revision.Properties[changedRequiredPropertyName] == nil {
								continue
							}
							if propertyDiff.Revision.Properties[changedRequiredPropertyName].Value.ReadOnly {
								continue
							}
							if propertyDiff.Base.Properties[changedRequiredPropertyName] == nil {
								// it is a new property, checked by the new-required-request-property check
								continue
							}

							propName := propertyFullName(propertyPath, propertyFullName(propertyName, changedRequiredPropertyName))

							result = append(result, ApiChange{
								Id:          RequestPropertyBecameOptionalId,
								Level:       INFO,
								Text:        config.Localize(RequestPropertyBecameOptionalId, ColorizedValue(propName)),
								Args:        []any{propName},
								Operation:   operation,
								OperationId: operationItem.Revision.OperationID,
								Path:        path,
								Source:      source,
							})
						}
					})
			}
		}
	}
	return result
}
