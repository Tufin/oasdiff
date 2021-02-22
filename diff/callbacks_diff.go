package diff

// CallbacksDiff is a diff between two sets of callbacks
type CallbacksDiff struct {
	Added    StringList        `json:"added,omitempty"`
	Deleted  StringList        `json:"deleted,omitempty"`
	Modified ModifiedCallbacks `json:"modified,omitempty"`
}

func (callbackDiff *CallbacksDiff) empty() bool {
	return len(callbackDiff.Added) == 0 &&
		len(callbackDiff.Deleted) == 0 &&
		len(callbackDiff.Modified) == 0
}

// ModifiedCallbacks is map of callback names to their respective diffs
type ModifiedCallbacks map[string]CallbackDiff
