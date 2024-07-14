package checker

type BackwardCompatibilityRule struct {
	Id          string
	Level       Level
	Description string
	Required    bool
	Handler     BackwardCompatibilityCheck `json:"-" yaml:"-"`
}

func newBackwardCompatibilityRule(id string, level Level, required bool, handler BackwardCompatibilityCheck) BackwardCompatibilityRule {
	return BackwardCompatibilityRule{
		Id:          id,
		Level:       level,
		Description: descriptionId(id),
		Required:    required,
		Handler:     handler,
	}
}

func GetAllRules() []BackwardCompatibilityRule {
	return []BackwardCompatibilityRule{
		// APIAddedCheck
		newBackwardCompatibilityRule(EndpointAddedId, INFO, true, APIAddedCheck), //
		// APIComponentsSecurityUpdatedCheck
		newBackwardCompatibilityRule(APIComponentsSecurityRemovedId, INFO, true, APIComponentsSecurityUpdatedCheck),                  //
		newBackwardCompatibilityRule(APIComponentsSecurityAddedId, INFO, true, APIComponentsSecurityUpdatedCheck),                    //
		newBackwardCompatibilityRule(APIComponentsSecurityComponentOauthUrlUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck), //
		newBackwardCompatibilityRule(APIComponentsSecurityTypeUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),              //
		newBackwardCompatibilityRule(APIComponentsSecurityOauthTokenUrlUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),     //
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeAddedId, INFO, true, APIComponentsSecurityUpdatedCheck),           //
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeRemovedId, INFO, true, APIComponentsSecurityUpdatedCheck),         //
		newBackwardCompatibilityRule(APIComponentSecurityOauthScopeUpdatedId, INFO, true, APIComponentsSecurityUpdatedCheck),         //
		// APISecurityUpdatedCheck
		newBackwardCompatibilityRule(APISecurityRemovedCheckId, INFO, true, APISecurityUpdatedCheck),       //
		newBackwardCompatibilityRule(APISecurityAddedCheckId, INFO, true, APISecurityUpdatedCheck),         //
		newBackwardCompatibilityRule(APISecurityScopeAddedId, INFO, true, APISecurityUpdatedCheck),         //
		newBackwardCompatibilityRule(APISecurityScopeRemovedId, INFO, true, APISecurityUpdatedCheck),       //
		newBackwardCompatibilityRule(APIGlobalSecurityRemovedCheckId, INFO, true, APISecurityUpdatedCheck), //
		newBackwardCompatibilityRule(APIGlobalSecurityAddedCheckId, INFO, true, APISecurityUpdatedCheck),   //
		newBackwardCompatibilityRule(APIGlobalSecurityScopeAddedId, INFO, true, APISecurityUpdatedCheck),   //
		newBackwardCompatibilityRule(APIGlobalSecurityScopeRemovedId, INFO, true, APISecurityUpdatedCheck), //
		// Stability Descreased Check is run as part of CheckBackwardCompatibility
		newBackwardCompatibilityRule(APIStabilityDecreasedId, ERR, true, nil), //
		// APIDeprecationCheck
		newBackwardCompatibilityRule(EndpointReactivatedId, INFO, true, APIDeprecationCheck),     //
		newBackwardCompatibilityRule(APIDeprecatedSunsetParseId, ERR, true, APIDeprecationCheck), //
		newBackwardCompatibilityRule(APIInvalidStabilityLevelId, ERR, true, APIDeprecationCheck), //
		newBackwardCompatibilityRule(APISunsetDateTooSmallId, ERR, true, APIDeprecationCheck),    //
		newBackwardCompatibilityRule(EndpointDeprecatedId, INFO, true, APIDeprecationCheck),      //
		// APIRemovedCheck
		newBackwardCompatibilityRule(APIPathRemovedWithoutDeprecationId, ERR, true, APIRemovedCheck), //
		newBackwardCompatibilityRule(APIPathSunsetParseId, ERR, true, APIRemovedCheck),               //
		newBackwardCompatibilityRule(APIPathRemovedBeforeSunsetId, ERR, true, APIRemovedCheck),       //
		newBackwardCompatibilityRule(APIRemovedWithoutDeprecationId, ERR, true, APIRemovedCheck),     //
		newBackwardCompatibilityRule(APIRemovedBeforeSunsetId, ERR, true, APIRemovedCheck),           //
		// APISunsetChangedCheck
		newBackwardCompatibilityRule(APISunsetDeletedId, ERR, true, APISunsetChangedCheck),             //
		newBackwardCompatibilityRule(APISunsetDateChangedTooSmallId, ERR, true, APISunsetChangedCheck), //
		// AddedRequiredRequestBodyCheck
		newBackwardCompatibilityRule(AddedRequiredRequestBodyId, ERR, true, AddedRequestBodyCheck),  //
		newBackwardCompatibilityRule(AddedOptionalRequestBodyId, INFO, true, AddedRequestBodyCheck), //
		// NewRequestNonPathDefaultParameterCheck
		newBackwardCompatibilityRule(NewRequiredRequestDefaultParameterToExistingPathId, ERR, true, NewRequestNonPathDefaultParameterCheck),  //
		newBackwardCompatibilityRule(NewOptionalRequestDefaultParameterToExistingPathId, INFO, true, NewRequestNonPathDefaultParameterCheck), //
		// NewRequestNonPathParameterCheck
		newBackwardCompatibilityRule(NewRequiredRequestParameterId, ERR, true, NewRequestNonPathParameterCheck),  //
		newBackwardCompatibilityRule(NewOptionalRequestParameterId, INFO, true, NewRequestNonPathParameterCheck), //
		// NewRequestPathParameterCheck
		newBackwardCompatibilityRule(NewRequestPathParameterId, ERR, true, NewRequestPathParameterCheck), //
		// NewRequiredRequestHeaderPropertyCheck
		newBackwardCompatibilityRule(NewRequiredRequestHeaderPropertyId, ERR, true, NewRequiredRequestHeaderPropertyCheck), //
		// RequestBodyBecameEnumCheck
		newBackwardCompatibilityRule(RequestBodyBecameEnumId, ERR, true, RequestBodyBecameEnumCheck), //
		// RequestBodyMediaTypeChangedCheck
		newBackwardCompatibilityRule(RequestBodyMediaTypeAddedId, INFO, true, RequestBodyMediaTypeChangedCheck),  //
		newBackwardCompatibilityRule(RequestBodyMediaTypeRemovedId, ERR, true, RequestBodyMediaTypeChangedCheck), //
		// RequestBodyRequiredUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyBecameOptionalId, INFO, true, RequestBodyRequiredUpdatedCheck), //
		newBackwardCompatibilityRule(RequestBodyBecameRequiredId, ERR, true, RequestBodyRequiredUpdatedCheck),  //
		// RequestDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyDiscriminatorAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),                   //
		newBackwardCompatibilityRule(RequestBodyDiscriminatorRemovedId, INFO, true, RequestDiscriminatorUpdatedCheck),                 //
		newBackwardCompatibilityRule(RequestBodyDiscriminatorPropertyNameChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),     //
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),            //
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingDeletedId, INFO, true, RequestDiscriminatorUpdatedCheck),          //
		newBackwardCompatibilityRule(RequestBodyDiscriminatorMappingChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),          //
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),               //
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorRemovedId, INFO, true, RequestDiscriminatorUpdatedCheck),             //
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorPropertyNameChangedId, INFO, true, RequestDiscriminatorUpdatedCheck), //
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingAddedId, INFO, true, RequestDiscriminatorUpdatedCheck),        //
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingDeletedId, INFO, true, RequestDiscriminatorUpdatedCheck),      //
		newBackwardCompatibilityRule(RequestPropertyDiscriminatorMappingChangedId, INFO, true, RequestDiscriminatorUpdatedCheck),      //
		// RequestHeaderPropertyBecameEnumCheck
		newBackwardCompatibilityRule(RequestHeaderPropertyBecameEnumId, ERR, true, RequestHeaderPropertyBecameEnumCheck), //
		// RequestHeaderPropertyBecameRequiredCheck
		newBackwardCompatibilityRule(RequestHeaderPropertyBecameRequiredId, ERR, true, RequestHeaderPropertyBecameRequiredCheck), //
		// RequestParameterBecameEnumCheck
		newBackwardCompatibilityRule(RequestParameterBecameEnumId, ERR, true, RequestParameterBecameEnumCheck), //
		// RequestParameterDefaultValueChangedCheck
		newBackwardCompatibilityRule(RequestParameterDefaultValueChangedId, ERR, true, RequestParameterDefaultValueChangedCheck), //
		newBackwardCompatibilityRule(RequestParameterDefaultValueAddedId, ERR, true, RequestParameterDefaultValueChangedCheck),   //
		newBackwardCompatibilityRule(RequestParameterDefaultValueRemovedId, ERR, true, RequestParameterDefaultValueChangedCheck), //
		// RequestParameterEnumValueUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterEnumValueAddedId, INFO, true, RequestParameterEnumValueUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestParameterEnumValueRemovedId, ERR, true, RequestParameterEnumValueUpdatedCheck), //
		// RequestParameterMaxItemsUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxItemsIncreasedId, INFO, true, RequestParameterMaxItemsUpdatedCheck), //
		newBackwardCompatibilityRule(RequestParameterMaxItemsDecreasedId, ERR, true, RequestParameterMaxItemsUpdatedCheck),  //
		// RequestParameterMaxLengthSetCheck
		newBackwardCompatibilityRule(RequestParameterMaxLengthSetId, WARN, true, RequestParameterMaxLengthSetCheck), //
		// RequestParameterMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxLengthIncreasedId, INFO, true, RequestParameterMaxLengthUpdatedCheck), //
		newBackwardCompatibilityRule(RequestParameterMaxLengthDecreasedId, ERR, true, RequestParameterMaxLengthUpdatedCheck),  //
		// RequestParameterMaxSetCheck
		newBackwardCompatibilityRule(RequestParameterMaxSetId, WARN, true, RequestParameterMaxSetCheck), //
		// RequestParameterMaxUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMaxIncreasedId, INFO, true, RequestParameterMaxUpdatedCheck), //
		newBackwardCompatibilityRule(RequestParameterMaxDecreasedId, ERR, true, RequestParameterMaxUpdatedCheck),  //
		// RequestParameterMinItemsSetCheck
		newBackwardCompatibilityRule(RequestParameterMinItemsSetId, WARN, true, RequestParameterMinItemsSetCheck), //
		// RequestParameterMinItemsUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinItemsIncreasedId, ERR, true, RequestParameterMinItemsUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestParameterMinItemsDecreasedId, INFO, true, RequestParameterMinItemsUpdatedCheck), //
		// RequestParameterMinLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinLengthIncreasedId, ERR, true, RequestParameterMinLengthUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestParameterMinLengthDecreasedId, INFO, true, RequestParameterMinLengthUpdatedCheck), //
		// RequestParameterMinSetCheck
		newBackwardCompatibilityRule(RequestParameterMinSetId, WARN, true, RequestParameterMinSetCheck), //
		// RequestParameterMinUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterMinIncreasedId, ERR, true, RequestParameterMinUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestParameterMinDecreasedId, INFO, true, RequestParameterMinUpdatedCheck), //
		// RequestParameterPatternAddedOrChangedCheck
		newBackwardCompatibilityRule(RequestParameterPatternAddedId, WARN, true, RequestParameterPatternAddedOrChangedCheck),       //
		newBackwardCompatibilityRule(RequestParameterPatternRemovedId, INFO, true, RequestParameterPatternAddedOrChangedCheck),     //
		newBackwardCompatibilityRule(RequestParameterPatternChangedId, WARN, true, RequestParameterPatternAddedOrChangedCheck),     //
		newBackwardCompatibilityRule(RequestParameterPatternGeneralizedId, INFO, true, RequestParameterPatternAddedOrChangedCheck), //
		// RequestParameterRemovedCheck
		newBackwardCompatibilityRule(RequestParameterRemovedId, WARN, true, RequestParameterRemovedCheck), //
		// RequestParameterRequiredValueUpdatedCheck
		newBackwardCompatibilityRule(RequestParameterBecomeRequiredId, ERR, true, RequestParameterRequiredValueUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestParameterBecomeOptionalId, INFO, true, RequestParameterRequiredValueUpdatedCheck), //
		// RequestParameterTypeChangedCheck
		newBackwardCompatibilityRule(RequestParameterTypeChangedId, ERR, true, RequestParameterTypeChangedCheck),              //
		newBackwardCompatibilityRule(RequestParameterTypeGeneralizedId, INFO, true, RequestParameterTypeChangedCheck),         //
		newBackwardCompatibilityRule(RequestParameterPropertyTypeChangedId, WARN, true, RequestParameterTypeChangedCheck),     //
		newBackwardCompatibilityRule(RequestParameterPropertyTypeGeneralizedId, INFO, true, RequestParameterTypeChangedCheck), //
		newBackwardCompatibilityRule(RequestParameterPropertyTypeSpecializedId, ERR, true, RequestParameterTypeChangedCheck),  //
		// RequestParameterXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestParameterXExtensibleEnumValueRemovedId, ERR, true, RequestParameterXExtensibleEnumValueRemovedCheck), //
		// RequestPropertyAllOfUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyAllOfAddedId, ERR, true, RequestPropertyAllOfUpdatedCheck),        //
		newBackwardCompatibilityRule(RequestBodyAllOfRemovedId, WARN, true, RequestPropertyAllOfUpdatedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyAllOfAddedId, ERR, true, RequestPropertyAllOfUpdatedCheck),    //
		newBackwardCompatibilityRule(RequestPropertyAllOfRemovedId, WARN, true, RequestPropertyAllOfUpdatedCheck), //
		// RequestPropertyAnyOfUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyAnyOfAddedId, INFO, true, RequestPropertyAnyOfUpdatedCheck),      //
		newBackwardCompatibilityRule(RequestBodyAnyOfRemovedId, ERR, true, RequestPropertyAnyOfUpdatedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyAnyOfAddedId, INFO, true, RequestPropertyAnyOfUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestPropertyAnyOfRemovedId, ERR, true, RequestPropertyAnyOfUpdatedCheck), //
		// RequestPropertyBecameEnumCheck
		newBackwardCompatibilityRule(RequestPropertyBecameEnumId, ERR, true, RequestPropertyBecameEnumCheck), //
		// RequestPropertyBecameNotNullableCheck
		newBackwardCompatibilityRule(RequestBodyBecomeNotNullableId, ERR, true, RequestPropertyBecameNotNullableCheck),     //
		newBackwardCompatibilityRule(RequestBodyBecomeNullableId, INFO, true, RequestPropertyBecameNotNullableCheck),       //
		newBackwardCompatibilityRule(RequestPropertyBecomeNotNullableId, ERR, true, RequestPropertyBecameNotNullableCheck), //
		newBackwardCompatibilityRule(RequestPropertyBecomeNullableId, INFO, true, RequestPropertyBecameNotNullableCheck),   //
		// RequestPropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(RequestBodyDefaultValueAddedId, INFO, true, RequestPropertyDefaultValueChangedCheck),       //
		newBackwardCompatibilityRule(RequestBodyDefaultValueRemovedId, INFO, true, RequestPropertyDefaultValueChangedCheck),     //
		newBackwardCompatibilityRule(RequestBodyDefaultValueChangedId, INFO, true, RequestPropertyDefaultValueChangedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyDefaultValueAddedId, INFO, true, RequestPropertyDefaultValueChangedCheck),   //
		newBackwardCompatibilityRule(RequestPropertyDefaultValueRemovedId, INFO, true, RequestPropertyDefaultValueChangedCheck), //
		newBackwardCompatibilityRule(RequestPropertyDefaultValueChangedId, INFO, true, RequestPropertyDefaultValueChangedCheck), //
		// RequestPropertyEnumValueUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyEnumValueRemovedId, ERR, true, RequestPropertyEnumValueUpdatedCheck),          //
		newBackwardCompatibilityRule(RequestReadOnlyPropertyEnumValueRemovedId, INFO, true, RequestPropertyEnumValueUpdatedCheck), //
		newBackwardCompatibilityRule(RequestPropertyEnumValueAddedId, INFO, true, RequestPropertyEnumValueUpdatedCheck),           //
		// RequestPropertyMaxDecreasedCheck
		newBackwardCompatibilityRule(RequestBodyMaxDecreasedId, ERR, true, RequestPropertyMaxDecreasedCheck),      //
		newBackwardCompatibilityRule(RequestBodyMaxIncreasedId, INFO, true, RequestPropertyMaxDecreasedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMaxDecreasedId, ERR, true, RequestPropertyMaxDecreasedCheck),  //
		newBackwardCompatibilityRule(RequestPropertyMaxIncreasedId, INFO, true, RequestPropertyMaxDecreasedCheck), //
		// RequestPropertyMaxLengthSetCheck
		newBackwardCompatibilityRule(RequestBodyMaxLengthSetId, WARN, true, RequestPropertyMaxLengthSetCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMaxLengthSetId, WARN, true, RequestPropertyMaxLengthSetCheck), //
		// RequestPropertyMaxLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyMaxLengthDecreasedId, ERR, true, RequestPropertyMaxLengthUpdatedCheck),              //
		newBackwardCompatibilityRule(RequestBodyMaxLengthIncreasedId, INFO, true, RequestPropertyMaxLengthUpdatedCheck),             //
		newBackwardCompatibilityRule(RequestPropertyMaxLengthDecreasedId, ERR, true, RequestPropertyMaxLengthUpdatedCheck),          //
		newBackwardCompatibilityRule(RequestReadOnlyPropertyMaxLengthDecreasedId, INFO, true, RequestPropertyMaxLengthUpdatedCheck), //
		newBackwardCompatibilityRule(RequestPropertyMaxLengthIncreasedId, INFO, true, RequestPropertyMaxLengthUpdatedCheck),         //
		// RequestPropertyMaxSetCheck
		newBackwardCompatibilityRule(RequestBodyMaxSetId, WARN, true, RequestPropertyMaxSetCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMaxSetId, WARN, true, RequestPropertyMaxSetCheck), //
		// RequestPropertyMinIncreasedCheck
		newBackwardCompatibilityRule(RequestBodyMinIncreasedId, ERR, true, RequestPropertyMinIncreasedCheck),      //
		newBackwardCompatibilityRule(RequestBodyMinDecreasedId, INFO, true, RequestPropertyMinIncreasedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMinIncreasedId, ERR, true, RequestPropertyMinIncreasedCheck),  // INFO or ERR
		newBackwardCompatibilityRule(RequestPropertyMinDecreasedId, INFO, true, RequestPropertyMinIncreasedCheck), //
		// RequestPropertyMinItemsIncreasedCheck
		newBackwardCompatibilityRule(RequestBodyMinItemsIncreasedId, ERR, true, RequestPropertyMinItemsIncreasedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMinItemsIncreasedId, ERR, true, RequestPropertyMinItemsIncreasedCheck), //
		// RequestPropertyMinItemsSetCheck
		newBackwardCompatibilityRule(RequestBodyMinItemsSetId, WARN, true, RequestPropertyMinItemsSetCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMinItemsSetId, WARN, true, RequestPropertyMinItemsSetCheck), //
		// RequestPropertyMinLengthUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyMinLengthIncreasedId, ERR, true, RequestPropertyMinLengthUpdatedCheck),      //
		newBackwardCompatibilityRule(RequestBodyMinLengthDecreasedId, INFO, true, RequestPropertyMinLengthUpdatedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMinLengthIncreasedId, ERR, true, RequestPropertyMinLengthUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestPropertyMinLengthDecreasedId, INFO, true, RequestPropertyMinLengthUpdatedCheck), //
		// RequestPropertyMinSetCheck
		newBackwardCompatibilityRule(RequestBodyMinSetId, WARN, true, RequestPropertyMinSetCheck),     //
		newBackwardCompatibilityRule(RequestPropertyMinSetId, WARN, true, RequestPropertyMinSetCheck), //
		// RequestPropertyOneOfUpdatedCheck
		newBackwardCompatibilityRule(RequestBodyOneOfAddedId, INFO, true, RequestPropertyOneOfUpdatedCheck),      //
		newBackwardCompatibilityRule(RequestBodyOneOfRemovedId, ERR, true, RequestPropertyOneOfUpdatedCheck),     //
		newBackwardCompatibilityRule(RequestPropertyOneOfAddedId, INFO, true, RequestPropertyOneOfUpdatedCheck),  //
		newBackwardCompatibilityRule(RequestPropertyOneOfRemovedId, ERR, true, RequestPropertyOneOfUpdatedCheck), //
		// RequestPropertyPatternUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyPatternRemovedId, INFO, true, RequestPropertyPatternUpdatedCheck), //
		newBackwardCompatibilityRule(RequestPropertyPatternAddedId, WARN, true, RequestPropertyPatternUpdatedCheck),   //
		newBackwardCompatibilityRule(RequestPropertyPatternChangedId, WARN, true, RequestPropertyPatternUpdatedCheck), // INFO or WARN
		// RequestPropertyRequiredUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyBecameRequiredId, ERR, true, RequestPropertyRequiredUpdatedCheck),             //
		newBackwardCompatibilityRule(RequestPropertyBecameRequiredWithDefaultId, INFO, true, RequestPropertyRequiredUpdatedCheck), //
		newBackwardCompatibilityRule(RequestPropertyBecameOptionalId, INFO, true, RequestPropertyRequiredUpdatedCheck),            //
		// RequestPropertyTypeChangedCheck
		newBackwardCompatibilityRule(RequestBodyTypeChangedId, ERR, true, RequestPropertyTypeChangedCheck),      // INFO or ERR
		newBackwardCompatibilityRule(RequestPropertyTypeChangedId, INFO, true, RequestPropertyTypeChangedCheck), // INFO or ERR
		// RequestPropertyUpdatedCheck
		newBackwardCompatibilityRule(RequestPropertyRemovedId, WARN, true, RequestPropertyUpdatedCheck),                //
		newBackwardCompatibilityRule(NewRequiredRequestPropertyId, ERR, true, RequestPropertyUpdatedCheck),             //
		newBackwardCompatibilityRule(NewRequiredRequestPropertyWithDefaultId, INFO, true, RequestPropertyUpdatedCheck), //
		newBackwardCompatibilityRule(NewOptionalRequestPropertyId, INFO, true, RequestPropertyUpdatedCheck),            //
		// RequestPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameNonWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck), //
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),    //
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),     //
		newBackwardCompatibilityRule(RequestOptionalPropertyBecameNonReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),  //
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameNonWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck), //
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameWriteOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),    //
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),     //
		newBackwardCompatibilityRule(RequestRequiredPropertyBecameNonReadOnlyCheckId, INFO, true, RequestPropertyWriteOnlyReadOnlyCheck),  //
		// RequestPropertyXExtensibleEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestPropertyXExtensibleEnumValueRemovedId, ERR, true, RequestPropertyXExtensibleEnumValueRemovedCheck), //
		// ResponseDiscriminatorUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),                   //
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorRemovedId, INFO, true, ResponseDiscriminatorUpdatedCheck),                 //
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorPropertyNameChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),     //
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),            //
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingDeletedId, INFO, true, ResponseDiscriminatorUpdatedCheck),          //
		newBackwardCompatibilityRule(ResponseBodyDiscriminatorMappingChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),          //
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),               //
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorRemovedId, INFO, true, ResponseDiscriminatorUpdatedCheck),             //
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorPropertyNameChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck), //
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingAddedId, INFO, true, ResponseDiscriminatorUpdatedCheck),        //
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingDeletedId, INFO, true, ResponseDiscriminatorUpdatedCheck),      //
		newBackwardCompatibilityRule(ResponsePropertyDiscriminatorMappingChangedId, INFO, true, ResponseDiscriminatorUpdatedCheck),      //
		// ResponseHeaderBecameOptionalCheck
		newBackwardCompatibilityRule(ResponseHeaderBecameOptionalId, ERR, true, ResponseHeaderBecameOptionalCheck), //
		// ResponseHeaderRemovedCheck
		newBackwardCompatibilityRule(RequiredResponseHeaderRemovedId, ERR, true, ResponseHeaderRemovedCheck),  //
		newBackwardCompatibilityRule(OptionalResponseHeaderRemovedId, WARN, true, ResponseHeaderRemovedCheck), //
		// ResponseMediaTypeUpdatedCheck
		newBackwardCompatibilityRule(ResponseMediaTypeRemovedId, ERR, true, ResponseMediaTypeUpdatedCheck), //
		newBackwardCompatibilityRule(ResponseMediaTypeAddedId, INFO, true, ResponseMediaTypeUpdatedCheck),  //
		// ResponseOptionalPropertyUpdatedCheck
		newBackwardCompatibilityRule(ResponseOptionalPropertyRemovedId, WARN, true, ResponseOptionalPropertyUpdatedCheck),          //
		newBackwardCompatibilityRule(ResponseOptionalWriteOnlyPropertyRemovedId, INFO, true, ResponseOptionalPropertyUpdatedCheck), //
		newBackwardCompatibilityRule(ResponseOptionalPropertyAddedId, INFO, true, ResponseOptionalPropertyUpdatedCheck),            //
		newBackwardCompatibilityRule(ResponseOptionalWriteOnlyPropertyAddedId, INFO, true, ResponseOptionalPropertyUpdatedCheck),   //
		// ResponseOptionalPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameNonWriteOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck), //
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameWriteOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),    //
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameReadOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),     //
		newBackwardCompatibilityRule(ResponseOptionalPropertyBecameNonReadOnlyId, INFO, true, ResponseOptionalPropertyWriteOnlyReadOnlyCheck),  //
		// ResponsePatternAddedOrChangedCheck
		newBackwardCompatibilityRule(ResponsePropertyPatternAddedId, INFO, true, ResponsePatternAddedOrChangedCheck),   //
		newBackwardCompatibilityRule(ResponsePropertyPatternChangedId, INFO, true, ResponsePatternAddedOrChangedCheck), // shouldn't this depend on the pattern?
		newBackwardCompatibilityRule(ResponsePropertyPatternRemovedId, INFO, true, ResponsePatternAddedOrChangedCheck), //
		// ResponsePropertyAllOfUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyAllOfAddedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),       //
		newBackwardCompatibilityRule(ResponseBodyAllOfRemovedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyAllOfAddedId, INFO, true, ResponsePropertyAllOfUpdatedCheck),   //
		newBackwardCompatibilityRule(ResponsePropertyAllOfRemovedId, INFO, true, ResponsePropertyAllOfUpdatedCheck), //
		// ResponsePropertyAnyOfUpdatedCheck
		newBackwardCompatibilityRule(ResponseBodyAnyOfAddedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),       //
		newBackwardCompatibilityRule(ResponseBodyAnyOfRemovedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyAnyOfAddedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck),   //
		newBackwardCompatibilityRule(ResponsePropertyAnyOfRemovedId, INFO, true, ResponsePropertyAnyOfUpdatedCheck), //
		// ResponsePropertyBecameNullableCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameNullableId, ERR, true, ResponsePropertyBecameNullableCheck), //
		newBackwardCompatibilityRule(ResponseBodyBecameNullableId, ERR, true, ResponsePropertyBecameNullableCheck),     //
		// ResponsePropertyBecameOptionalCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameOptionalId, ERR, true, ResponsePropertyBecameOptionalCheck),          //
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyBecameOptionalId, ERR, true, ResponsePropertyBecameOptionalCheck), //
		// ResponsePropertyBecameRequiredCheck
		newBackwardCompatibilityRule(ResponsePropertyBecameRequiredId, INFO, true, ResponsePropertyBecameRequiredCheck),          //
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyBecameRequiredId, INFO, true, ResponsePropertyBecameRequiredCheck), //
		// ResponsePropertyDefaultValueChangedCheck
		newBackwardCompatibilityRule(ResponseBodyDefaultValueAddedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),       //
		newBackwardCompatibilityRule(ResponseBodyDefaultValueRemovedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),     //
		newBackwardCompatibilityRule(ResponseBodyDefaultValueChangedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueAddedId, INFO, true, ResponsePropertyDefaultValueChangedCheck),   //
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueRemovedId, INFO, true, ResponsePropertyDefaultValueChangedCheck), //
		newBackwardCompatibilityRule(ResponsePropertyDefaultValueChangedId, INFO, true, ResponsePropertyDefaultValueChangedCheck), //
		// ResponsePropertyEnumValueAddedCheck
		newBackwardCompatibilityRule(ResponsePropertyEnumValueAddedId, WARN, true, ResponsePropertyEnumValueAddedCheck),          //
		newBackwardCompatibilityRule(ResponseWriteOnlyPropertyEnumValueAddedId, INFO, true, ResponsePropertyEnumValueAddedCheck), //
		// ResponsePropertyMaxIncreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMaxIncreasedId, ERR, true, ResponsePropertyMaxIncreasedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMaxIncreasedId, ERR, true, ResponsePropertyMaxIncreasedCheck), //
		// ResponsePropertyMaxLengthIncreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMaxLengthIncreasedId, ERR, true, ResponsePropertyMaxLengthIncreasedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMaxLengthIncreasedId, ERR, true, ResponsePropertyMaxLengthIncreasedCheck), //
		// ResponsePropertyMaxLengthUnsetCheck
		newBackwardCompatibilityRule(ResponseBodyMaxLengthUnsetId, ERR, true, ResponsePropertyMaxLengthUnsetCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMaxLengthUnsetId, ERR, true, ResponsePropertyMaxLengthUnsetCheck), //
		// ResponsePropertyMinDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinDecreasedId, ERR, true, ResponsePropertyMinDecreasedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMinDecreasedId, ERR, true, ResponsePropertyMinDecreasedCheck), //
		// ResponsePropertyMinItemsDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinItemsDecreasedId, ERR, true, ResponsePropertyMinItemsDecreasedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMinItemsDecreasedId, ERR, true, ResponsePropertyMinItemsDecreasedCheck), //
		// ResponsePropertyMinItemsUnsetCheck
		newBackwardCompatibilityRule(ResponseBodyMinItemsUnsetId, ERR, true, ResponsePropertyMinItemsUnsetCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMinItemsUnsetId, ERR, true, ResponsePropertyMinItemsUnsetCheck), //
		// ResponsePropertyMinLengthDecreasedCheck
		newBackwardCompatibilityRule(ResponseBodyMinLengthDecreasedId, ERR, true, ResponsePropertyMinLengthDecreasedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyMinLengthDecreasedId, ERR, true, ResponsePropertyMinLengthDecreasedCheck), //
		// ResponsePropertyOneOfUpdated
		newBackwardCompatibilityRule(ResponseBodyOneOfAddedId, INFO, true, ResponsePropertyOneOfUpdated),       //
		newBackwardCompatibilityRule(ResponseBodyOneOfRemovedId, INFO, true, ResponsePropertyOneOfUpdated),     //
		newBackwardCompatibilityRule(ResponsePropertyOneOfAddedId, INFO, true, ResponsePropertyOneOfUpdated),   //
		newBackwardCompatibilityRule(ResponsePropertyOneOfRemovedId, INFO, true, ResponsePropertyOneOfUpdated), //
		// ResponsePropertyTypeChangedCheck
		newBackwardCompatibilityRule(ResponseBodyTypeChangedId, ERR, true, ResponsePropertyTypeChangedCheck),     //
		newBackwardCompatibilityRule(ResponsePropertyTypeChangedId, ERR, true, ResponsePropertyTypeChangedCheck), //
		// ResponseRequiredPropertyUpdatedCheck
		newBackwardCompatibilityRule(ResponseRequiredPropertyRemovedId, ERR, true, ResponseRequiredPropertyUpdatedCheck),           //
		newBackwardCompatibilityRule(ResponseRequiredWriteOnlyPropertyRemovedId, INFO, true, ResponseRequiredPropertyUpdatedCheck), //
		newBackwardCompatibilityRule(ResponseRequiredPropertyAddedId, INFO, true, ResponseRequiredPropertyUpdatedCheck),            //
		newBackwardCompatibilityRule(ResponseRequiredWriteOnlyPropertyAddedId, INFO, true, ResponseRequiredPropertyUpdatedCheck),   //
		// ResponseRequiredPropertyWriteOnlyReadOnlyCheck
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameNonWriteOnlyId, WARN, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck), //
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameWriteOnlyId, INFO, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),    //
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameReadOnlyId, INFO, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),     //
		newBackwardCompatibilityRule(ResponseRequiredPropertyBecameNonReadOnlyId, INFO, true, ResponseRequiredPropertyWriteOnlyReadOnlyCheck),  //
		// ResponseSuccessStatusUpdatedCheck
		newBackwardCompatibilityRule(ResponseSuccessStatusRemovedId, ERR, true, ResponseSuccessStatusUpdatedCheck), // INFO or ERR
		newBackwardCompatibilityRule(ResponseSuccessStatusAddedId, ERR, true, ResponseSuccessStatusUpdatedCheck),   // INFO or ERR
		// ResponseNonSuccessStatusUpdatedCheck
		newBackwardCompatibilityRule(ResponseNonSuccessStatusRemovedId, ERR, false, ResponseNonSuccessStatusUpdatedCheck), // INFO or ERR
		newBackwardCompatibilityRule(ResponseNonSuccessStatusAddedId, INFO, false, ResponseNonSuccessStatusUpdatedCheck),  // INFO or ERR
		// APIOperationIdUpdatedCheck
		newBackwardCompatibilityRule(APIOperationIdRemovedId, ERR, false, APIOperationIdUpdatedCheck), // INFO or ERR
		newBackwardCompatibilityRule(APIOperationIdAddId, INFO, false, APIOperationIdUpdatedCheck),    // INFO or ERR
		// APITagUpdatedCheck
		newBackwardCompatibilityRule(APITagRemovedId, ERR, false, APITagUpdatedCheck), // INFO or ERR
		newBackwardCompatibilityRule(APITagAddedId, INFO, false, APITagUpdatedCheck),  // INFO or ERR
		// APIComponentsSchemaRemovedCheck
		newBackwardCompatibilityRule(APISchemasRemovedId, ERR, false, APIComponentsSchemaRemovedCheck), // INFO or ERR
		// ResponseParameterEnumValueRemovedCheck
		newBackwardCompatibilityRule(ResponsePropertyEnumValueRemovedId, ERR, false, ResponseParameterEnumValueRemovedCheck), // INFO or ERR
		// ResponseMediaTypeEnumValueRemovedCheck
		newBackwardCompatibilityRule(ResponseMediaTypeEnumValueRemovedId, ERR, false, ResponseMediaTypeEnumValueRemovedCheck), // INFO or ERR
		// RequestBodyEnumValueRemovedCheck
		newBackwardCompatibilityRule(RequestBodyEnumValueRemovedId, ERR, false, RequestBodyEnumValueRemovedCheck), // INFO or ERR
	}
}

func GetOptionalRules() []BackwardCompatibilityRule {

	result := []BackwardCompatibilityRule{}
	for _, rule := range GetAllRules() {
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

func GetOptionalRuleIds() []string {

	result := []string{}
	for _, rule := range GetOptionalRules() {
		result = append(result, rule.Id)
	}
	return result
}

func GetRequiredRules() []BackwardCompatibilityRule {

	result := []BackwardCompatibilityRule{}
	for _, rule := range GetAllRules() {
		if rule.Required {
			result = append(result, rule)
		}
	}
	return result
}
