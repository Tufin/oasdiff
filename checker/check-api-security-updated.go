package checker

import (
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
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

func checkGlobalSecurity(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)
	if diffReport.SecurityDiff == nil {
		return result
	}

	for _, addedSecurity := range diffReport.SecurityDiff.Added {
		result = append(result, SecurityChange{
			Id:    APIGlobalSecurityAddedCheckId,
			Level: INFO,
			Args:  []any{addedSecurity},
		})
	}

	for _, removedSecurity := range diffReport.SecurityDiff.Deleted {
		result = append(result, SecurityChange{
			Id:    APIGlobalSecurityRemovedCheckId,
			Level: INFO,
			Args:  []any{removedSecurity},
		})
	}

	for _, updatedSecurity := range diffReport.SecurityDiff.Modified {
		for securitySchemeName, updatedSecuritySchemeScopes := range updatedSecurity {
			for _, addedScope := range updatedSecuritySchemeScopes.Added {
				result = append(result, SecurityChange{
					Id:    APIGlobalSecurityScopeAddedId,
					Level: INFO,
					Args:  []any{addedScope, securitySchemeName},
				})
			}
			for _, deletedScope := range updatedSecuritySchemeScopes.Deleted {
				result = append(result, SecurityChange{
					Id:    APIGlobalSecurityScopeRemovedId,
					Level: INFO,
					Args:  []any{deletedScope, securitySchemeName},
				})
			}
		}
	}

	return result

}

func APISecurityUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config *Config) Changes {
	result := make(Changes, 0)

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
				result = append(result, ApiChange{
					Id:          APISecurityAddedCheckId,
					Level:       INFO,
					Args:        []any{addedSecurity},
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}

			for _, deletedSecurity := range operationItem.SecurityDiff.Deleted {
				if deletedSecurity == "" {
					continue
				}
				result = append(result, ApiChange{
					Id:          APISecurityRemovedCheckId,
					Level:       INFO,
					Args:        []any{deletedSecurity},
					Operation:   operation,
					OperationId: operationItem.Revision.OperationID,
					Path:        path,
					Source:      load.NewSource(source),
				})
			}

			for _, updatedSecurity := range operationItem.SecurityDiff.Modified {
				if updatedSecurity.Empty() {
					continue
				}
				for securitySchemeName, updatedSecuritySchemeScopes := range updatedSecurity {
					for _, addedScope := range updatedSecuritySchemeScopes.Added {
						result = append(result, ApiChange{
							Id:          APISecurityScopeAddedId,
							Level:       INFO,
							Args:        []any{addedScope, securitySchemeName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      load.NewSource(source),
						})
					}
					for _, deletedScope := range updatedSecuritySchemeScopes.Deleted {
						result = append(result, ApiChange{
							Id:          APISecurityScopeRemovedId,
							Level:       INFO,
							Args:        []any{deletedScope, securitySchemeName},
							Operation:   operation,
							OperationId: operationItem.Revision.OperationID,
							Path:        path,
							Source:      load.NewSource(source),
						})
					}
				}
			}

		}
	}

	return result
}
