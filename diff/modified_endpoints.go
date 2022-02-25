package diff

// ModifiedEndpoints is a map of endpoints to their respective diffs
type ModifiedEndpoints map[Endpoint]*MethodDiff

// ToEndpoints returns the modified endpoints
func (modifiedEndpoints ModifiedEndpoints) ToEndpoints() Endpoints {
	keys := make(Endpoints, len(modifiedEndpoints))
	i := 0
	for k := range modifiedEndpoints {
		keys[i] = k
		i++
	}
	return keys
}
