package checker

import (
	"github.com/tufin/oasdiff/diff"
)

const (
	RequestParameterMaxItemsIncreasedId = "request-parameter-max-items-increased"
	RequestParameterMaxItemsDecreasedId = "request-parameter-max-items-decreased"
)

func RequestParameterMaxItemsUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	changeGetter := newApiChangeGetter(config, operationsSources)
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
			for paramLocation, paramDiffs := range operationItem.ParametersDiff.Modified {
				for paramName, paramDiff := range paramDiffs {
					if paramDiff.SchemaDiff == nil || paramDiff.SchemaDiff.ItemsDiff == nil {
						continue
					}
					maxItemsDiff := paramDiff.SchemaDiff.ItemsDiff.MaxItemsDiff
					if maxItemsDiff == nil {
						continue
					}
					if maxItemsDiff.From == nil ||
						maxItemsDiff.To == nil {
						continue
					}

					id := RequestParameterMaxItemsDecreasedId
					level := ERR
					if IsIncreasedValue(maxItemsDiff) {
						id = RequestParameterMaxItemsIncreasedId
						level = INFO
					}

					result = append(result, changeGetter(
						id,
						level,
						[]any{paramLocation, paramName, maxItemsDiff.From, maxItemsDiff.To},
						"",
						operation,
						operationItem.Revision,
						path,
						operationItem.Revision,
					))
				}
			}
		}
	}
	return result
}
