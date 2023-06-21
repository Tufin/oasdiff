package checker

import (
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/utils"
)

func GetDefaultChecks() BackwardCompatibilityCheckConfig {
	return GetChecks(utils.StringList{})
}

func GetChecks(includeChecks utils.StringList) BackwardCompatibilityCheckConfig {
	return getBackwardCompatibilityCheckConfig(allChecks(), LevelOverrides(includeChecks))
}

func LevelOverrides(includeChecks utils.StringList) map[string]Level {
	result := map[string]Level{}
	for _, s := range includeChecks {
		// if the checker was explicitly included with the `-include-checks`,
		// it means that it's output is considered a breaking change,
		// so the returned level should overwritten to ERR.
		result[s] = ERR
	}
	return result
}

func GetAllChecks(includeChecks utils.StringList) BackwardCompatibilityCheckConfig {
	return getBackwardCompatibilityCheckConfig(allChecks(), LevelOverrides(includeChecks))
}

func getBackwardCompatibilityCheckConfig(checks []BackwardCompatibilityCheck, levelOverrides map[string]Level) BackwardCompatibilityCheckConfig {
	return BackwardCompatibilityCheckConfig{
		Checks:              checks,
		LogLevelOverrides:   levelOverrides,
		MinSunsetBetaDays:   31,
		MinSunsetStableDays: 180,
		Localizer:           *localizations.New("en", "en"),
	}
}

var optionalChecks = map[string]BackwardCompatibilityCheck{
	"response-non-success-status-removed":   ResponseNonSuccessStatusUpdated,
	"api-operation-id-removed":              APIOperationIdUpdatedCheck,
	"api-tag-removed":                       APITagUpdatedCheck,
	"api-schema-removed":                    APIComponentsSchemaRemovedCheck,
	"response-property-enum-value-removed":  ResponseParameterEnumValueRemovedCheck,
	"response-mediatype-enum-value-removed": ResponseMediaTypeEnumValueRemovedCheck,
	"request-body-enum-value-removed":       RequestBodyEnumValueRemovedCheck,
}

func GetOptionalChecks() []string {
	result := make([]string, len(optionalChecks))
	i := 0
	for key := range optionalChecks {
		result[i] = key
		i++
	}

	return result
}

func defaultChecks() []BackwardCompatibilityCheck {
	return []BackwardCompatibilityCheck{
		RequestParameterRemovedCheck,
		NewRequiredRequestPropertyCheck,
		RequestParameterPatternAddedOrChangedCheck,
		RequestPropertyPatternAddedOrChangedCheck,
		AddedRequiredRequestBodyCheck,
		RequestParameterRequiredValueUpdatedCheck,
		RequestParameterBecameEnumCheck,
		RequestPropertyBecameRequiredCheck,
		RequestPropertyBecameEnumCheck,
		RequestHeaderPropertyBecameRequiredCheck,
		RequestHeaderPropertyBecameEnumCheck,
		ResponsePropertyBecameOptionalCheck,
		ResponsePropertyBecameNullableCheck,
		RequestPropertyBecameNotNullableCheck,
		RequestBodyRequiredUpdatedCheck,
		RequestBodyBecameEnumCheck,
		ResponseHeaderBecameOptional,
		ResponseHeaderRemoved,
		ResponseSuccessStatusUpdated,
		ResponseMediaTypeRemoved,
		NewRequestPathParameterCheck,
		NewRequestNonPathParameterCheck,
		NewRequiredRequestHeaderPropertyCheck,
		ResponseRequiredPropertyRemovedCheck,
		UncheckedRequestAllOfWarnCheck,
		UncheckedResponseAllOfWarnCheck,
		RequestPropertyRemovedCheck,
		ResponseRequiredPropertyBecameNonWriteOnlyCheck,
		RequestPropertyMaxLengthSetCheck,
		RequestParameterMaxLengthSetCheck,
		ResponsePropertyMaxLengthUnsetCheck,
		RequestParameterMaxLengthDecreasedCheck,
		RequestPropertyMaxLengthDecreasedCheck,
		ResponsePropertyMaxLengthIncreasedCheck,
		ResponsePropertyMinLengthDecreasedCheck,
		RequestPropertyMaxSetCheck,
		RequestPropertyMinSetCheck,
		RequestPropertyMaxDecreasedCheck,
		RequestPropertyMinIncreasedCheck,
		RequestParameterMaxSetCheck,
		RequestParameterMinSetCheck,
		RequestParameterMaxDecreasedCheck,
		RequestParameterMinIncreasedCheck,
		RequestParameterMinItemsSetCheck,
		RequestParameterMinItemsIncreasedCheck,
		RequestPropertyMinItemsSetCheck,
		RequestPropertyMinItemsIncreasedCheck,
		ResponsePropertyMinItemsUnsetCheck,
		ResponsePropertyMinItemsDecreasedCheck,
		RequestParameterEnumValueRemovedCheck,
		RequestPropertyEnumValueRemovedCheck,
		ResponsePropertyEnumValueAddedCheck,
		RequestParameterXExtensibleEnumValueRemovedCheck,
		RequestPropertyXExtensibleEnumValueRemovedCheck,
		RequestParameterTypeChangedCheck,
		RequestPropertyTypeChangedCheck,
		ResponsePropertyTypeChangedCheck,
		APIAddedCheck,
		APIRemovedCheck,
		APIDeprecationCheck,
		APISunsetChangedCheck,
		ResponsePropertyMaxIncreasedCheck,
		ResponsePropertyMinDecreasedCheck,
		RequestParameterDefaultValueChanged,
	}
}

func allChecks() []BackwardCompatibilityCheck {
	checks := defaultChecks()
	for _, v := range optionalChecks {
		checks = append(checks, v)
	}
	return checks
}
