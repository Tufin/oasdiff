package checker

func DefaultChecks() []BackwardCompatibilityCheck {
	checks := make([]BackwardCompatibilityCheck, 0)
	checks = append(checks, ParameterRemovedCheck)
	checks = append(checks, NewRequiredRequestPropertyCheck)
	return checks
}

