package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

// ParametersDiff is a diff between two lists of parameter objects: https://swagger.io/specification/#parameter-object
type ParametersDiff struct {
	Added    ParamNamesByLocation `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  ParamNamesByLocation `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ParamDiffByLocation  `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty return true if there is no diff
func (parametersDiff *ParametersDiff) Empty() bool {
	if parametersDiff == nil {
		return true
	}

	return len(parametersDiff.Added) == 0 &&
		len(parametersDiff.Deleted) == 0 &&
		len(parametersDiff.Modified) == 0
}

// ParamNamesByLocation maps param location (path, query, header or cookie) to the params in this location
type ParamNamesByLocation map[string]StringList

// ParamDiffByLocation maps param location (path, query, header or cookie) to param diffs in this location
type ParamDiffByLocation map[string]ParamDiffs

func newParametersDiff() *ParametersDiff {
	return &ParametersDiff{
		Added:    ParamNamesByLocation{},
		Deleted:  ParamNamesByLocation{},
		Modified: ParamDiffByLocation{},
	}
}

// ParamDiffs is map of parameter names to their respective diffs
type ParamDiffs map[string]*ParameterDiff

func (parametersDiff *ParametersDiff) addAddedParam(param *openapi3.Parameter) {

	if paramNames, ok := parametersDiff.Added[param.In]; ok {
		parametersDiff.Added[param.In] = append(paramNames, param.Name)
	} else {
		parametersDiff.Added[param.In] = StringList{param.Name}
	}
}

func (parametersDiff *ParametersDiff) addDeletedParam(param *openapi3.Parameter) {

	if paramNames, ok := parametersDiff.Deleted[param.In]; ok {
		parametersDiff.Deleted[param.In] = append(paramNames, param.Name)
	} else {
		parametersDiff.Deleted[param.In] = StringList{param.Name}
	}
}

func (parametersDiff *ParametersDiff) addModifiedParam(param *openapi3.Parameter, diff *ParameterDiff) {

	if paramDiffs, ok := parametersDiff.Modified[param.In]; ok {
		paramDiffs[param.Name] = diff
	} else {
		parametersDiff.Modified[param.In] = ParamDiffs{param.Name: diff}
	}
}

func getParametersDiff(config *Config, params1, params2 openapi3.Parameters) *ParametersDiff {
	diff := getParametersDiffInternal(config, params1, params2)
	if diff.Empty() {
		return nil
	}
	return diff
}

func getParametersDiffInternal(config *Config, params1, params2 openapi3.Parameters) *ParametersDiff {

	result := newParametersDiff()

	for _, paramRef1 := range params1 {
		value1 := deref(paramRef1)

		if paramValue2, ok := findParam(value1, params2); ok {
			if diff := getParameterDiff(config, value1, paramValue2); !diff.Empty() {
				result.addModifiedParam(value1, diff)
			}
		} else {
			result.addDeletedParam(value1)
		}
	}

	for _, paramRef2 := range params2 {
		value2 := deref(paramRef2)

		if _, ok := findParam(value2, params1); !ok {
			result.addAddedParam(value2)
		}
	}

	return result
}

func deref(ref *openapi3.ParameterRef) *openapi3.Parameter {
	// TODO: check if ref == nil
	return ref.Value
}

func findParam(param1 *openapi3.Parameter, params2 openapi3.Parameters) (*openapi3.Parameter, bool) {
	// TODO: optimize with a map
	for _, paramRef2 := range params2 {
		value2 := deref(paramRef2)

		if equalParams(param1, value2) {
			return value2, true
		}
	}

	return nil, false
}

func equalParams(param1 *openapi3.Parameter, param2 *openapi3.Parameter) bool {
	return param1.Name == param2.Name && param1.In == param2.In
}

func (parametersDiff *ParametersDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(parametersDiff.Added),
		Deleted:  len(parametersDiff.Deleted),
		Modified: len(parametersDiff.Modified),
	}
}
