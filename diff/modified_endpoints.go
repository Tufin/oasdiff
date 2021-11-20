package diff

// ModifiedEndpoints is a map of endpoints to their respective diffs
type ModifiedEndpoints map[Endpoint]*MethodDiff

// Breaking indicates whether this element includes a breaking change
func (diff ModifiedEndpoints) Breaking() bool {
	for _, methodDiff := range diff {
		if methodDiff.Breaking() {
			return true
		}
	}

	return false
}

// ToEndpoints returns a list of ModifiedEndpoints keys
func (modifiedEndpoints ModifiedEndpoints) ToEndpoints() Endpoints {
	keys := make(Endpoints, 0, len(modifiedEndpoints))
	for k := range modifiedEndpoints {
		keys = append(keys, k)
	}
	return keys
}
