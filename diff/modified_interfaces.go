package diff

import (
	"encoding/json"

	"github.com/wI2L/jsondiff"
)

// ModifiedInterfaces is map of interface names to their respective diffs
type ModifiedInterfaces map[string]jsondiff.Patch

// Empty indicates whether a change was found in this element
func (modifiedInterfaces ModifiedInterfaces) Empty() bool {
	return len(modifiedInterfaces) == 0
}

// GetJsonOrigValue returns the original value of the diff, only if there is exactly one diff
func GetJsonOrigValue(patch jsondiff.Patch) (json.RawMessage, bool) {
	if len(patch) != 1 {
		return nil, false
	}
	result, ok := patch[0].OldValue.(json.RawMessage)
	if !ok {
		return nil, false
	}
	return result, true
}

// GetJsonOrigValue returns the new value of the diff, only if there is exactly one diff
func GetJsonNewValue(patch jsondiff.Patch) (json.RawMessage, bool) {
	if len(patch) != 1 {
		return nil, false
	}
	result, ok := patch[0].Value.(json.RawMessage)
	if !ok {
		return nil, false
	}
	return result, true
}
