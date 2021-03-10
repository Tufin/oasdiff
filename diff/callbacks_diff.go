package diff

import "github.com/getkin/kin-openapi/openapi3"

// CallbacksDiff is a diff between callback objects: https://swagger.io/specification/#callback-object
type CallbacksDiff struct {
	Added    StringList        `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  StringList        `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ModifiedCallbacks `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty return true if there is no diff
func (callbacksDiff *CallbacksDiff) Empty() bool {

	if callbacksDiff == nil {
		return true
	}

	return len(callbacksDiff.Added) == 0 &&
		len(callbacksDiff.Deleted) == 0 &&
		len(callbacksDiff.Modified) == 0
}

// ModifiedCallbacks is map of callback names to their respective diffs
type ModifiedCallbacks map[string]*PathsDiff

func newCallbacksDiff() *CallbacksDiff {
	return &CallbacksDiff{
		Added:    StringList{},
		Deleted:  StringList{},
		Modified: ModifiedCallbacks{},
	}
}

func getCallbacksDiff(config *Config, callbacks1, callbacks2 openapi3.Callbacks) *CallbacksDiff {
	diff := getCallbacksDiffInternal(config, callbacks1, callbacks2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getCallbacksDiffInternal(config *Config, callbacks1, callbacks2 openapi3.Callbacks) *CallbacksDiff {

	result := newCallbacksDiff()

	for callbackValue1, callbackRef1 := range callbacks1 {
		if callbackRef1 != nil && callbackRef1.Value != nil {
			if callbackValue2, ok := callbacks2[callbackValue1]; ok {
				if diff := diffCallbackValues(config, callbackRef1.Value, callbackValue2.Value); !diff.Empty() {
					result.Modified[callbackValue1] = diff
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

func diffCallbackValues(config *Config, callback1, callback2 *openapi3.Callback) *PathsDiff {
	return getPathsDiff(config, openapi3.Paths(*callback1), openapi3.Paths(*callback2))
}

func (callbacksDiff *CallbacksDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(callbacksDiff.Added),
		Deleted:  len(callbacksDiff.Deleted),
		Modified: len(callbacksDiff.Modified),
	}
}
