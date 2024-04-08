package diff

import (
	"github.com/wI2L/jsondiff"
)

// jsonPatch is a wrapper to jsondiff.jsonPatch with proper serialization for json and yaml
type jsonPatch []*jsonOperation

// Empty indicates whether a change was found in this element
func (p jsonPatch) Empty() bool {
	return len(p) == 0
}

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
