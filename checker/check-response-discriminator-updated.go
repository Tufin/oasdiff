package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func ResponseDiscriminatorUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}

		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ResponsesDiff == nil || operationItem.ResponsesDiff.Modified == nil {
				continue
			}
			source := (*operationsSources)[operationItem.Revision]

			appendResultItem := func(messageId string, a ...any) {
				result = append(result, ApiChange{
					Id:          messageId,
					Level:       INFO,
					Text:        config.Localize(messageId, a...),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}

			for responseStatus, responsesDiff := range operationItem.ResponsesDiff.Modified {
				if responsesDiff.ContentDiff == nil || responsesDiff.ContentDiff.MediaTypeModified == nil {
					continue
				}

				modifiedMediaTypes := responsesDiff.ContentDiff.MediaTypeModified
				for _, mediaTypeDiff := range modifiedMediaTypes {
					if mediaTypeDiff.SchemaDiff == nil {
						continue
					}

					processDiscriminatorDiff(
						mediaTypeDiff.SchemaDiff.DiscriminatorDiff,
						responseStatus,
						"",
						appendResultItem)

					CheckModifiedPropertiesDiff(
						mediaTypeDiff.SchemaDiff,
						func(propertyPath string, propertyName string, propertyDiff *diff.SchemaDiff, parent *diff.SchemaDiff) {
							processDiscriminatorDiff(
								propertyDiff.DiscriminatorDiff,
								responseStatus,
								propertyFullName(propertyPath, propertyName),
								appendResultItem)
						})
				}
			}
		}
	}
	return result
}

func processDiscriminatorDiff(
	discriminatorDiff *diff.DiscriminatorDiff,
	responseStatus string,
	propertyName string,
	appendResultItem func(messageId string, a ...any)) {

	if discriminatorDiff == nil {
		return
	}

	messageIdPrefix := "response-body-discriminator"
	if propertyName != "" {
		messageIdPrefix = "response-property-discriminator"
	}

	if discriminatorDiff.Added {
		if propertyName == "" {
			appendResultItem(messageIdPrefix+"-added", responseStatus)
		} else {
			appendResultItem(messageIdPrefix+"-added", ColorizedValue(propertyName), responseStatus)
		}
		return
	}
	if discriminatorDiff.Deleted {
		if propertyName == "" {
			appendResultItem(messageIdPrefix+"-removed", responseStatus)
		} else {
			appendResultItem(messageIdPrefix+"-removed", ColorizedValue(propertyName), responseStatus)
		}
		return
	}

	if discriminatorDiff.PropertyNameDiff != nil {
		if propertyName == "" {
			appendResultItem(messageIdPrefix+"-property-name-changed",
				ColorizedValue(discriminatorDiff.PropertyNameDiff.From),
				ColorizedValue(discriminatorDiff.PropertyNameDiff.To),
				responseStatus)
		} else {
			appendResultItem(messageIdPrefix+"-property-name-changed",
				ColorizedValue(propertyName),
				ColorizedValue(discriminatorDiff.PropertyNameDiff.From),
				ColorizedValue(discriminatorDiff.PropertyNameDiff.To),
				responseStatus)
		}
	}

	if discriminatorDiff.MappingDiff != nil {
		if len(discriminatorDiff.MappingDiff.Added) > 0 {
			if propertyName == "" {
				appendResultItem(messageIdPrefix+"-mapping-added",
					ColorizedValue(discriminatorDiff.MappingDiff.Added),
					responseStatus)
			} else {
				appendResultItem(messageIdPrefix+"-mapping-added",
					ColorizedValue(discriminatorDiff.MappingDiff.Added),
					ColorizedValue(propertyName),
					responseStatus)
			}
		}

		if len(discriminatorDiff.MappingDiff.Deleted) > 0 {
			if propertyName == "" {
				appendResultItem(messageIdPrefix+"-mapping-deleted",
					ColorizedValue(discriminatorDiff.MappingDiff.Deleted),
					responseStatus)
			} else {
				appendResultItem(messageIdPrefix+"-mapping-deleted",
					ColorizedValue(discriminatorDiff.MappingDiff.Deleted),
					ColorizedValue(propertyName),
					responseStatus)
			}
		}

		for k, v := range discriminatorDiff.MappingDiff.Modified {
			if propertyName == "" {
				appendResultItem(messageIdPrefix+"-mapping-changed",
					ColorizedValue(k),
					ColorizedValue(v.From),
					ColorizedValue(v.To),
					responseStatus)
			} else {
				appendResultItem(messageIdPrefix+"-mapping-changed",
					ColorizedValue(k),
					ColorizedValue(v.From),
					ColorizedValue(v.To),
					ColorizedValue(propertyName),
					responseStatus)

			}
		}
	}
}
