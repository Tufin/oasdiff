package checker

import (
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/utils"
)

func GetDefaultChecks() BackwardCompatibilityCheckConfig {
	return GetChecks(utils.StringList{})
}

func GetChecks(includeChecks utils.StringList) BackwardCompatibilityCheckConfig {
	return getBackwardCompatibilityCheckConfig(append(defaultChecks(), includedChecks(includeChecks)...))
}

func GetAllChecks() BackwardCompatibilityCheckConfig {
	return getBackwardCompatibilityCheckConfig(allChecks())
}

func getBackwardCompatibilityCheckConfig(checks []BackwardCompatibilityCheck) BackwardCompatibilityCheckConfig {
	return BackwardCompatibilityCheckConfig{
		Checks:              checks,
		MinSunsetBetaDays:   31,
		MinSunsetStableDays: 180,
		Localizer:           *localizations.New("en", "en"),
	}
}

var optionalChecks = map[string]BackwardCompatibilityCheck{
	"response-non-success-status-removed":   ResponseNonSuccessStatusRemoved,
	"api-operation-id-removed":              APIOperationIdRemovedCheck,
	"api-tag-removed":                       APITagRemovedCheck,
	"response-property-enum-value-removed":  ResponseParameterEnumValueRemovedCheck,
	"response-mediatype-enum-value-removed": ResponseMediaTypeEnumValueRemovedCheck,
	"request-body-enum-value-removed":       RequestBodyEnumValueRemovedCheck,
	"api-path-added":                        APIAddedCheck,
}

func ValidateIncludeChecks(includeChecks utils.StringList) utils.StringList {
	result := utils.StringList{}
	for _, s := range includeChecks {
		if _, ok := optionalChecks[s]; !ok {
			result = append(result, s)
		}
	}

	return result.Sort()
}

func includedChecks(includeChecks utils.StringList) []BackwardCompatibilityCheck {
	result := []BackwardCompatibilityCheck{}
	for _, s := range includeChecks {
		result = append(result, optionalChecks[s])
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
		RequestParameterBecameRequiredCheck,
		RequestParameterBecameEnumCheck,
		RequestPropertyBecameRequiredCheck,
		RequestPropertyBecameEnumCheck,
		RequestHeaderPropertyBecameRequiredCheck,
		RequestHeaderPropertyBecameEnumCheck,
		ResponsePropertyBecameOptionalCheck,
		ResponsePropertyBecameNullableCheck,
		RequestPropertyBecameNotNullableCheck,
		RequestBodyBecameRequiredCheck,
		RequestBodyBecameEnumCheck,
		ResponseHeaderBecameOptional,
		ResponseHeaderRemoved,
		ResponseSuccessStatusRemoved,
		ResponseMediaTypeRemoved,
		NewRequestPathParameterCheck,
		NewRequiredRequestNonPathParameterCheck,
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
