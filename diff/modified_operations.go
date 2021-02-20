package diff

// ModifiedOperations maps HTTP methods to thier diff
type ModifiedOperations map[string]*MethodDiff

func (modifiedOperations ModifiedOperations) empty() bool {
	return len(modifiedOperations) == 0
}
