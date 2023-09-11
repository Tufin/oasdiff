package checker

import (
	"github.com/tufin/oasdiff/diff"
)

func NewNonPathRequestDefaultParameterCheck(diffReport *diff.Diff, _ *diff.OperationsSourcesMap, config Config) Changes {
	result := make(Changes, 0)
	if diffReport.PathsDiff == nil || len(diffReport.PathsDiff.Modified) == 0 {
		return result
	}
	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.ParametersDiff == nil {
			continue
		}
		for paramLoc, paramNameList := range pathItem.ParametersDiff.Added {
			if paramLoc == "path" {
				continue
			}
			for _, param := range pathItem.Revision.Parameters {
				if !paramNameList.Contains(param.Value.Name) {
					continue
				}
				id := "new-required-request-default-parameter-to-existing-path"
				level := ERR
				if !param.Value.Required {
					id = "new-optional-request-default-parameter-to-existing-path"
					level = INFO
				}
				result = append(result, ApiChange{
					Id:    id,
					Level: level,
					Text:  config.Localize(id, ColorizedValue(paramLoc), ColorizedValue(param.Value.Name)),
					Path:  path,
				})
			}
		}
	}
	return result
}
