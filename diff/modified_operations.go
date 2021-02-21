package diff

// ModifiedOperations is a map of HTTP methods to their respective diffs
type ModifiedOperations map[string]*MethodDiff

func (modifiedOperations ModifiedOperations) empty() bool {
	return len(modifiedOperations) == 0
}
