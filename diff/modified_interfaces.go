package diff

// ModifiedInterfaces is map of interface names to their respective diffs
type ModifiedInterfaces map[string]JsonPatch

// Empty indicates whether a change was found in this element
func (modifiedInterfaces ModifiedInterfaces) Empty() bool {
	return len(modifiedInterfaces) == 0
}
