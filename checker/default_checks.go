package checker

func DefaultChecks() []BackwardCompatibilityCheck {
	checks := make([]BackwardCompatibilityCheck, 0)
	checks = append(checks, ParameterRemovedCheck)
	checks = append(checks, NewRequiredRequestPropertyCheck)
	checks = append(checks, RequestParameterPatternAddedOrChangedCheck)
	checks = append(checks, RequestPropertyPatternAddedOrChangedCheck)
	checks = append(checks, AddedRequiredBodyCheck)
	return checks
}
