package checker

func getDescription(l Localizer, id string) string {
	return l(id + "-description")
}

func newBackwardCompatibilityRule(l Localizer, id string, level Level, required bool, handler BackwardCompatibilityCheck) BackwardCompatibilityRule {
	return BackwardCompatibilityRule{
		Id:          id,
		Level:       level,
		Description: getDescription(l, id),
		Required:    required,
		Handler:     handler,
	}
}

func GetAllRules(l Localizer) []BackwardCompatibilityRule {
	return []BackwardCompatibilityRule{
		// APIAddedCheck
		newBackwardCompatibilityRule(l, EndpointAddedId, INFO, true, APIAddedCheck),
		// APIComponentsSecurityUpdatedCheck
		newBackwardCompatibilityRule(l, APIComponentsSecurityRemovedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentsSecurityAddedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentsSecurityComponentOauthUrlUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentsSecurityTypeUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentsSecurityOauthTokenUrlUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentSecurityOauthScopeAddedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentSecurityOauthScopeRemovedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIComponentSecurityOauthScopeUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),
		// APISecurityUpdatedCheck
		newBackwardCompatibilityRule(l, APISecurityRemovedCheckId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APISecurityAddedCheckId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APISecurityScopeAddedId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APISecurityScopeRemovedId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIGlobalSecurityRemovedCheckId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIGlobalSecurityAddedCheckId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIGlobalSecurityScopeAddedId, INFO, true, APISecurityUpdatedCheck),
		newBackwardCompatibilityRule(l, APIGlobalSecurityScopeRemovedId, INFO, true, APISecurityUpdatedCheck),
		// Stability Descreased Check is run as part of CheckBackwardCompatibility
		newBackwardCompatibilityRule(l, APIStabilityDecreasedId, ERR, true, nil),
		// APIDeprecationCheck
		newBackwardCompatibilityRule(l, EndpointReactivatedId, INFO, true, APIDeprecationCheck),
		newBackwardCompatibilityRule(l, APIDeprecatedSunsetParseId, ERR, true, APIDeprecationCheck),
		newBackwardCompatibilityRule(l, ParseErrorId, ERR, true, APIDeprecationCheck),
		newBackwardCompatibilityRule(l, APISunsetDateTooSmallId, ERR, true, APIDeprecationCheck),
		newBackwardCompatibilityRule(l, EndpointDeprecatedId, INFO, true, APIDeprecationCheck),
		// APIRemovedCheck
		newBackwardCompatibilityRule(l, APIPathRemovedWithoutDeprecationId, ERR, true, APIRemovedCheck),
		newBackwardCompatibilityRule(l, APIPathSunsetParseId, ERR, true, APIRemovedCheck),
		newBackwardCompatibilityRule(l, APIPathRemovedBeforeSunsetId, ERR, true, APIRemovedCheck),
		newBackwardCompatibilityRule(l, APIRemovedWithoutDeprecationId, ERR, true, APIRemovedCheck),
		newBackwardCompatibilityRule(l, APIRemovedBeforeSunsetId, ERR, true, APIRemovedCheck),
		// APISunsetChangedCheck
		newBackwardCompatibilityRule(l, APISunsetDeletedId, ERR, true, APISunsetChangedCheck),
		newBackwardCompatibilityRule(l, APISunsetDateChangedTooSmallId, ERR, true, APISunsetChangedCheck),
		// AddedRequiredRequestBodyCheck
		newBackwardCompatibilityRule(l, AddedRequiredRequestBodyId, ERR, true, AddedRequiredRequestBodyCheck),
		// NewRequestNonPathDefaultParameterCheck
		newBackwardCompatibilityRule(l, NewRequiredRequestDefaultParameterToExistingPathId, ERR, true, NewRequestNonPathDefaultParameterCheck),
		newBackwardCompatibilityRule(l, NewOptionalRequestDefaultParameterToExistingPathId, INFO, true, NewRequestNonPathDefaultParameterCheck),
		// NewRequestNonPathParameterCheck
		newBackwardCompatibilityRule(l, NewRequiredRequestParameterId, ERR, true, NewRequestNonPathParameterCheck),
		newBackwardCompatibilityRule(l, NewOptionalRequestParameterId, INFO, true, NewRequestNonPathParameterCheck),
		// NewRequestPathParameterCheck
		newBackwardCompatibilityRule(l, NewRequestPathParameterId, ERR, true, NewRequestPathParameterCheck),
		// NewRequiredRequestHeaderPropertyCheck
		newBackwardCompatibilityRule(l, NewRequiredRequestHeaderPropertyId, ERR, true, NewRequiredRequestHeaderPropertyCheck),
		// RequestBodyBecameEnumCheck
		newBackwardCompatibilityRule(l, RequestBodyBecameEnumId, ERR, true, RequestBodyBecameEnumCheck),
		// RequestBodyMediaTypeChangedCheck
		newBackwardCompatibilityRule(l, RequestBodyMediaTypeAddedId, INFO, true, RequestBodyMediaTypeChangedCheck),
		newBackwardCompatibilityRule(l, RequestBodyMediaTypeRemovedId, ERR, true, RequestBodyMediaTypeChangedCheck),
		// RequestBodyRequiredUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyBecameOptionalId, INFO, true, RequestBodyRequiredUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyBecameRequiredId, ERR, true, RequestBodyRequiredUpdatedCheck),
		// RequestDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyDiscriminatorAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDiscriminatorRemovedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDiscriminatorPropertyNameChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDiscriminatorMappingAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDiscriminatorMappingDeletedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDiscriminatorMappingChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDiscriminatorAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDiscriminatorRemovedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDiscriminatorPropertyNameChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDiscriminatorMappingAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDiscriminatorMappingDeletedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDiscriminatorMappingChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),
		// RequestHeaderPropertyBecameEnumCheck
		newBackwardCompatibilityRule(l, RequestHeaderPropertyBecameEnumId, ERR, true, RequestHeaderPropertyBecameEnumCheck),
		// RequestHeaderPropertyBecameRequiredCheck
		newBackwardCompatibilityRule(l, RequestHeaderPropertyBecameRequiredId, ERR, true, RequestHeaderPropertyBecameRequiredCheck),
		// RequestParameterBecameEnumCheck
		newBackwardCompatibilityRule(l, RequestParameterBecameEnumId, ERR, true, RequestParameterBecameEnumCheck),
		// RequestParameterDefaultValueChangedCheck
		newBackwardCompatibilityRule(l, RequestParameterDefaultValueChangedId, ERR, true, RequestParameterDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestParameterDefaultValueAddedId, ERR, true, RequestParameterDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestParameterDefaultValueRemovedId, ERR, true, RequestParameterDefaultValueChangedCheck),
		// RequestParameterEnumValueUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterEnumValueAddedId, INFO, true, RequestParameterEnumValueUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterEnumValueRemovedId, ERR, true, RequestParameterEnumValueUpdatedCheck),
		// RequestParameterMaxItemsUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterMaxItemsIncreasedId, INFO, true, RequestParameterMaxItemsUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterMaxItemsDecreasedId, ERR, true, RequestParameterMaxItemsUpdatedCheck),
		// RequestParameterMaxLengthSetCheck
		newBackwardCompatibilityRule(l, RequestParameterMaxLengthSetId, WARN, true, RequestParameterMaxLengthSetCheck),
		// RequestParameterMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterMaxLengthIncreasedId, INFO, true, RequestParameterMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterMaxLengthDecreasedId, ERR, true, RequestParameterMaxLengthUpdatedCheck),
		// RequestParameterMaxSetCheck
		newBackwardCompatibilityRule(l, RequestParameterMaxSetId, WARN, true, RequestParameterMaxSetCheck),
		// RequestParameterMaxUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterMaxIncreasedId, INFO, true, RequestParameterMaxUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterMaxDecreasedId, ERR, true, RequestParameterMaxUpdatedCheck),
		// RequestParameterMinItemsSetCheck
		newBackwardCompatibilityRule(l, RequestParameterMinItemsSetId, WARN, true, RequestParameterMinItemsSetCheck),
		// RequestParameterMinItemsUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterMinItemsIncreasedId, ERR, true, RequestParameterMinItemsUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterMinItemsDecreasedId, INFO, true, RequestParameterMinItemsUpdatedCheck),
		// RequestParameterMinLengthUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterMinLengthIncreasedId, ERR, true, RequestParameterMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterMinLengthDecreasedId, INFO, true, RequestParameterMinLengthUpdatedCheck),
		// RequestParameterMinSetCheck
		newBackwardCompatibilityRule(l, RequestParameterMinSetId, WARN, true, RequestParameterMinSetCheck),
		// RequestParameterMinUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterMinIncreasedId, ERR, true, RequestParameterMinUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterMinDecreasedId, INFO, true, RequestParameterMinUpdatedCheck),
		// RequestParameterPatternAddedOrChangedCheck
		newBackwardCompatibilityRule(l, RequestParameterPatternAddedId, WARN, true, RequestParameterPatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(l, RequestParameterPatternRemovedId, INFO, true, RequestParameterPatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(l, RequestParameterPatternChangedId, WARN, true, RequestParameterPatternAddedOrChangedCheck),
		// RequestParameterRemovedCheck
		newBackwardCompatibilityRule(l, RequestParameterRemovedId, WARN, true, RequestParameterRemovedCheck),
		// RequestParameterRequiredValueUpdatedCheck
		newBackwardCompatibilityRule(l, RequestParameterBecomeRequiredId, ERR, true, RequestParameterRequiredValueUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestParameterBecomeOptionalId, INFO, true, RequestParameterRequiredValueUpdatedCheck),
		// RequestParameterTypeChangedCheck
		newBackwardCompatibilityRule(l, RequestParameterTypeChangedId, ERR, true, RequestParameterTypeChangedCheck),
		// RequestParameterXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(l, UnparsableParameterFromXExtensibleEnumId, ERR, true, RequestParameterXExtensibleEnumValueRemovedCheck),
		newBackwardCompatibilityRule(l, UnparsableParameterToXExtensibleEnumId, ERR, true, RequestParameterXExtensibleEnumValueRemovedCheck),
		newBackwardCompatibilityRule(l, RequestParameterXExtensibleEnumValueRemovedId, ERR, true, RequestParameterXExtensibleEnumValueRemovedCheck),
		// RequestPropertyAllOfUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyAllOfAddedId, ERR, true, RequestPropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyAllOfRemovedId, WARN, true, RequestPropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyAllOfAddedId, ERR, true, RequestPropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyAllOfRemovedId, WARN, true, RequestPropertyAllOfUpdatedCheck),
		// RequestPropertyAnyOfUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyAnyOfAddedId, INFO, true, RequestPropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyAnyOfRemovedId, ERR, true, RequestPropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyAnyOfAddedId, INFO, true, RequestPropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyAnyOfRemovedId, ERR, true, RequestPropertyAnyOfUpdatedCheck),
		// RequestPropertyBecameEnumCheck
		newBackwardCompatibilityRule(l, RequestPropertyBecameEnumId, ERR, true, RequestPropertyBecameEnumCheck),
		// RequestPropertyBecameNotNullableCheck
		newBackwardCompatibilityRule(l, RequestBodyBecomeNotNullableId, ERR, true, RequestPropertyBecameNotNullableCheck),
		newBackwardCompatibilityRule(l, RequestBodyBecomeNullableId, INFO, true, RequestPropertyBecameNotNullableCheck),
		newBackwardCompatibilityRule(l, RequestPropertyBecomeNotNullableId, ERR, true, RequestPropertyBecameNotNullableCheck),
		newBackwardCompatibilityRule(l, RequestPropertyBecomeNullableId, INFO, true, RequestPropertyBecameNotNullableCheck),
		// RequestPropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(l, RequestBodyDefaultValueAddedId, INFO, true, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDefaultValueRemovedId, INFO, true, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestBodyDefaultValueChangedId, INFO, true, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDefaultValueAddedId, INFO, true, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDefaultValueRemovedId, INFO, true, RequestPropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyDefaultValueChangedId, INFO, true, RequestPropertyDefaultValueChangedCheck),
		// RequestPropertyEnumValueUpdatedCheck
		newBackwardCompatibilityRule(l, RequestPropertyEnumValueRemovedId, ERR, true, RequestPropertyEnumValueUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyEnumValueAddedId, INFO, true, RequestPropertyEnumValueUpdatedCheck),
		// RequestPropertyMaxDecreasedCheck
		newBackwardCompatibilityRule(l, RequestBodyMaxDecreasedId, ERR, true, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(l, RequestBodyMaxIncreasedId, INFO, true, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMaxDecreasedId, ERR, true, RequestPropertyMaxDecreasedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMaxIncreasedId, INFO, true, RequestPropertyMaxDecreasedCheck),
		// RequestPropertyMaxLengthSetCheck
		newBackwardCompatibilityRule(l, RequestBodyMaxLengthSetId, WARN, true, RequestPropertyMaxLengthSetCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMaxLengthSetId, WARN, true, RequestPropertyMaxLengthSetCheck),
		// RequestPropertyMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyMaxLengthDecreasedId, ERR, true, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyMaxLengthIncreasedId, INFO, true, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMaxLengthDecreasedId, ERR, true, RequestPropertyMaxLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMaxLengthIncreasedId, INFO, true, RequestPropertyMaxLengthUpdatedCheck),
		// RequestPropertyMaxSetCheck
		newBackwardCompatibilityRule(l, RequestBodyMaxSetId, WARN, true, RequestPropertyMaxSetCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMaxSetId, WARN, true, RequestPropertyMaxSetCheck),
		// RequestPropertyMinIncreasedCheck
		newBackwardCompatibilityRule(l, RequestBodyMinIncreasedId, ERR, true, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(l, RequestBodyMinDecreasedId, INFO, true, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinIncreasedId, ERR, true, RequestPropertyMinIncreasedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinDecreasedId, INFO, true, RequestPropertyMinIncreasedCheck),
		// RequestPropertyMinItemsIncreasedCheck
		newBackwardCompatibilityRule(l, RequestBodyMinItemsIncreasedId, ERR, true, RequestPropertyMinItemsIncreasedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinItemsIncreasedId, ERR, true, RequestPropertyMinItemsIncreasedCheck),
		// RequestPropertyMinItemsSetCheck
		newBackwardCompatibilityRule(l, RequestBodyMinItemsSetId, WARN, true, RequestPropertyMinItemsSetCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinItemsSetId, WARN, true, RequestPropertyMinItemsSetCheck),
		// RequestPropertyMinLengthUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyMinLengthIncreasedId, ERR, true, RequestPropertyMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyMinLengthDecreasedId, INFO, true, RequestPropertyMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinLengthIncreasedId, ERR, true, RequestPropertyMinLengthUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinLengthDecreasedId, INFO, true, RequestPropertyMinLengthUpdatedCheck),
		// RequestPropertyMinSetCheck
		newBackwardCompatibilityRule(l, RequestBodyMinSetId, WARN, true, RequestPropertyMinSetCheck),
		newBackwardCompatibilityRule(l, RequestPropertyMinSetId, WARN, true, RequestPropertyMinSetCheck),
		// RequestPropertyOneOfUpdatedCheck
		newBackwardCompatibilityRule(l, RequestBodyOneOfAddedId, INFO, true, RequestPropertyOneOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestBodyOneOfRemovedId, ERR, true, RequestPropertyOneOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyOneOfAddedId, INFO, true, RequestPropertyOneOfUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyOneOfRemovedId, ERR, true, RequestPropertyOneOfUpdatedCheck),
		// RequestPropertyPatternUpdatedCheck
		newBackwardCompatibilityRule(l, RequestPropertyPatternRemovedId, INFO, true, RequestPropertyPatternUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyPatternAddedId, WARN, true, RequestPropertyPatternUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyPatternChangedId, WARN, true, RequestPropertyPatternUpdatedCheck),
		// RequestPropertyRequiredUpdatedCheck
		newBackwardCompatibilityRule(l, RequestPropertyBecameRequiredId, ERR, true, RequestPropertyRequiredUpdatedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyBecameOptionalId, INFO, true, RequestPropertyRequiredUpdatedCheck),
		// RequestPropertyTypeChangedCheck
		newBackwardCompatibilityRule(l, RequestBodyTypeChangedId, ERR, true, RequestPropertyTypeChangedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyTypeChangedId, INFO, true, RequestPropertyTypeChangedCheck),
		// RequestPropertyUpdatedCheck
		newBackwardCompatibilityRule(l, RequestPropertyRemovedId, WARN, true, RequestPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, NewRequiredRequestPropertyId, ERR, true, RequestPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, NewOptionalRequestPropertyId, INFO, true, RequestPropertyUpdatedCheck),
		// RequestPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(l, RequestOptionalPropertyBecameNonWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestOptionalPropertyBecameWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestOptionalPropertyBecameReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestOptionalPropertyBecameNonReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestRequiredPropertyBecameNonWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestRequiredPropertyBecameWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestRequiredPropertyBecameReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, RequestRequiredPropertyBecameNonReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),
		// RequestPropertyXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(l, UnparseablePropertyFromXExtensibleEnumId, ERR, true, RequestPropertyXExtensibleEnumValueRemovedCheck),
		newBackwardCompatibilityRule(l, UnparseablePropertyToXExtensibleEnumId, ERR, true, RequestPropertyXExtensibleEnumValueRemovedCheck),
		newBackwardCompatibilityRule(l, RequestPropertyXExtensibleEnumValueRemovedId, ERR, true, RequestPropertyXExtensibleEnumValueRemovedCheck),
		// ResponseDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseBodyDiscriminatorAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDiscriminatorRemovedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDiscriminatorPropertyNameChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDiscriminatorMappingAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDiscriminatorMappingDeletedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDiscriminatorMappingChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDiscriminatorAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDiscriminatorRemovedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDiscriminatorPropertyNameChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDiscriminatorMappingAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDiscriminatorMappingDeletedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDiscriminatorMappingChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),
		// ResponseHeaderBecameOptionalCheck
		newBackwardCompatibilityRule(l, ResponseHeaderBecameOptionalId, ERR, true, ResponseHeaderBecameOptionalCheck),
		// ResponseHeaderRemovedCheck
		newBackwardCompatibilityRule(l, RequiredResponseHeaderRemovedId, ERR, true, ResponseHeaderRemovedCheck),
		newBackwardCompatibilityRule(l, OptionalResponseHeaderRemovedId, WARN, true, ResponseHeaderRemovedCheck),
		// ResponseMediaTypeUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseMediaTypeUpdatedId, ERR, true, ResponseMediaTypeUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseMediaTypeAddedId, INFO, true, ResponseMediaTypeUpdatedCheck),
		// ResponseOptionalPropertyUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseOptionalPropertyRemovedId, WARN, true, ResponseOptionalPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseOptionalWriteOnlyPropertyRemovedId, INFO, true, ResponseOptionalPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseOptionalPropertyAddedId, INFO, true, ResponseOptionalPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseOptionalWriteOnlyPropertyAddedId, INFO, true, ResponseOptionalPropertyUpdatedCheck),
		// ResponseOptionalPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(l, ResponseOptionalPropertyBecameNonWriteOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, ResponseOptionalPropertyBecameWriteOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, ResponseOptionalPropertyBecameReadOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, ResponseOptionalPropertyBecameNonReadOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),
		// ResponsePatternAddedOrChangedCheck
		newBackwardCompatibilityRule(l, ResponsePropertyPatternAddedId, INFO, true, ResponsePatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyPatternChangedId, INFO, true, ResponsePatternAddedOrChangedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyPatternRemovedId, INFO, true, ResponsePatternAddedOrChangedCheck),
		// ResponsePropertyAllOfUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseBodyAllOfAddedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyAllOfRemovedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyAllOfAddedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyAllOfRemovedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),
		// ResponsePropertyAnyOfUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseBodyAnyOfAddedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyAnyOfRemovedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyAnyOfAddedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyAnyOfRemovedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),
		// ResponsePropertyBecameNullableCheck
		newBackwardCompatibilityRule(l, ResponsePropertyBecameNullableId, ERR, true, ResponsePropertyBecameNullableCheck),
		newBackwardCompatibilityRule(l, ResponseBodyBecameNullableId, ERR, true, ResponsePropertyBecameNullableCheck),
		// ResponsePropertyBecameOptionalCheck
		newBackwardCompatibilityRule(l, ResponsePropertyBecameOptionalId, ERR, true, ResponsePropertyBecameOptionalCheck),
		newBackwardCompatibilityRule(l, ResponseWriteOnlyPropertyBecameOptionalId, ERR, true, ResponsePropertyBecameOptionalCheck),
		// ResponsePropertyBecameRequiredCheck
		newBackwardCompatibilityRule(l, ResponsePropertyBecameRequiredId, INFO, true, ResponsePropertyBecameRequiredCheck),
		newBackwardCompatibilityRule(l, ResponseWriteOnlyPropertyBecameRequiredId, INFO, true, ResponsePropertyBecameRequiredCheck),
		// ResponsePropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(l, ResponseBodyDefaultValueAddedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDefaultValueRemovedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, ResponseBodyDefaultValueChangedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDefaultValueAddedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDefaultValueRemovedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyDefaultValueChangedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),
		// ResponsePropertyEnumValueAddedCheck
		newBackwardCompatibilityRule(l, ResponsePropertyEnumValueAddedId, WARN, true, ResponsePropertyEnumValueAddedCheck),
		newBackwardCompatibilityRule(l, ResponseWriteOnlyPropertyEnumValueAddedId, INFO, true, ResponsePropertyEnumValueAddedCheck),
		// ResponsePropertyMaxIncreasedCheck
		newBackwardCompatibilityRule(l, ResponseBodyMaxIncreasedId, ERR, true, ResponsePropertyMaxIncreasedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMaxIncreasedId, ERR, true, ResponsePropertyMaxIncreasedCheck),
		// ResponsePropertyMaxLengthIncreasedCheck
		newBackwardCompatibilityRule(l, ResponseBodyMaxLengthIncreasedId, ERR, true, ResponsePropertyMaxLengthIncreasedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMaxLengthIncreasedId, ERR, true, ResponsePropertyMaxLengthIncreasedCheck),
		// ResponsePropertyMaxLengthUnsetCheck
		newBackwardCompatibilityRule(l, ResponseBodyMaxLengthUnsetId, ERR, true, ResponsePropertyMaxLengthUnsetCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMaxLengthUnsetId, ERR, true, ResponsePropertyMaxLengthUnsetCheck),
		// ResponsePropertyMinDecreasedCheck
		newBackwardCompatibilityRule(l, ResponseBodyMinDecreasedId, ERR, true, ResponsePropertyMinDecreasedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMinDecreasedId, ERR, true, ResponsePropertyMinDecreasedCheck),
		// ResponsePropertyMinItemsDecreasedCheck
		newBackwardCompatibilityRule(l, ResponseBodyMinItemsDecreasedId, ERR, true, ResponsePropertyMinItemsDecreasedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMinItemsDecreasedId, ERR, true, ResponsePropertyMinItemsDecreasedCheck),
		// ResponsePropertyMinItemsUnsetCheck
		newBackwardCompatibilityRule(l, ResponseBodyMinItemsUnsetId, ERR, true, ResponsePropertyMinItemsUnsetCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMinItemsUnsetId, ERR, true, ResponsePropertyMinItemsUnsetCheck),
		// ResponsePropertyMinLengthDecreasedCheck
		newBackwardCompatibilityRule(l, ResponseBodyMinLengthDecreasedId, ERR, true, ResponsePropertyMinLengthDecreasedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyMinLengthDecreasedId, ERR, true, ResponsePropertyMinLengthDecreasedCheck),
		// ResponsePropertyOneOfUpdated
		newBackwardCompatibilityRule(l, ResponseBodyOneOfAddedId, INFO, true, ResponsePropertyOneOfUpdated),
		newBackwardCompatibilityRule(l, ResponseBodyOneOfRemovedId, INFO, true, ResponsePropertyOneOfUpdated),
		newBackwardCompatibilityRule(l, ResponsePropertyOneOfAddedId, INFO, true, ResponsePropertyOneOfUpdated),
		newBackwardCompatibilityRule(l, ResponsePropertyOneOfRemovedId, INFO, true, ResponsePropertyOneOfUpdated),
		// ResponsePropertyTypeChangedCheck
		newBackwardCompatibilityRule(l, ResponseBodyTypeChangedId, ERR, true, ResponsePropertyTypeChangedCheck),
		newBackwardCompatibilityRule(l, ResponsePropertyTypeChangedId, ERR, true, ResponsePropertyTypeChangedCheck),
		// ResponseRequiredPropertyUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseRequiredPropertyRemovedId, ERR, true, ResponseRequiredPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseRequiredWriteOnlyPropertyRemovedId, INFO, true, ResponseRequiredPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseRequiredPropertyAddedId, ERR, true, ResponseRequiredPropertyUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseRequiredWriteOnlyPropertyAddedId, INFO, true, ResponseRequiredPropertyUpdatedCheck),
		// ResponseRequiredPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(l, ResponseRequiredPropertyBecameNonWriteOnlyId, WARN, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, ResponseRequiredPropertyBecameWriteOnlyId, INFO, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, ResponseRequiredPropertyBecameReadOnlyId, INFO, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		newBackwardCompatibilityRule(l, ResponseRequiredPropertyBecameNonReadOnlyId, INFO, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),
		// ResponseSuccessStatusUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseSuccessStatusRemovedId, ERR, true, ResponseSuccessStatusUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseSuccessStatusAddedId, ERR, true, ResponseSuccessStatusUpdatedCheck),
		// ResponseNonSuccessStatusUpdatedCheck
		newBackwardCompatibilityRule(l, ResponseNonSuccessStatusRemovedId, ERR, false, ResponseNonSuccessStatusUpdatedCheck),
		newBackwardCompatibilityRule(l, ResponseNonSuccessStatusAddedId, INFO, false, ResponseNonSuccessStatusUpdatedCheck),
		// APIOperationIdUpdatedCheck
		newBackwardCompatibilityRule(l, APIOperationIdRemovedId, ERR, false, APIOperationIdUpdatedCheck),
		newBackwardCompatibilityRule(l, APIOperationIdAddId, INFO, false, APIOperationIdUpdatedCheck),
		// APITagUpdatedCheck
		newBackwardCompatibilityRule(l, APITagRemovedId, ERR, false, APITagUpdatedCheck),
		newBackwardCompatibilityRule(l, APITagAddedId, INFO, false, APITagUpdatedCheck),
		// APIComponentsSchemaRemovedCheck
		newBackwardCompatibilityRule(l, APISchemasRemovedId, ERR, false, APIComponentsSchemaRemovedCheck),
		// ResponseParameterEnumValueRemovedCheck
		newBackwardCompatibilityRule(l, ResponsePropertyEnumValueRemovedId, ERR, false, ResponseParameterEnumValueRemovedCheck),
		// ResponseMediaTypeEnumValueRemovedCheck
		newBackwardCompatibilityRule(l, ResponseMediaTypeEnumValueRemovedId, ERR, false, ResponseMediaTypeEnumValueRemovedCheck),
		// RequestBodyEnumValueRemovedCheck
		newBackwardCompatibilityRule(l, RequestBodyEnumValueRemovedId, ERR, false, RequestBodyEnumValueRemovedCheck),
	}
}

func GetOptionalRules(l Localizer) []BackwardCompatibilityRule {

	result := []BackwardCompatibilityRule{}
	for _, rule := range GetAllRules(l) {
		if rule.Required {
			continue
		}

		if rule.Level == INFO {
			// rules with level INFO are not breaking
			continue
		}

		result = append(result, rule)
	}
	return result
}

func GetRequiredRules(l Localizer) []BackwardCompatibilityRule {

	result := []BackwardCompatibilityRule{}
	for _, rule := range GetAllRules(l) {
		if rule.Required {
			result = append(result, rule)
		}
	}
	return result
}
