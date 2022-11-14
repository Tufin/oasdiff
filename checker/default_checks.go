package checker

func DefaultChecks() []BackwardCompatibilityCheck {
	checks := []BackwardCompatibilityCheck{
		ParameterRemovedCheck,
		NewRequiredRequestPropertyCheck,
		RequestParameterPatternAddedOrChangedCheck,
		RequestPropertyPatternAddedOrChangedCheck,
		AddedRequiredRequestBodyCheck,
		RequestParameterBecameRequiredCheck,
		RequestPropertyBecameRequiredCheck,
		RequestHeaderPropertyBecameRequiredCheck,
		ResponsePropertyBecameOptionalCheck,
		RequestBodyBecameRequiredCheck,
	}
	return checks
}
