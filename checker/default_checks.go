package checker

import (
	"github.com/tufin/oasdiff/utils"
)

func GetDefaultChecks() Config {
	return GetChecks(utils.StringList{})
}

func GetChecks(includeChecks utils.StringList) Config {
	return getBackwardCompatibilityCheckConfig(allChecks(), LevelOverrides(includeChecks), BetaDeprecationDays, StableDeprecationDays)
}

func LevelOverrides(includeChecks utils.StringList) map[string]Level {
	result := map[string]Level{}
	for _, s := range includeChecks {
		// if the checker was explicitly included with the `--include-checks`,
		// it means that it's output is considered a breaking change,
		// so the returned level should overwritten to ERR.
		result[s] = ERR
	}
	return result
}

func GetAllChecks(includeChecks utils.StringList, deprecationDaysBeta int, deprecationDaysStable int) Config {
	return getBackwardCompatibilityCheckConfig(allChecks(), LevelOverrides(includeChecks), deprecationDaysBeta, deprecationDaysStable)
}

func getBackwardCompatibilityCheckConfig(checks []BackwardCompatibilityCheck, levelOverrides map[string]Level, minSunsetBetaDays int, minSunsetStableDays int) Config {
	return Config{
		Checks:              checks,
		LogLevelOverrides:   levelOverrides,
		MinSunsetBetaDays:   minSunsetBetaDays,
		MinSunsetStableDays: minSunsetStableDays,
		Localize:            NewLocalizer("en", "en"),
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
		RequestParameterPatternAddedOrChangedCheck,
		RequestPropertyPatternUpdatedCheck,
		AddedRequiredRequestBodyCheck,
		RequestParameterRequiredValueUpdatedCheck,
		RequestParameterBecameEnumCheck,
		RequestPropertyRequiredUpdatedCheck,
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
		ResponseMediaTypeUpdated,
		NewRequestPathParameterCheck,
		NewRequestNonPathParameterCheck,
		NewRequiredRequestHeaderPropertyCheck,
		ResponseRequiredPropertyUpdatedCheck,
		UncheckedRequestAllOfWarnCheck,
		UncheckedResponseAllOfWarnCheck,
		RequestPropertyUpdatedCheck,
		ResponseRequiredPropertyWriteOnlyReadOnlyCheck,
		RequestPropertyMaxLengthSetCheck,
		RequestParameterMaxLengthSetCheck,
		ResponsePropertyMaxLengthUnsetCheck,
		RequestPropertyMaxLengthUpdatedCheck,
		RequestParameterMaxLengthUpdatedCheck,
		ResponsePropertyMaxLengthIncreasedCheck,
		ResponsePropertyMinLengthDecreasedCheck,
		RequestPropertyMaxSetCheck,
		RequestPropertyMinSetCheck,
		RequestPropertyMaxDecreasedCheck,
		RequestPropertyMinIncreasedCheck,
		RequestParameterMaxSetCheck,
		RequestParameterMinSetCheck,
		RequestParameterMaxUpdatedCheck,
		RequestParameterMinUpdatedCheck,
		RequestParameterMinItemsSetCheck,
		RequestParameterMinItemsUpdatedCheck,
		RequestPropertyMinItemsSetCheck,
		RequestPropertyMinItemsIncreasedCheck,
		ResponsePropertyMinItemsUnsetCheck,
		ResponsePropertyMinItemsDecreasedCheck,
		RequestPropertyEnumValueUpdatedCheck,
		RequestParameterEnumValueUpdatedCheck,
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
		APIComponentsSecurityUpdatedCheck,
		APISecurityUpdatedCheck,
		RequestBodyMediaTypeChangedCheck,
		ResponseOptionalPropertyWriteOnlyReadOnlyCheck,
		ResponsePatternAddedOrChangedCheck,
		ResponsePropertyDefaultValueChangedCheck,
		RequestPropertyAnyOfUpdated,
		ResponsePropertyAnyOfUpdated,
		RequestPropertyOneOfUpdated,
		ResponsePropertyOneOfUpdated,
		RequestPropertyAllOfUpdated,
		ResponsePropertyAllOfUpdated,
		RequestPropertyWriteOnlyReadOnlyCheck,
		RequestDiscriminatorUpdatedCheck,
		ResponseDiscriminatorUpdatedCheck,
		RequestPropertyDefaultValueChangedCheck,
	}
}

func allChecks() []BackwardCompatibilityCheck {
	checks := defaultChecks()
	for _, v := range optionalChecks {
		checks = append(checks, v)
	}
	return checks
}
