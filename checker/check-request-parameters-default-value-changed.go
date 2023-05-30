package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

func RequestParameterDefaultValueChanged(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
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

					source := (*operationsSources)[operationItem.Revision]

					result = append(result, BackwardCompatibilityError{
						Id:          "request-parameter-default-value-changed",
						Level:       ERR,
						Text:        fmt.Sprintf(config.i18n("request-parameter-default-value-changed"), ColorizedValue(paramLocation), ColorizedValue(paramName), ColorizedValue(defaultValueDiff.From), ColorizedValue(defaultValueDiff.To)),
						Operation:   operation,
						OperationId: operationItem.Revision.OperationID,
						Path:        path,
						Source:      source,
					})
				}
			}
		}
	}
	return result
}
