package diff

// ModifiedEndpoints is a map of endpoints to their respective diffs
type ModifiedEndpoints map[Endpoint]*MethodDiff

// ToEndpoints returns a list of ModifiedEndpoints keys
func (modifiedEndpoints ModifiedEndpoints) ToEndpoints() Endpoints {
	keys := make(Endpoints, 0, len(modifiedEndpoints))
	for k := range modifiedEndpoints {
		keys = append(keys, k)
	}
	return keys
}
