package diff

type ModifiedOperations map[string]*MethodDiff // key is HTTP method

func (modifiedOperations ModifiedOperations) empty() bool {
	return len(modifiedOperations) == 0
}
