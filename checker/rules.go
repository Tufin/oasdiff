package checker

import (
	"fmt"

	"github.com/tufin/oasdiff/utils"
)

type BackwardCompatibilityRule struct {
	Id          string
	Level       Level
	Description string
	Handler     BackwardCompatibilityCheck `json:"-" yaml:"-"`
}

func newBackwardCompatibilityRule(id string, level Level, handler BackwardCompatibilityCheck) BackwardCompatibilityRule {
	return BackwardCompatibilityRule{
		Id:          id,
		Level:       level,
		Description: descriptionId(id),
		Handler:     handler,
	}
}

type BackwardCompatibilityRules []BackwardCompatibilityRule

func GetAllRules() BackwardCompatibilityRules {
	return BackwardCompatibilityRules{
		// APIAddedCheck
		newBackwardCompatibilityRule(EndpointAddedId, INFO, APIAddedCheck),
		// APIComponentsSecurityUpdatedCheck
		newBackwardCompatibilityRule(APIComponentsSecurityRemovedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentsSecurityAddedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentsSecurityComponentOauthUrlUpdatedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentsSecurityTypeUpdatedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentsSecurityOauthTokenUrlUpdatedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeAddedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeRemovedId, INFO, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeUpdatedId, INFO, APIComponentsSecurityUpdatedCheck),
		// APISecurityUpdatedCheck
		newBackwardCompatibilityRule(APISecurityRemovedCheckId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APISecurityAddedCheckId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APISecurityScopeAddedId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APISecurityScopeRemovedId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIGlobalSecurityRemovedCheckId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIGlobalSecurityAddedCheckId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIGlobalSecurityScopeAddedId, INFO, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(APIGlobalSecurityScopeRemovedId, INFO, APISecurityUpdatedCheck),
		// Stability Descreased Check is run as part of CheckBackwardCompatibility
		newBackwardCompatibilityRule(APIStabilityDecreasedId, ERR, nil),
		// APIDeprecationCheck
		newBackwardCompatibilityRule(EndpointReactivatedId, INFO, APIDeprecationCheck),
		newBackwardCompatibilityRule(APIDeprecatedSunsetParseId, ERR, APIDeprecationCheck),
		newBackwardCompatibilityRule(APIDeprecatedSunsetMissingId, ERR, APIDeprecationCheck),
		newBackwardCompatibilityRule(APIInvalidStabilityLevelId, ERR, APIDeprecationCheck),
		newBackwardCompatibilityRule(APISunsetDateTooSmallId, ERR, APIDeprecationCheck),
		newBackwardCompatibilityRule(EndpointDeprecatedId, INFO, APIDeprecationCheck),
		// APIRemovedCheck
		newBackwardCompatibilityRule(APIPathRemovedWithoutDeprecationId, ERR, APIRemovedCheck),
		newBackwardCompatibilityRule(APIPathSunsetParseId, ERR, APIRemovedCheck),
		newBackwardCompatibilityRule(APIPathRemovedBeforeSunsetId, ERR, APIRemovedCheck),
		newBackwardCompatibilityRule(APIRemovedWithoutDeprecationId, ERR, APIRemovedCheck),
		newBackwardCompatibilityRule(APIRemovedBeforeSunsetId, ERR, APIRemovedCheck),
		// APISunsetChangedCheck
		newBackwardCompatibilityRule(APISunsetDeletedId, ERR, APISunsetChangedCheck),
		newBackwardCompatibilityRule(APISunsetDateChangedTooSmallId, ERR, APISunsetChangedCheck),
		// AddedRequiredRequestBodyCheck
		newBackwardCompatibilityRule(AddedRequiredRequestBodyId, ERR, AddedRequestBodyCheck),
		newBackwardCompatibilityRule(AddedOptionalRequestBodyId, INFO, AddedRequestBodyCheck),
		// NewRequestNonPathDefaultParameterCheck
		newBackwardCompatibilityRule(NewRequiredRequestDefaultParameterToExistingPathId, ERR, NewRequestNonPathDefaultParameterCheck),
		newBackwardCompatibilityRule(NewOptionalRequestDefaultParameterToExistingPathId, INFO, NewRequestNonPathDefaultParameterCheck),
		// NewRequestNonPathParameterCheck
		newBackwardCompatibilityRule(NewRequiredRequestParameterId, ERR, NewRequestNonPathParameterCheck),
		newBackwardCompatibilityRule(NewOptionalRequestParameterId, INFO, NewRequestNonPathParameterCheck),
		// NewRequestPathParameterCheck
		newBackwardCompatibilityRule(NewRequestPathParameterId, ERR, NewRequestPathParameterCheck),
		// NewRequiredRequestHeaderPropertyCheck
		newBackwardCompatibilityRule(NewRequiredRequestHeaderPropertyId, ERR, NewRequiredRequestHeaderPropertyCheck),
		// RequestBodyBecameEnumCheck
		newBackwardCompatibilityRule(RequestBodyBecameEnumId, ERR, RequestBodyBecameEnumCheck),
		// RequestBodyMediaTypeChangedCheck
		newBackwardCompatibilityRule(RequestBodyMediaTypeAddedId, INFO, RequestBodyMediaTypeChangedCheck),
		newBackwardCompatibilityRule(RequestBodyMediaTypeRemovedId, ERR, RequestBodyMediaTypeChangedCheck),
		// RequestBodyRequiredUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyBecameOptionalId, INFO, RequestBodyRequiredUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyBecameRequiredId, ERR, RequestBodyRequiredUpdatedCheck),
		// RequestDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyDiscriminatorAddedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorRemovedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorPropertyNameChangedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingAddedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingDeletedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingChangedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorAddedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorRemovedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorPropertyNameChangedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingAddedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingDeletedId, INFO, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingChangedId, INFO, RequestDiscriminatorUpdatedCheck),
		// RequestHeaderPropertyBecameEnumCheck
		newBackwardCompatibilityRule(RequestHeaderPropertyBecameEnumId, ERR, RequestHeaderPropertyBecameEnumCheck),
		// RequestHeaderPropertyBecameRequiredCheck
		newBackwardCompatibilityRule(RequestHeaderPropertyBecameRequiredId, ERR, RequestHeaderPropertyBecameRequiredCheck),
		// RequestParameterBecameEnumCheck
		newBackwardCompatibilityRule(RequestParameterBecameEnumId, ERR, RequestParameterBecameEnumCheck),
		// RequestParameterDefaultValueChangedCheck
		newBackwardCompatibilityRule(RequestParameterDefaultValueChangedId, ERR, RequestParameterDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestParameterDefaultValueAddedId, ERR, RequestParameterDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestParameterDefaultValueRemovedId, ERR, RequestParameterDefaultValueChangedCheck),
		// RequestParameterEnumValueUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterEnumValueAddedId, INFO, RequestParameterEnumValueUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterEnumValueRemovedId, ERR, RequestParameterEnumValueUpdatedCheck),
		// RequestParameterMaxItemsUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxItemsIncreasedId, INFO, RequestParameterMaxItemsUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterMaxItemsDecreasedId, ERR, RequestParameterMaxItemsUpdatedCheck),
		// RequestParameterMaxLengthSetCheck
		newBackwardCompatibilityRule(RequestParameterMaxLengthSetId, WARN, RequestParameterMaxLengthSetCheck),
		// RequestParameterMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxLengthIncreasedId, INFO, RequestParameterMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterMaxLengthDecreasedId, ERR, RequestParameterMaxLengthUpdatedCheck),
		// RequestParameterMaxSetCheck
		newBackwardCompatibilityRule(RequestParameterMaxSetId, WARN, RequestParameterMaxSetCheck),
		// RequestParameterMaxUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxIncreasedId, INFO, RequestParameterMaxUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterMaxDecreasedId, ERR, RequestParameterMaxUpdatedCheck),
		// RequestParameterMinItemsSetCheck
		newBackwardCompatibilityRule(RequestParameterMinItemsSetId, WARN, RequestParameterMinItemsSetCheck),
		// RequestParameterMinItemsUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinItemsIncreasedId, ERR, RequestParameterMinItemsUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterMinItemsDecreasedId, INFO, RequestParameterMinItemsUpdatedCheck),
		// RequestParameterMinLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinLengthIncreasedId, ERR, RequestParameterMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterMinLengthDecreasedId, INFO, RequestParameterMinLengthUpdatedCheck),
		// RequestParameterMinSetCheck
		newBackwardCompatibilityRule(RequestParameterMinSetId, WARN, RequestParameterMinSetCheck),
		// RequestParameterMinUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinIncreasedId, ERR, RequestParameterMinUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterMinDecreasedId, INFO, RequestParameterMinUpdatedCheck),
		// RequestParameterPatternAddedOrChangedCheck
		newBackwardCompatibilityRule(RequestParameterPatternAddedId, WARN, RequestParameterPatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(RequestParameterPatternRemovedId, INFO, RequestParameterPatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(RequestParameterPatternChangedId, WARN, RequestParameterPatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(RequestParameterPatternGeneralizedId, INFO, RequestParameterPatternAddedOrChangedCheck),
		// RequestParameterRemovedCheck
		newBackwardCompatibilityRule(RequestParameterRemovedId, WARN, RequestParameterRemovedCheck),
		// RequestParameterRequiredValueUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterBecomeRequiredId, ERR, RequestParameterRequiredValueUpdatedCheck),
		newBackwardCompatibilityRule(RequestParameterBecomeOptionalId, INFO, RequestParameterRequiredValueUpdatedCheck),
		// RequestParameterTypeChangedCheck
		newBackwardCompatibilityRule(RequestParameterTypeChangedId, ERR, RequestParameterTypeChangedCheck),
		newBackwardCompatibilityRule(RequestParameterTypeGeneralizedId, INFO, RequestParameterTypeChangedCheck),
		newBackwardCompatibilityRule(RequestParameterPropertyTypeChangedId, WARN, RequestParameterTypeChangedCheck),
		newBackwardCompatibilityRule(RequestParameterPropertyTypeGeneralizedId, INFO, RequestParameterTypeChangedCheck),
		newBackwardCompatibilityRule(RequestParameterPropertyTypeSpecializedId, ERR, RequestParameterTypeChangedCheck),
		// RequestParameterXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestParameterXExtensibleEnumValueRemovedId, ERR, RequestParameterXExtensibleEnumValueRemovedCheck),
		// RequestPropertyAllOfUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyAllOfAddedId, ERR, RequestPropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyAllOfRemovedId, WARN, RequestPropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyAllOfAddedId, ERR, RequestPropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyAllOfRemovedId, WARN, RequestPropertyAllOfUpdatedCheck),
		// RequestPropertyAnyOfUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyAnyOfAddedId, INFO, RequestPropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyAnyOfRemovedId, ERR, RequestPropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyAnyOfAddedId, INFO, RequestPropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyAnyOfRemovedId, ERR, RequestPropertyAnyOfUpdatedCheck),
		// RequestPropertyBecameEnumCheck
		newBackwardCompatibilityRule(RequestPropertyBecameEnumId, ERR, RequestPropertyBecameEnumCheck),
		// RequestPropertyBecameNotNullableCheck
		newBackwardCompatibilityRule(RequestBodyBecomeNotNullableId, ERR, RequestPropertyBecameNotNullableCheck),
		newBackwardCompatibilityRule(RequestBodyBecomeNullableId, INFO, RequestPropertyBecameNotNullableCheck),
		newBackwardCompatibilityRule(RequestPropertyBecomeNotNullableId, ERR, RequestPropertyBecameNotNullableCheck),
		newBackwardCompatibilityRule(RequestPropertyBecomeNullableId, INFO, RequestPropertyBecameNotNullableCheck),
		// RequestPropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(RequestBodyDefaultValueAddedId, INFO, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestBodyDefaultValueRemovedId, INFO, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestBodyDefaultValueChangedId, INFO, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestPropertyDefaultValueAddedId, INFO, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestPropertyDefaultValueRemovedId, INFO, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(RequestPropertyDefaultValueChangedId, INFO, RequestPropertyDefaultValueChangedCheck),
		// RequestPropertyEnumValueUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyEnumValueRemovedId, ERR, RequestPropertyEnumValueUpdatedCheck),
		newBackwardCompatibilityRule(RequestReadOnlyPropertyEnumValueRemovedId, INFO, RequestPropertyEnumValueUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyEnumValueAddedId, INFO, RequestPropertyEnumValueUpdatedCheck),
		// RequestPropertyMaxDecreasedCheck
		newBackwardCompatibilityRule(RequestBodyMaxDecreasedId, ERR, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(RequestBodyMaxIncreasedId, INFO, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(RequestPropertyMaxDecreasedId, ERR, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(RequestReadOnlyPropertyMaxDecreasedId, INFO, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(RequestPropertyMaxIncreasedId, INFO, RequestPropertyMaxDecreasedCheck),
		// RequestPropertyMaxLengthSetCheck
		newBackwardCompatibilityRule(RequestBodyMaxLengthSetId, WARN, RequestPropertyMaxLengthSetCheck),
		newBackwardCompatibilityRule(RequestPropertyMaxLengthSetId, WARN, RequestPropertyMaxLengthSetCheck),
		// RequestPropertyMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyMaxLengthDecreasedId, ERR, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyMaxLengthIncreasedId, INFO, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyMaxLengthDecreasedId, ERR, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestReadOnlyPropertyMaxLengthDecreasedId, INFO, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyMaxLengthIncreasedId, INFO, RequestPropertyMaxLengthUpdatedCheck),
		// RequestPropertyMaxSetCheck
		newBackwardCompatibilityRule(RequestBodyMaxSetId, WARN, RequestPropertyMaxSetCheck),
		newBackwardCompatibilityRule(RequestPropertyMaxSetId, WARN, RequestPropertyMaxSetCheck),
		// RequestPropertyMinIncreasedCheck
		newBackwardCompatibilityRule(RequestBodyMinIncreasedId, ERR, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(RequestBodyMinDecreasedId, INFO, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(RequestPropertyMinIncreasedId, ERR, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(RequestReadOnlyPropertyMinIncreasedId, INFO, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(RequestPropertyMinDecreasedId, INFO, RequestPropertyMinIncreasedCheck),
		// RequestPropertyMinItemsIncreasedCheck
		newBackwardCompatibilityRule(RequestBodyMinItemsIncreasedId, ERR, RequestPropertyMinItemsIncreasedCheck),
		newBackwardCompatibilityRule(RequestPropertyMinItemsIncreasedId, ERR, RequestPropertyMinItemsIncreasedCheck),
		// RequestPropertyMinItemsSetCheck
		newBackwardCompatibilityRule(RequestBodyMinItemsSetId, WARN, RequestPropertyMinItemsSetCheck),
		newBackwardCompatibilityRule(RequestPropertyMinItemsSetId, WARN, RequestPropertyMinItemsSetCheck),
		// RequestPropertyMinLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyMinLengthIncreasedId, ERR, RequestPropertyMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyMinLengthDecreasedId, INFO, RequestPropertyMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyMinLengthIncreasedId, ERR, RequestPropertyMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyMinLengthDecreasedId, INFO, RequestPropertyMinLengthUpdatedCheck),
		// RequestPropertyMinSetCheck
		newBackwardCompatibilityRule(RequestBodyMinSetId, WARN, RequestPropertyMinSetCheck),
		newBackwardCompatibilityRule(RequestPropertyMinSetId, WARN, RequestPropertyMinSetCheck),
		// RequestPropertyOneOfUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyOneOfAddedId, INFO, RequestPropertyOneOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestBodyOneOfRemovedId, ERR, RequestPropertyOneOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyOneOfAddedId, INFO, RequestPropertyOneOfUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyOneOfRemovedId, ERR, RequestPropertyOneOfUpdatedCheck),
		// RequestPropertyPatternUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyPatternRemovedId, INFO, RequestPropertyPatternUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyPatternAddedId, WARN, RequestPropertyPatternUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyPatternChangedId, WARN, RequestPropertyPatternUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyPatternGeneralizedId, INFO, RequestPropertyPatternUpdatedCheck),
		// RequestPropertyRequiredUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyBecameRequiredId, ERR, RequestPropertyRequiredUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyBecameRequiredWithDefaultId, INFO, RequestPropertyRequiredUpdatedCheck),
		newBackwardCompatibilityRule(RequestPropertyBecameOptionalId, INFO, RequestPropertyRequiredUpdatedCheck),
		// RequestPropertyTypeChangedCheck
		newBackwardCompatibilityRule(RequestBodyTypeGeneralizedId, INFO, RequestPropertyTypeChangedCheck),
		newBackwardCompatibilityRule(RequestBodyTypeChangedId, ERR, RequestPropertyTypeChangedCheck),
		newBackwardCompatibilityRule(RequestPropertyTypeGeneralizedId, INFO, RequestPropertyTypeChangedCheck),
		newBackwardCompatibilityRule(RequestPropertyTypeChangedId, ERR, RequestPropertyTypeChangedCheck),
		// RequestPropertyUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyRemovedId, WARN, RequestPropertyUpdatedCheck),
		newBackwardCompatibilityRule(NewRequiredRequestPropertyId, ERR, RequestPropertyUpdatedCheck),
		newBackwardCompatibilityRule(NewRequiredRequestPropertyWithDefaultId, INFO, RequestPropertyUpdatedCheck),
		newBackwardCompatibilityRule(NewOptionalRequestPropertyId, INFO, RequestPropertyUpdatedCheck),
		// RequestPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameNonWriteOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameWriteOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameReadOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameNonReadOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameNonWriteOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameWriteOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameReadOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameNonReadOnlyCheckId, INFO, RequestPropertyWriteOnlyReadOnlyCheck),
		// RequestPropertyXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestPropertyXExtensibleEnumValueRemovedId, ERR, RequestPropertyXExtensibleEnumValueRemovedCheck),
		// ResponseDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorAddedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorRemovedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorPropertyNameChangedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingAddedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingDeletedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingChangedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorAddedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorRemovedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorPropertyNameChangedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingAddedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingDeletedId, INFO, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingChangedId, INFO, ResponseDiscriminatorUpdatedCheck),
		// ResponseHeaderBecameOptionalCheck
		newBackwardCompatibilityRule(ResponseHeaderBecameOptionalId, ERR, ResponseHeaderBecameOptionalCheck),
		// ResponseHeaderRemovedCheck
		newBackwardCompatibilityRule(RequiredResponseHeaderRemovedId, ERR, ResponseHeaderRemovedCheck),
		newBackwardCompatibilityRule(OptionalResponseHeaderRemovedId, WARN, ResponseHeaderRemovedCheck),
		// ResponseMediaTypeUpdatedCheck
		newBackwardCompatibilityRule(ResponseMediaTypeRemovedId, ERR, ResponseMediaTypeUpdatedCheck),
		newBackwardCompatibilityRule(ResponseMediaTypeAddedId, INFO, ResponseMediaTypeUpdatedCheck),
		// ResponseOptionalPropertyUpdatedCheck
		newBackwardCompatibilityRule(ResponseOptionalPropertyRemovedId, WARN, ResponseOptionalPropertyUpdatedCheck),
		newBackwardCompatibilityRule(ResponseOptionalWriteOnlyPropertyRemovedId, INFO, ResponseOptionalPropertyUpdatedCheck),
		newBackwardCompatibilityRule(ResponseOptionalPropertyAddedId, INFO, ResponseOptionalPropertyUpdatedCheck),
		newBackwardCompatibilityRule(ResponseOptionalWriteOnlyPropertyAddedId, INFO, ResponseOptionalPropertyUpdatedCheck),
		// ResponseOptionalPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameNonWriteOnlyId, INFO, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameWriteOnlyId, INFO, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameReadOnlyId, INFO, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameNonReadOnlyId, INFO, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		// ResponsePatternAddedOrChangedCheck
		newBackwardCompatibilityRule(ResponsePropertyPatternAddedId, INFO, ResponsePatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(ResponsePropertyPatternChangedId, INFO, ResponsePatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(ResponsePropertyPatternRemovedId, INFO, ResponsePatternAddedOrChangedCheck),
		// ResponsePropertyAllOfUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyAllOfAddedId, INFO, ResponsePropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyAllOfRemovedId, INFO, ResponsePropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyAllOfAddedId, INFO, ResponsePropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyAllOfRemovedId, INFO, ResponsePropertyAllOfUpdatedCheck),
		// ResponsePropertyAnyOfUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyAnyOfAddedId, INFO, ResponsePropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(ResponseBodyAnyOfRemovedId, INFO, ResponsePropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyAnyOfAddedId, INFO, ResponsePropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(ResponsePropertyAnyOfRemovedId, INFO, ResponsePropertyAnyOfUpdatedCheck),
		// ResponsePropertyBecameNullableCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameNullableId, ERR, ResponsePropertyBecameNullableCheck),
		newBackwardCompatibilityRule(ResponseBodyBecameNullableId, ERR, ResponsePropertyBecameNullableCheck),
		// ResponsePropertyBecameOptionalCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameOptionalId, ERR, ResponsePropertyBecameOptionalCheck),
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyBecameOptionalId, INFO, ResponsePropertyBecameOptionalCheck),
		// ResponsePropertyBecameRequiredCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameRequiredId, INFO, ResponsePropertyBecameRequiredCheck),
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyBecameRequiredId, INFO, ResponsePropertyBecameRequiredCheck),
		// ResponsePropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(ResponseBodyDefaultValueAddedId, INFO, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(ResponseBodyDefaultValueRemovedId, INFO, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(ResponseBodyDefaultValueChangedId, INFO, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueAddedId, INFO, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueRemovedId, INFO, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueChangedId, INFO, ResponsePropertyDefaultValueChangedCheck),
		// ResponsePropertyEnumValueAddedCheck
		newBackwardCompatibilityRule(ResponsePropertyEnumValueAddedId, WARN, ResponsePropertyEnumValueAddedCheck),
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyEnumValueAddedId, INFO, ResponsePropertyEnumValueAddedCheck),
		// ResponsePropertyMaxIncreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMaxIncreasedId, ERR, ResponsePropertyMaxIncreasedCheck),
		newBackwardCompatibilityRule(ResponsePropertyMaxIncreasedId, ERR, ResponsePropertyMaxIncreasedCheck),
		// ResponsePropertyMaxLengthIncreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMaxLengthIncreasedId, ERR, ResponsePropertyMaxLengthIncreasedCheck),
		newBackwardCompatibilityRule(ResponsePropertyMaxLengthIncreasedId, ERR, ResponsePropertyMaxLengthIncreasedCheck),
		// ResponsePropertyMaxLengthUnsetCheck
		newBackwardCompatibilityRule(ResponseBodyMaxLengthUnsetId, ERR, ResponsePropertyMaxLengthUnsetCheck),
		newBackwardCompatibilityRule(ResponsePropertyMaxLengthUnsetId, ERR, ResponsePropertyMaxLengthUnsetCheck),
		// ResponsePropertyMinDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinDecreasedId, ERR, ResponsePropertyMinDecreasedCheck),
		newBackwardCompatibilityRule(ResponsePropertyMinDecreasedId, ERR, ResponsePropertyMinDecreasedCheck),
		// ResponsePropertyMinItemsDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinItemsDecreasedId, ERR, ResponsePropertyMinItemsDecreasedCheck),
		newBackwardCompatibilityRule(ResponsePropertyMinItemsDecreasedId, ERR, ResponsePropertyMinItemsDecreasedCheck),
		// ResponsePropertyMinItemsUnsetCheck
		newBackwardCompatibilityRule(ResponseBodyMinItemsUnsetId, ERR, ResponsePropertyMinItemsUnsetCheck),
		newBackwardCompatibilityRule(ResponsePropertyMinItemsUnsetId, ERR, ResponsePropertyMinItemsUnsetCheck),
		// ResponsePropertyMinLengthDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinLengthDecreasedId, ERR, ResponsePropertyMinLengthDecreasedCheck),
		newBackwardCompatibilityRule(ResponsePropertyMinLengthDecreasedId, ERR, ResponsePropertyMinLengthDecreasedCheck),
		// ResponsePropertyOneOfUpdated
		newBackwardCompatibilityRule(ResponseBodyOneOfAddedId, INFO, ResponsePropertyOneOfUpdated),
		newBackwardCompatibilityRule(ResponseBodyOneOfRemovedId, INFO, ResponsePropertyOneOfUpdated),
		newBackwardCompatibilityRule(ResponsePropertyOneOfAddedId, INFO, ResponsePropertyOneOfUpdated),
		newBackwardCompatibilityRule(ResponsePropertyOneOfRemovedId, INFO, ResponsePropertyOneOfUpdated),
		// ResponsePropertyTypeChangedCheck
		newBackwardCompatibilityRule(ResponseBodyTypeChangedId, ERR, ResponsePropertyTypeChangedCheck),
		newBackwardCompatibilityRule(ResponsePropertyTypeChangedId, ERR, ResponsePropertyTypeChangedCheck),
		// ResponseRequiredPropertyUpdatedCheck
		newBackwardCompatibilityRule(ResponseRequiredPropertyRemovedId, ERR, ResponseRequiredPropertyUpdatedCheck),
		newBackwardCompatibilityRule(ResponseRequiredWriteOnlyPropertyRemovedId, INFO, ResponseRequiredPropertyUpdatedCheck),
		newBackwardCompatibilityRule(ResponseRequiredPropertyAddedId, INFO, ResponseRequiredPropertyUpdatedCheck),
		newBackwardCompatibilityRule(ResponseRequiredWriteOnlyPropertyAddedId, INFO, ResponseRequiredPropertyUpdatedCheck),
		// ResponseRequiredPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameNonWriteOnlyId, WARN, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameWriteOnlyId, INFO, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameReadOnlyId, INFO, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameNonReadOnlyId, INFO, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		// ResponseSuccessStatusUpdatedCheck
		newBackwardCompatibilityRule(ResponseSuccessStatusRemovedId, ERR, ResponseSuccessStatusUpdatedCheck),
		newBackwardCompatibilityRule(ResponseSuccessStatusAddedId, INFO, ResponseSuccessStatusUpdatedCheck),
		// ResponseNonSuccessStatusUpdatedCheck
		newBackwardCompatibilityRule(ResponseNonSuccessStatusRemovedId, INFO, ResponseNonSuccessStatusUpdatedCheck), // optional
		newBackwardCompatibilityRule(ResponseNonSuccessStatusAddedId, INFO, ResponseNonSuccessStatusUpdatedCheck),
		// APIOperationIdUpdatedCheck
		newBackwardCompatibilityRule(APIOperationIdRemovedId, INFO, APIOperationIdUpdatedCheck), // optional
		newBackwardCompatibilityRule(APIOperationIdAddId, INFO, APIOperationIdUpdatedCheck),
		// APITagUpdatedCheck
		newBackwardCompatibilityRule(APITagRemovedId, INFO, APITagUpdatedCheck), // optional
		newBackwardCompatibilityRule(APITagAddedId, INFO, APITagUpdatedCheck),
		// APIComponentsSchemaRemovedCheck
		newBackwardCompatibilityRule(APISchemasRemovedId, INFO, APIComponentsSchemaRemovedCheck), // optional
		// ResponseParameterEnumValueRemovedCheck
		newBackwardCompatibilityRule(ResponsePropertyEnumValueRemovedId, INFO, ResponseParameterEnumValueRemovedCheck), // optional
		// ResponseMediaTypeEnumValueRemovedCheck
		newBackwardCompatibilityRule(ResponseMediaTypeEnumValueRemovedId, INFO, ResponseMediaTypeEnumValueRemovedCheck), // optional
		// RequestBodyEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestBodyEnumValueRemovedId, INFO, RequestBodyEnumValueRemovedCheck), // optional
	}
}

func GetOptionalRules() BackwardCompatibilityRules {
	return BackwardCompatibilityRules{
		newBackwardCompatibilityRule(ResponseNonSuccessStatusRemovedId, INFO, ResponseNonSuccessStatusUpdatedCheck),
		newBackwardCompatibilityRule(APIOperationIdRemovedId, INFO, APIOperationIdUpdatedCheck),
		newBackwardCompatibilityRule(APITagRemovedId, INFO, APITagUpdatedCheck),
		newBackwardCompatibilityRule(APISchemasRemovedId, INFO, APIComponentsSchemaRemovedCheck),
		newBackwardCompatibilityRule(ResponsePropertyEnumValueRemovedId, INFO, ResponseParameterEnumValueRemovedCheck),
		newBackwardCompatibilityRule(ResponseMediaTypeEnumValueRemovedId, INFO, ResponseMediaTypeEnumValueRemovedCheck),
		newBackwardCompatibilityRule(RequestBodyEnumValueRemovedId, INFO, RequestBodyEnumValueRemovedCheck),
	}
}

// GetCheckLevels gets levels for all backward compatibility checks
func GetCheckLevels() map[string]Level {
	return rulesToLevels(GetAllRules())
}

// GetAllChecks gets all backward compatibility checks
func GetAllChecks() BackwardCompatibilityChecks {
	return rulesToChecks(GetAllRules())
}

// rulesToChecks return a unique list of checks from a list of rules
func rulesToChecks(rules BackwardCompatibilityRules) BackwardCompatibilityChecks {
	result := BackwardCompatibilityChecks{}
	m := utils.StringSet{}
	for _, rule := range rules {
		// functions are not comparable, so we convert them to strings
		pStr := fmt.Sprintf("%v", rule.Handler)
		if !m.Contains(pStr) {
			m.Add(pStr)
			result = append(result, rule.Handler)
		}
	}
	return result
}

func GetOptionalRuleIds() []string {
	return rulesToIIs(GetOptionalRules())
}

func GetAllRuleIds() []string {
	return rulesToIIs(GetAllRules())
}

// rulesToLevels return a map of check IDs to levels
func rulesToLevels(rules BackwardCompatibilityRules) map[string]Level {
	result := map[string]Level{}
	for _, rule := range rules {
		result[rule.Id] = rule.Level
	}
	return result
}

func rulesToIIs(rules BackwardCompatibilityRules) []string {
	result := []string{}
	for _, rule := range rules {
		result = append(result, rule.Id)
	}
	return result
}
