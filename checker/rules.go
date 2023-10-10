package checker

func newBackwardCompatibilityRule(id string, level Level, required bool) BackwardCompatibilityRule {
	return BackwardCompatibilityRule{
		Id:          id,
		Level:       level,
		Description: id + "-description", // TODO: rule descriptions are still missing as of now
		Required:    required,
	}
}

func GetAllRules() []BackwardCompatibilityRule {
	return []BackwardCompatibilityRule{
		// APIAddedCheck
		newBackwardCompatibilityRule(EndpointAddedId, INFO, true),
		// APIComponentsSecurityUpdatedCheck, APISecurityUpdatedCheck
		newBackwardCompatibilityRule(APIComponentsSecurityRemovedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentsSecurityAddedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentsSecurityComponentOauthUrlUpdatedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentsSecurityTypeUpdatedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentsSecurityOauthTokenUrlUpdatedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeAddedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeRemovedId, INFO, true),
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeUpdatedId, INFO, true),
		// APIDeprecationCheck
		newBackwardCompatibilityRule(EndpointReactivatedId, INFO, true),
		newBackwardCompatibilityRule(APIDeprecatedSunsetParseId, ERR, true),
		newBackwardCompatibilityRule(ParseErrorId, ERR, true),
		newBackwardCompatibilityRule(APISunsetDateTooSmallId, ERR, true),
		newBackwardCompatibilityRule(EndpointDeprecatedId, INFO, true),
		// APIRemovedCheck
		newBackwardCompatibilityRule(APIPathRemovedWithoutDeprecationId, ERR, true),
		newBackwardCompatibilityRule(APIPathSunsetParseId, ERR, true),
		newBackwardCompatibilityRule(APIPathRemovedBeforeSunsetId, ERR, true),
		newBackwardCompatibilityRule(APIRemovedWithoutDeprecationId, ERR, true),
		newBackwardCompatibilityRule(APIRemovedBeforeSunsetId, ERR, true),
		// APISunsetChangedCheck
		newBackwardCompatibilityRule(APISunsetDeletedId, ERR, true),
		newBackwardCompatibilityRule(APISunsetDateChangedTooSmallId, ERR, true),
		// AddedRequiredRequestBodyCheck
		newBackwardCompatibilityRule(AddedRequiredRequestBodyId, ERR, true),
		// NewRequestNonPathDefaultParameterCheck
		newBackwardCompatibilityRule(NewRequiredRequestDefaultParameterToExistingPathId, ERR, true),
		newBackwardCompatibilityRule(NewOptionalRequestDefaultParameterToExistingPathId, INFO, true),
		// NewRequestNonPathParameterCheck
		newBackwardCompatibilityRule(NewRequiredRequestParameterId, ERR, true),
		newBackwardCompatibilityRule(NewOptionalRequestParameterId, INFO, true),
		// NewRequestPathParameterCheck
		newBackwardCompatibilityRule(NewRequestPathParameterId, ERR, true),
		// NewRequiredRequestHeaderPropertyCheck
		newBackwardCompatibilityRule(NewRequiredRequestHeaderPropertyId, ERR, true),
		// RequestBodyBecameEnumCheck
		newBackwardCompatibilityRule(RequestBodyBecameEnumId, ERR, true),
		// RequestBodyMediaTypeChangedCheck
		newBackwardCompatibilityRule(RequestBodyMediaTypeAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyMediaTypeRemovedId, ERR, true),
		// RequestBodyRequiredUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyBecameOptionalId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyBecameRequiredId, ERR, true),
		// RequestDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyDiscriminatorAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorRemovedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorPropertyNameChangedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingDeletedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingChangedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorRemovedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorPropertyNameChangedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingDeletedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingChangedId, INFO, true),
		// RequestHeaderPropertyBecameEnumCheck
		newBackwardCompatibilityRule(RequestHeaderPropertyBecameEnumId, ERR, true),
		// RequestHeaderPropertyBecameRequiredCheck
		newBackwardCompatibilityRule(RequestHeaderPropertyBecameRequiredId, ERR, true),
		// RequestParameterBecameEnumCheck
		newBackwardCompatibilityRule(RequestParameterBecameEnumId, ERR, true),
		// RequestParameterDefaultValueChangedCheck
		newBackwardCompatibilityRule(RequestParameterDefaultValueChangedId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterDefaultValueAddedId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterDefaultValueRemovedId, ERR, true),
		// RequestParameterEnumValueUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterEnumValueAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestParameterEnumValueRemovedId, ERR, true),
		// RequestParameterMaxItemsUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxItemsIncreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestParameterMaxItemsDecreasedId, ERR, true),
		// RequestParameterMaxLengthSetCheck
		newBackwardCompatibilityRule(RequestParameterMaxLengthSetId, WARN, true),
		// RequestParameterMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxLengthIncreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestParameterMaxLengthDecreasedId, ERR, true),
		// RequestParameterMaxSetCheck
		newBackwardCompatibilityRule(RequestParameterMaxSetId, WARN, true),
		// RequestParameterMaxUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxIncreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestParameterMaxDecreasedId, ERR, true),
		// RequestParameterMinItemsSetCheck
		newBackwardCompatibilityRule(RequestParameterMinItemsSetId, WARN, true),
		// RequestParameterMinItemsUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinItemsIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterMinItemsDecreasedId, INFO, true),
		// RequestParameterMinLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinLengthIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterMinLengthDecreasedId, INFO, true),
		// RequestParameterMinSetCheck
		newBackwardCompatibilityRule(RequestParameterMinSetId, WARN, true),
		// RequestParameterMinUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterMinDecreasedId, INFO, true),
		// RequestParameterPatternAddedOrChangedCheck
		newBackwardCompatibilityRule(RequestParameterPatternAddedId, WARN, true),
		newBackwardCompatibilityRule(RequestParameterPatternRemovedId, INFO, true),
		newBackwardCompatibilityRule(RequestParameterPatternChangedId, WARN, true),
		// RequestParameterRemovedCheck
		newBackwardCompatibilityRule(RequestParameterRemovedId, WARN, true),
		// RequestParameterRequiredValueUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterBecomeRequiredId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterBecomeOptionalId, INFO, true),
		// RequestParameterTypeChangedCheck
		newBackwardCompatibilityRule(RequestParameterTypeChangedId, ERR, true),
		// RequestParameterXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(UnparsableParameterFromXExtensibleEnumId, ERR, true),
		newBackwardCompatibilityRule(UnparsableParameterToXExtensibleEnumId, ERR, true),
		newBackwardCompatibilityRule(RequestParameterXExtensibleEnumValueRemovedId, ERR, true),
		// RequestPropertyAllOfUpdated
		newBackwardCompatibilityRule(RequestBodyAllOfAddedId, ERR, true),
		newBackwardCompatibilityRule(RequestBodyAllOfRemovedId, WARN, true),
		newBackwardCompatibilityRule(RequestPropertyAllOfAddedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyAllOfRemovedId, WARN, true),
		// RequestPropertyAnyOfUpdated
		newBackwardCompatibilityRule(RequestBodyAnyOfAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyAnyOfRemovedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyAnyOfAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyAnyOfRemovedId, ERR, true),
		// RequestPropertyBecameEnumCheck
		newBackwardCompatibilityRule(RequestPropertyBecameEnumId, ERR, true),
		// RequestPropertyBecameNotNullableCheck
		newBackwardCompatibilityRule(RequestBodyBecomeNotNullableId, ERR, true),
		newBackwardCompatibilityRule(RequestBodyBecomeNullableId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyBecomeNotNullableId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyBecomeNullableId, INFO, true),
		// RequestPropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(RequestBodyDefaultValueAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDefaultValueRemovedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyDefaultValueChangedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDefaultValueAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDefaultValueRemovedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyDefaultValueChangedId, INFO, true),
		// RequestPropertyEnumValueUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyEnumValueRemovedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyEnumValueAddedId, INFO, true),
		// RequestPropertyMaxDecreasedCheck
		newBackwardCompatibilityRule(RequestBodyMaxDecreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestBodyMaxIncreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyMaxDecreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyMaxIncreasedId, INFO, true),
		// RequestPropertyMaxLengthSetCheck
		newBackwardCompatibilityRule(RequestBodyMaxLengthSetId, WARN, true),
		newBackwardCompatibilityRule(RequestPropertyMaxLengthSetId, WARN, true),
		// RequestPropertyMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyMaxLengthDecreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestBodyMaxLengthIncreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyMaxLengthDecreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyMaxLengthIncreasedId, INFO, true),
		// RequestPropertyMaxSetCheck
		newBackwardCompatibilityRule(RequestBodyMaxSetId, WARN, true),
		newBackwardCompatibilityRule(RequestPropertyMaxSetId, WARN, true),
		// RequestPropertyMinIncreasedCheck
		newBackwardCompatibilityRule(RequestBodyMinIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestBodyMinDecreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyMinIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyMinDecreasedId, INFO, true),
		// RequestPropertyMinItemsIncreasedCheck
		newBackwardCompatibilityRule(RequestBodyMinItemsIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyMinItemsIncreasedId, ERR, true),
		// RequestPropertyMinItemsSetCheck
		newBackwardCompatibilityRule(RequestBodyMinItemsSetId, WARN, true),
		newBackwardCompatibilityRule(RequestPropertyMinItemsSetId, WARN, true),
		// RequestPropertyMinLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyMinLengthIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestBodyMinLengthDecreasedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyMinLengthIncreasedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyMinLengthDecreasedId, INFO, true),
		// RequestPropertyMinSetCheck
		newBackwardCompatibilityRule(RequestBodyMinSetId, WARN, true),
		newBackwardCompatibilityRule(RequestPropertyMinSetId, WARN, true),
		// RequestPropertyOneOfUpdated
		newBackwardCompatibilityRule(RequestBodyOneOfAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestBodyOneOfRemovedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyOneOfAddedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyOneOfRemovedId, ERR, true),
		// RequestPropertyPatternUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyPatternRemovedId, INFO, true),
		newBackwardCompatibilityRule(RequestPropertyPatternAddedId, WARN, true),
		newBackwardCompatibilityRule(RequestPropertyPatternChangedId, WARN, true),
		// RequestPropertyRequiredUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyBecameRequiredId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyBecameOptionalId, INFO, true),
		// RequestPropertyTypeChangedCheck
		newBackwardCompatibilityRule(RequestBodyTypeChangedId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyTypeChangedId, INFO, true),
		// RequestPropertyUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyRemovedId, WARN, true),
		newBackwardCompatibilityRule(NewRequiredRequestPropertyId, ERR, true),
		newBackwardCompatibilityRule(NewOptionalRequestPropertyId, INFO, true),
		// RequestPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameNonWriteOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameWriteOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameReadOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameNonReadOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameNonWriteOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameWriteOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameReadOnlyCheckId, INFO, true),
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameNonReadOnlyCheckId, INFO, true),
		// RequestPropertyXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(UnparseablePropertyFromXExtensibleEnumId, ERR, true),
		newBackwardCompatibilityRule(UnparseablePropertyToXExtensibleEnumId, ERR, true),
		newBackwardCompatibilityRule(RequestPropertyXExtensibleEnumValueRemovedId, ERR, true),
		// ResponseDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorPropertyNameChangedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingDeletedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingChangedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorPropertyNameChangedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingDeletedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingChangedId, INFO, true),
		// ResponseHeaderBecameOptional
		newBackwardCompatibilityRule(ResponseHeaderBecameOptionalId, ERR, true),
		// ResponseHeaderRemoved
		newBackwardCompatibilityRule(RequiredResponseHeaderRemovedId, ERR, true),
		newBackwardCompatibilityRule(OptionalResponseHeaderRemovedId, WARN, true),
		// ResponseMediaTypeUpdated
		newBackwardCompatibilityRule(ResponseMediaTypeUpdatedId, ERR, true),
		newBackwardCompatibilityRule(ResponseMediaTypeAddedId, INFO, true),
		// ResponseOptionalPropertyUpdatedCheck
		newBackwardCompatibilityRule(ResponseOptionalPropertyRemovedId, WARN, true),
		newBackwardCompatibilityRule(ResponseOptionalWriteOnlyPropertyRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponseOptionalPropertyAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseOptionalWriteOnlyPropertyAddedId, INFO, true),
		// ResponseOptionalPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameNonWriteOnlyId, INFO, true),
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameWriteOnlyId, INFO, true),
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameReadOnlyId, INFO, true),
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameNonReadOnlyId, INFO, true),
		// ResponsePatternAddedOrChangedCheck
		newBackwardCompatibilityRule(ResponsePropertyPatternAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyPatternChangedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyPatternRemovedId, INFO, true),
		// ResponsePropertyAllOfUpdated
		newBackwardCompatibilityRule(ResponseBodyAllOfAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyAllOfRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyAllOfAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyAllOfRemovedId, INFO, true),
		// ResponsePropertyAnyOfUpdated
		newBackwardCompatibilityRule(ResponseBodyAnyOfAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyAnyOfRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyAnyOfAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyAnyOfRemovedId, INFO, true),
		// ResponsePropertyBecameNullableCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameNullableId, ERR, true),
		newBackwardCompatibilityRule(ResponseBodyBecameNullableId, ERR, true),
		// ResponsePropertyBecameOptionalCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameOptionalId, ERR, true),
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyBecameOptionalId, ERR, true),
		// ResponsePropertyBecameRequiredCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameRequiredId, INFO, true),
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyBecameRequiredId, INFO, true),
		// ResponsePropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(ResponseBodyDefaultValueAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDefaultValueRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyDefaultValueChangedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueChangedId, INFO, true),
		// ResponsePropertyEnumValueAddedCheck
		newBackwardCompatibilityRule(ResponsePropertyEnumValueAddedId, WARN, true),
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyEnumValueAddedId, INFO, true),
		// ResponsePropertyMaxIncreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMaxIncreasedId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMaxIncreasedId, ERR, true),
		// ResponsePropertyMaxLengthIncreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMaxLengthIncreasedId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMaxLengthIncreasedId, ERR, true),
		// ResponsePropertyMaxLengthUnsetCheck
		newBackwardCompatibilityRule(ResponseBodyMaxLengthUnsetId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMaxLengthUnsetId, ERR, true),
		// ResponsePropertyMinDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinDecreasedId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMinDecreasedId, ERR, true),
		// ResponsePropertyMinItemsDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinItemsDecreasedId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMinItemsDecreasedId, ERR, true),
		// ResponsePropertyMinItemsUnsetCheck
		newBackwardCompatibilityRule(ResponseBodyMinItemsUnsetId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMinItemsUnsetId, ERR, true),
		// ResponsePropertyMinLengthDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinLengthDecreasedId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyMinLengthDecreasedId, ERR, true),
		// ResponsePropertyOneOfUpdated
		newBackwardCompatibilityRule(ResponseBodyOneOfAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponseBodyOneOfRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyOneOfAddedId, INFO, true),
		newBackwardCompatibilityRule(ResponsePropertyOneOfRemovedId, INFO, true),
		// ResponsePropertyTypeChangedCheck
		newBackwardCompatibilityRule(ResponseBodyTypeChangedId, ERR, true),
		newBackwardCompatibilityRule(ResponsePropertyTypeChangedId, ERR, true),
		// ResponseRequiredPropertyUpdatedCheck
		newBackwardCompatibilityRule(ResponseRequiredPropertyRemovedId, ERR, true),
		newBackwardCompatibilityRule(ResponseRequiredWriteOnlyPropertyRemovedId, INFO, true),
		newBackwardCompatibilityRule(ResponseRequiredPropertyAddedId, ERR, true),
		newBackwardCompatibilityRule(ResponseRequiredWriteOnlyPropertyAddedId, INFO, true),
		// ResponseRequiredPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameNonWriteOnlyId, WARN, true),
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameWriteOnlyId, INFO, true),
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameReadOnlyId, INFO, true),
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameNonReadOnlyId, INFO, true),
		// ResponseSuccessStatusUpdated / ResponseNonSuccessStatusUpdated
		newBackwardCompatibilityRule(ResponseSuccessStatusRemovedId, ERR, true),
		newBackwardCompatibilityRule(ResponseNonSuccessStatusRemovedId, INFO, false),
		// UncheckedRequestAllOfWarnCheck
		newBackwardCompatibilityRule(RequestAllOfModifiedId, WARN, true),
		// UncheckedResponseAllOfWarnCheck
		newBackwardCompatibilityRule(ResponseAllOfModifiedId, WARN, true),
		// APIOperationIdUpdatedCheck
		newBackwardCompatibilityRule(APIOperationIdRemovedId, ERR, false),
		newBackwardCompatibilityRule(APIOperationIdAddId, ERR, false),
		// APITagUpdatedCheck
		newBackwardCompatibilityRule(APITagRemovedId, ERR, false),
		newBackwardCompatibilityRule(APITagAddedId, ERR, false),
		// APIComponentsSchemaRemovedCheck
		newBackwardCompatibilityRule(APISchemasRemovedId, ERR, false),
		// ResponseParameterEnumValueRemovedCheck
		newBackwardCompatibilityRule(ResponsePropertyEnumValueRemovedId, ERR, false),
		// ResponseMediaTypeEnumValueRemovedCheck
		newBackwardCompatibilityRule(ResponseMediaTypeEnumValueRemovedId, ERR, false),
		// RequestBodyEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestBodyEnumValueRemovedId, ERR, false),
	}
}
