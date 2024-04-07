package diff

import (
	"encoding/json"

	"github.com/wI2L/jsondiff"
)

// jsonPatch is a wrapper to jsondiff.jsonPatch with proper serialization for json and yaml
type jsonPatch []*jsonOperation

// jsonOperation is a wrapper to jsondiff.jsonOperation with proper serialization for json and yaml
type jsonOperation struct {
	OldValue interface{} `json:"oldValue" yaml:"oldValue"`
	Value    interface{} `json:"value" yaml:"value"`
	Type     string      `json:"op" yaml:"op"`
	From     string      `json:"from" yaml:"from"`
	Path     string      `json:"path" yaml:"path"`
}

func toJsonPatch(patch jsondiff.Patch) jsonPatch {
	result := make(jsonPatch, len(patch))
	for i, op := range patch {
		result[i] = newJsonOperation(op)
	}
	return result
}

func newJsonOperation(op jsondiff.Operation) *jsonOperation {
	return &jsonOperation{
		OldValue: op.OldValue,
		Value:    op.Value,
		Type:     op.Type,
		From:     op.From,
		Path:     op.Path,
	}
}

func compareJson(source, target interface{}, opts ...jsondiff.Option) (jsonPatch, error) {
	patch, err := jsondiff.Compare(source, target, opts...)
	if err != nil {
		return nil, err
	}
	return toJsonPatch(patch), nil
}

// GetJsonOrigValue returns the original value of the diff, only if there is exactly one diff
func GetJsonOrigValue(patch jsonPatch) (json.RawMessage, bool) {
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
func GetJsonNewValue(patch jsonPatch) (json.RawMessage, bool) {
	if len(patch) != 1 {
		return nil, false
	}
	result, ok := patch[0].Value.(json.RawMessage)
	if !ok {
		return nil, false
	}
	return result, true
}
