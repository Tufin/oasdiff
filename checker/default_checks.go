package checker

func DefaultChecks() []BackwardCompatibilityCheck {
	checks := []BackwardCompatibilityCheck{
		RequestParameterRemovedCheck,
		NewRequiredRequestPropertyCheck,
		RequestParameterPatternAddedOrChangedCheck,
		RequestPropertyPatternAddedOrChangedCheck,
		AddedRequiredRequestBodyCheck,
		RequestParameterBecameRequiredCheck,
		RequestPropertyBecameRequiredCheck,
		RequestHeaderPropertyBecameRequiredCheck,
		ResponsePropertyBecameOptionalCheck,
		RequestBodyBecameRequiredCheck,
		ResponseHeaderBecameOptional,
		ResponseHeaderRemoved,
		ResponseSuccessStatusRemoved,
		ResponseMediaTypeRemoved,
		NewRequestPathParameterCheck,
		NewRequiredRequestNonPathParameterCheck,
		NewRequiredRequestHeaderPropertyCheck,
	}
	return checks
}
