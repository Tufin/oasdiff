package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestBodyDiscriminatorAddedId                   = "request-body-discriminator-added"
	RequestBodyDiscriminatorRemovedId                 = "request-body-discriminator-removed"
	RequestBodyDiscriminatorPropertyNameChangedId     = "request-body-discriminator-property-name-changed"
	RequestBodyDiscriminatorMappingAddedId            = "request-body-discriminator-mapping-added"
	RequestBodyDiscriminatorMappingDeletedId          = "request-body-discriminator-mapping-deleted"
	RequestBodyDiscriminatorMappingChangedId          = "request-body-discriminator-mapping-changed"
	RequestPropertyDiscriminatorAddedId               = "request-property-discriminator-added"
	RequestPropertyDiscriminatorRemovedId             = "request-property-discriminator-removed"
	RequestPropertyDiscriminatorPropertyNameChangedId = "request-property-discriminator-property-name-changed"
	RequestPropertyDiscriminatorMappingAddedId        = "request-property-discriminator-mapping-added"
	RequestPropertyDiscriminatorMappingDeletedId      = "request-property-discriminator-mapping-deleted"
	RequestPropertyDiscriminatorMappingChangedId      = "request-property-discriminator-mapping-changed"
)

func RequestDiscriminatorUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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

			appendResultItem := func(messageId string, a ...any) {
				result = append(result, ApiChange{
					Id:          messageId,
					Level:       INFO,
					Text:        config.Localize(messageId, a...),
					Args:        []any{},
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}

			for _, mediaTypeDiff := range operationItem.RequestBodyDiff.ContentDiff.MediaTypeModified {
				if mediaTypeDiff.SchemaDiff == nil {
					continue
				}

				processDiscriminatorDiffForRequest(
					mediaTypeDiff.SchemaDiff.DiscriminatorDiff,
					"",
					appendResultItem)

				CheckModifiedPropertiesDiff(
					mediaTypeDiff.SchemaDiff,
					func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
						processDiscriminatorDiffForRequest(
							propertyDiff.DiscriminatorDiff,
							propertyFullName(propertyPath, propertyName),
							appendResultItem)
					})

			}
		}
	}
	return result
}

func processDiscriminatorDiffForRequest(
	discriminatorDiff *diff.DiscriminatorDiff,
	propertyName string,
	appendResultItem func(messageId string, a ...any)) {

	if discriminatorDiff == nil {
		return
	}

	messageIdPrefix := "request-body-discriminator"
	if propertyName != "" {
		messageIdPrefix = "request-property-discriminator"
	}

	if discriminatorDiff.Added {
		if propertyName == "" {
			appendResultItem(messageIdPrefix + "-added")
		} else {
			appendResultItem(messageIdPrefix+"-added", ColorizedValue(propertyName))
		}
		return
	}
	if discriminatorDiff.Deleted {
		if propertyName == "" {
			appendResultItem(messageIdPrefix + "-removed")
		} else {
			appendResultItem(messageIdPrefix+"-removed", ColorizedValue(propertyName))
		}
		return
	}

	if discriminatorDiff.PropertyNameDiff != nil {
		if propertyName == "" {
			appendResultItem(messageIdPrefix+"-property-name-changed",
				ColorizedValue(discriminatorDiff.PropertyNameDiff.From),
				ColorizedValue(discriminatorDiff.PropertyNameDiff.To))
		} else {
			appendResultItem(messageIdPrefix+"-property-name-changed",
				ColorizedValue(propertyName),
				ColorizedValue(discriminatorDiff.PropertyNameDiff.From),
				ColorizedValue(discriminatorDiff.PropertyNameDiff.To))
		}
	}

	if discriminatorDiff.MappingDiff != nil {
		if len(discriminatorDiff.MappingDiff.Added) > 0 {
			if propertyName == "" {
				appendResultItem(messageIdPrefix+"-mapping-added",
					ColorizedValue(discriminatorDiff.MappingDiff.Added))
			} else {
				appendResultItem(messageIdPrefix+"-mapping-added",
					ColorizedValue(discriminatorDiff.MappingDiff.Added),
					ColorizedValue(propertyName))
			}
		}

		if len(discriminatorDiff.MappingDiff.Deleted) > 0 {
			if propertyName == "" {
				appendResultItem(messageIdPrefix+"-mapping-deleted",
					ColorizedValue(discriminatorDiff.MappingDiff.Deleted))
			} else {
				appendResultItem(messageIdPrefix+"-mapping-deleted",
					ColorizedValue(discriminatorDiff.MappingDiff.Deleted),
					ColorizedValue(propertyName))
			}
		}

		for k, v := range discriminatorDiff.MappingDiff.Modified {
			if propertyName == "" {
				appendResultItem(messageIdPrefix+"-mapping-changed",
					ColorizedValue(k),
					ColorizedValue(v.From),
					ColorizedValue(v.To))
			} else {
				appendResultItem(messageIdPrefix+"-mapping-changed",
					ColorizedValue(k),
					ColorizedValue(v.From),
					ColorizedValue(v.To),
					ColorizedValue(propertyName))

			}
		}
	}
}
