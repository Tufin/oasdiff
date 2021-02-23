package diff

import "github.com/getkin/kin-openapi/openapi3"

// CallbacksDiff is a diff between two sets of callbacks
type CallbacksDiff struct {
	Added   StringList `json:"added,omitempty"`
	Deleted StringList `json:"deleted,omitempty"`
	// Modified ModifiedCallbacks `json:"modified,omitempty"`
}

func (callbackDiff *CallbacksDiff) empty() bool {
	return len(callbackDiff.Added) == 0 &&
		len(callbackDiff.Deleted) == 0 //&& len(callbackDiff.Modified) == 0
}

// ModifiedCallbacks is map of callback names to their respective diffs
type ModifiedCallbacks map[string]CallbackDiff

func newCallbacksDiff() *CallbacksDiff {
	return &CallbacksDiff{
		Added:   StringList{},
		Deleted: StringList{},
		// Modified: ModifiedCallbacks{},
	}
}

func getCallbacksDiff(callbacks1, callbacks2 openapi3.Callbacks) *CallbacksDiff {

	result := newCallbacksDiff()

	for callbackValue1, callbackRef1 := range callbacks1 {
		if callbackRef1 != nil && callbackRef1.Value != nil {
			if callbackValue2, ok := callbacks2[callbackValue1]; ok {
				if diff := diffCallbackValues(callbackRef1.Value, callbackValue2.Value); !diff.empty() {
					// result.Modified[callbackValue1] = diff
				}
			} else {
				result.Deleted = append(result.Deleted, callbackValue1)
			}
		}
	}

	for callbackValue2, callbackRef2 := range callbacks2 {
		if callbackRef2 != nil && callbackRef2.Value != nil {
			if _, ok := callbacks1[callbackValue2]; !ok {
				result.Added = append(result.Added, callbackValue2)
			}
		}
	}

	return result

}
