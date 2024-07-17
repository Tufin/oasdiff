package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterDefaultValueChangedId = "request-parameter-default-value-changed"
	RequestParameterDefaultValueAddedId   = "request-parameter-default-value-added"
	RequestParameterDefaultValueRemovedId = "request-parameter-default-value-removed"
)

func RequestParameterDefaultValueChangedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {
			if operationItem.ParametersDiff == nil {
				continue
			}
			appendResultItem := func(messageId string, a ...any) {
				result = append(result, NewApiChange(
					messageId,
					config,
					a,
					"",
					operationsSources,
					operationItem.Revision,
					operation,
					path,
				))
			}
			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {

					baseParam := operationItem.Base.Parameters.GetByInAndName(paramLocation, paramName)
					if baseParam == nil || baseParam.Required {
						continue
					}

					revisionParam := operationItem.Revision.Parameters.GetByInAndName(paramLocation, paramName)
					if revisionParam == nil || revisionParam.Required {
						continue
					}

					if paramDiff.SchemaDiff == nil {
						continue
					}

					defaultValueDiff := paramDiff.SchemaDiff.DefaultDiff
					if defaultValueDiff.Empty() {
						continue
					}

					if defaultValueDiff.From == nil {
						appendResultItem(RequestParameterDefaultValueAddedId, paramLocation, paramName, defaultValueDiff.To)
					} else if defaultValueDiff.To == nil {
						appendResultItem(RequestParameterDefaultValueRemovedId, paramLocation, paramName, defaultValueDiff.From)
					} else {
						appendResultItem(RequestParameterDefaultValueChangedId, paramLocation, paramName, defaultValueDiff.From, defaultValueDiff.To)
					}
				}
			}
		}
	}
	return result
}
