package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/diff"
)

const (
	APIComponentsSecurityRemovedCheckId           = "api-security-component-removed"
	APIComponentsSecurityAddedCheckId             = "api-security-component-added"
	APIComponentsSecurityComponentOauthUrlUpdated = "api-security-component-oauth-url-changed"
	APIComponentsSecurityTyepUpdated              = "api-security-component-type-changed"
	APIComponentsSecurityOauthTokenUrlUpdated     = "api-security-component-oauth-token-url-changed"
	APIComponentSecurityOauthScopeAdded           = "api-security-component-oauth-scope-added"
	APIComponentSecurityOauthScopeRemoved         = "api-security-component-oauth-scope-removed"
	APIComponentSecurityOauthScopeUpdated         = "api-security-component-oauth-scope-changed"
)

func checkOAuthUpdates(updatedSecurity *diff.SecuritySchemeDiff, config BackwardCompatibilityCheckConfig, updatedSecurityName string) IBackwardCompatibilityErrors {
	result := make(IBackwardCompatibilityErrors, 0)

	if updatedSecurity.OAuthFlowsDiff == nil {
		return result
	}

	if updatedSecurity.OAuthFlowsDiff.ImplicitDiff == nil {
		return result
	}

	if urlDiff := updatedSecurity.OAuthFlowsDiff.ImplicitDiff.AuthorizationURLDiff; urlDiff != nil {
		result = append(result, BackwardCompatibilityError{
			Id:          APIComponentsSecurityComponentOauthUrlUpdated,
			Level:       INFO,
			Text:        fmt.Sprintf(config.i18n(APIComponentsSecurityComponentOauthUrlUpdated), ColorizedValue(updatedSecurityName), ColorizedValue(urlDiff.From), ColorizedValue(urlDiff.To)),
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		})
	}

	if tokenDiff := updatedSecurity.OAuthFlowsDiff.ImplicitDiff.TokenURLDiff; tokenDiff != nil {
		result = append(result, BackwardCompatibilityError{
			Id:          APIComponentsSecurityOauthTokenUrlUpdated,
			Level:       INFO,
			Text:        fmt.Sprintf(config.i18n(APIComponentsSecurityOauthTokenUrlUpdated), ColorizedValue(updatedSecurityName), ColorizedValue(tokenDiff.From), ColorizedValue(tokenDiff.To)),
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		})
	}

	if scopesDiff := updatedSecurity.OAuthFlowsDiff.ImplicitDiff.ScopesDiff; scopesDiff != nil {
		for _, addedScope := range scopesDiff.Added {
			result = append(result, BackwardCompatibilityError{
				Id:          APIComponentSecurityOauthScopeAdded,
				Level:       INFO,
				Text:        fmt.Sprintf(config.i18n(APIComponentSecurityOauthScopeAdded), ColorizedValue(updatedSecurityName), ColorizedValue(addedScope)),
				Operation:   "N/A",
				Path:        "N/A",
				Source:      "N/A",
				OperationId: "N/A",
			})
		}

		for _, removedScope := range scopesDiff.Deleted {
			result = append(result, BackwardCompatibilityError{
				Id:          APIComponentSecurityOauthScopeRemoved,
				Level:       INFO,
				Text:        fmt.Sprintf(config.i18n(APIComponentSecurityOauthScopeRemoved), ColorizedValue(updatedSecurityName), ColorizedValue(removedScope)),
				Operation:   "N/A",
				Path:        "N/A",
				Source:      "N/A",
				OperationId: "N/A",
			})
		}

		for name, modifiedScope := range scopesDiff.Modified {
			result = append(result, BackwardCompatibilityError{
				Id:          APIComponentSecurityOauthScopeUpdated,
				Level:       INFO,
				Text:        fmt.Sprintf(config.i18n(APIComponentSecurityOauthScopeUpdated), ColorizedValue(updatedSecurityName), ColorizedValue(name), ColorizedValue(modifiedScope.From), ColorizedValue(modifiedScope.To)),
				Operation:   "N/A",
				Path:        "N/A",
				Source:      "N/A",
				OperationId: "N/A",
			})
		}

	}

	return result
}

func APIComponentsSecurityUpdatedCheck(diffReport *diff.Diff, operationsSources *diff.OperationsSourcesMap, config BackwardCompatibilityCheckConfig) IBackwardCompatibilityErrors {
	result := make(IBackwardCompatibilityErrors, 0)
	if diffReport.ComponentsDiff.SecuritySchemesDiff == nil {
		return result
	}

	for _, updatedSecurity := range diffReport.ComponentsDiff.SecuritySchemesDiff.Added {
		result = append(result, BackwardCompatibilityError{
			Id:          APIComponentsSecurityAddedCheckId,
			Level:       INFO,
			Text:        fmt.Sprintf(config.i18n(APIComponentsSecurityAddedCheckId), ColorizedValue(updatedSecurity)),
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		})
	}

	for _, updatedSecurity := range diffReport.ComponentsDiff.SecuritySchemesDiff.Deleted {
		result = append(result, BackwardCompatibilityError{
			Id:          APIComponentsSecurityRemovedCheckId,
			Level:       INFO,
			Text:        fmt.Sprintf(config.i18n(APIComponentsSecurityRemovedCheckId), ColorizedValue(updatedSecurity)),
			Operation:   "N/A",
			Path:        "N/A",
			Source:      "N/A",
			OperationId: "N/A",
		})
	}

	for updatedSecurityName, updatedSecurity := range diffReport.ComponentsDiff.SecuritySchemesDiff.Modified {
		result = append(result, checkOAuthUpdates(updatedSecurity, config, updatedSecurityName)...)

		if updatedSecurity.TypeDiff != nil {
			result = append(result, BackwardCompatibilityError{
				Id:          APIComponentsSecurityTyepUpdated,
				Level:       INFO,
				Text:        fmt.Sprintf(config.i18n(APIComponentsSecurityTyepUpdated), ColorizedValue(updatedSecurityName), ColorizedValue(updatedSecurity.TypeDiff.From), ColorizedValue(updatedSecurity.TypeDiff.To)),
				Operation:   "N/A",
				Path:        "N/A",
				Source:      "N/A",
				OperationId: "N/A",
			})
		}
	}

	return result
}
