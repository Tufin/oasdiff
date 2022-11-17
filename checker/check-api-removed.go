package checker

import (
	"fmt"
	"time"

	"cloud.google.com/go/civil"
	"github.com/tufin/oasdiff/diff"
)

func APIRemovedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) []BackwardCompatibilityError {
	result := make([]BackwardCompatibilityError, 0)
	if diffReport.PathsDiff == nil {
		return result
	}
	
	for _, path := range diffReport.PathsDiff.Deleted {
		if diffReport.PathsDiff.Base[path] == nil || diffReport.PathsDiff.Base[path].Operations()  == nil {
			continue
		}
		for operation, _ := range diffReport.PathsDiff.Base[path].Operations() {
			op := diffReport.PathsDiff.Base[path].Operations()[operation]
			if !op.Deprecated {
				source := "original_source=" + (*operationsSources)[op]
				result = append(result, BackwardCompatibilityError{
					Id:        "api-path-removed-without-deprecation",
					Level:     ERR,
					Text:      "api path removed without deprecation",
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			date, err := diff.GetSunsetDate(op.ExtensionProps)
			if err != nil {
				source := "original_source=" + (*operationsSources)[op]
				result = append(result, BackwardCompatibilityError{
					Id:        "api-path-sunset-parse",
					Level:     ERR,
					Text:      "api path sunset date can't be parsed",
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			if !civil.DateOf(time.Now()).After(date) {
				source := (*operationsSources)[op]
				result = append(result, BackwardCompatibilityError{
					Id:        "api-path-removed-before-sunset",
					Level:     ERR,
					Text:      fmt.Sprintf("api path removed before the sunset date %s", date),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
			}
		}
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for _, operation := range pathItem.OperationsDiff.Deleted {
			op := pathItem.Base.Operations()[operation]
			if !op.Deprecated {
				source := (*operationsSources)[op]
				result = append(result, BackwardCompatibilityError{
					Id:        "api-removed-without-deprecation",
					Level:     ERR,
					Text:      "api removed without deprecation",
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			date, err := diff.GetSunsetDate(op.ExtensionProps)
			if err != nil {
				source := (*operationsSources)[op]
				result = append(result, BackwardCompatibilityError{
					Id:        "api-sunset-parse",
					Level:     ERR,
					Text:      "api sunset date can't be parsed",
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
				continue
			}
			if !civil.DateOf(time.Now()).After(date) {
				source := (*operationsSources)[op]
				result = append(result, BackwardCompatibilityError{
					Id:        "api-removed-before-sunset",
					Level:     ERR,
					Text:      fmt.Sprintf("api removed before the sunset date %s", date),
					Operation: operation,
					Path:      path,
					Source:    source,
					ToDo:      "Add to exceptions-list.md",
				})
			}
		}
	}

	return result
}
