package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	APISecurityRemovedCheckId       = "api-security-removed"
	APISecurityAddedCheckId         = "api-security-added"
	APISecurityScopeAddedId         = "api-security-scope-added"
	APISecurityScopeRemovedId       = "api-security-scope-removed"
	APIGlobalSecurityRemovedCheckId = "api-global-security-removed"
	APIGlobalSecurityAddedCheckId   = "api-global-security-added"
	APIGlobalSecurityScopeAddedId   = "api-global-security-scope-added"
	APIGlobalSecurityScopeRemovedId = "api-global-security-scope-removed"
)

func checkGlobalSecurity(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) IBackwardCompatibilityErrors {
	result := make(IBackwardCompatibilityErrors, 0)
	if diffReport.SecurityDiff == nil {
		return result
	}

	for _, addedSecurity := range diffReport.SecurityDiff.Added {
		result = append(result, BackwardCompatibilityComponentError{
			Id:     APIGlobalSecurityAddedCheckId,
			Level:  INFO,
			Text:   fmt.Sprintf(config.i18n(APIGlobalSecurityAddedCheckId), ColorizedValue(addedSecurity)),
			Source: "security." + addedSecurity,
		})
	}

	for _, removedSecurity := range diffReport.SecurityDiff.Deleted {
		result = append(result, BackwardCompatibilityError{
			Id:          APIGlobalSecurityRemovedCheckId,
			Level:       INFO,
			Text:        fmt.Sprintf(config.i18n(APIGlobalSecurityRemovedCheckId), ColorizedValue(removedSecurity)),
			Operation:   "N/A",
			Path:        "",
			Source:      "security." + removedSecurity,
			OperationId: "N/A",
		})
	}

	for _, updatedSecurity := range diffReport.SecurityDiff.Modified {
		for securitySchemeName, updatedSecuritySchemeScopes := range updatedSecurity {
			for _, addedScope := range updatedSecuritySchemeScopes.Added {
				result = append(result, BackwardCompatibilityError{
					Id:          APIGlobalSecurityScopeAddedId,
					Level:       INFO,
					Text:        fmt.Sprintf(config.i18n(APIGlobalSecurityScopeAddedId), ColorizedValue(addedScope), ColorizedValue(securitySchemeName)),
					Operation:   "N/A",
					Path:        "",
					Source:      "security.scopes." + addedScope,
					OperationId: "N/A",
				})
			}
			for _, deletedScope := range updatedSecuritySchemeScopes.Deleted {
				result = append(result, BackwardCompatibilityError{
					Id:          APIGlobalSecurityScopeRemovedId,
					Level:       INFO,
					Text:        fmt.Sprintf(config.i18n(APIGlobalSecurityScopeRemovedId), ColorizedValue(deletedScope), ColorizedValue(securitySchemeName)),
					Operation:   "N/A",
					Path:        "",
					Source:      "security.scopes." + deletedScope,
					OperationId: "N/A",
				})
			}
		}
	}

	return result

}

func APISecurityUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) IBackwardCompatibilityErrors {
	result := make(IBackwardCompatibilityErrors, 0)

	result = append(result, checkGlobalSecurity(diffReport, operationsSources, config)...)

	if diffReport.PathsDiff == nil || diffReport.PathsDiff.Modified == nil {
		return result
	}

	for path, pathItem := range diffReport.PathsDiff.Modified {
		if pathItem.OperationsDiff == nil {
			continue
		}
		for operation, operationItem := range pathItem.OperationsDiff.Modified {

			source := (*operationsSources)[operationItem.Revision]

			if operationItem.SecurityDiff == nil {
				continue
			}

			for _, addedSecurity := range operationItem.SecurityDiff.Added {
				if addedSecurity == "" {
					continue
				}
				result = append(result, BackwardCompatibilityError{
					Id:          APISecurityAddedCheckId,
					Level:       INFO,
					Text:        fmt.Sprintf(config.i18n(APISecurityAddedCheckId), ColorizedValue(addedSecurity)),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}

			for _, deletedSecurity := range operationItem.SecurityDiff.Deleted {
				if deletedSecurity == "" {
					continue
				}
				result = append(result, BackwardCompatibilityError{
					Id:          APISecurityRemovedCheckId,
					Level:       INFO,
					Text:        fmt.Sprintf(config.i18n(APISecurityRemovedCheckId), ColorizedValue(deletedSecurity)),
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      source,
				})
			}

			for _, updatedSecurity := range operationItem.SecurityDiff.Modified {
				if updatedSecurity.Empty() {
					continue
				}
				for securitySchemeName, updatedSecuritySchemeScopes := range updatedSecurity {
					for _, addedScope := range updatedSecuritySchemeScopes.Added {
						result = append(result, BackwardCompatibilityError{
							Id:          APISecurityScopeAddedId,
							Level:       INFO,
							Text:        fmt.Sprintf(config.i18n(APISecurityScopeAddedId), ColorizedValue(addedScope), ColorizedValue(securitySchemeName)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
					for _, deletedScope := range updatedSecuritySchemeScopes.Deleted {
						result = append(result, BackwardCompatibilityError{
							Id:          APISecurityScopeRemovedId,
							Level:       INFO,
							Text:        fmt.Sprintf(config.i18n(APISecurityScopeRemovedId), ColorizedValue(deletedScope), ColorizedValue(securitySchemeName)),
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      source,
						})
					}
				}
			}

		}
	}

	return result
}
