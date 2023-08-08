package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterDefaultValueChanged(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config Config) Changes {
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
			source := (*operationsSources)[operationItem.Revision]
			appendResultItem := func(messageId string, a ...any) {
				result = append(result, ApiChange{
					Id:          messageId,
					Level:       ERR,
					Text:        fmt.Sprintf(config.i18n(messageId), a...),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
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
						appendResultItem("request-parameter-default-value-added", ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(defaultValueDiff.To))
					} else if defaultValueDiff.To == nil {
						appendResultItem("request-parameter-default-value-removed", ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(defaultValueDiff.From))
					} else {
						appendResultItem("request-parameter-default-value-changed", ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(defaultValueDiff.From), ColorizedValue(defaultValueDiff.To))
					}
				}
			}
		}
	}
	return result
}
