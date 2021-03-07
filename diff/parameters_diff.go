package diff

import "github.com/getkin/kin-openapi/openapi3"

// ParametersDiff is a diff between two lists of parameter objects: https://swagger.io/specification/#parameter-object
type ParametersDiff struct {
	Added    ParamNamesByLocation `json:"added,omitempty"`
	Deleted  ParamNamesByLocation `json:"deleted,omitempty"`
	Modified ParamDiffByLocation  `json:"modified,omitempty"`
}

func (parametersDiff *ParametersDiff) empty() bool {
	if parametersDiff == nil {
		return true
	}

	return len(parametersDiff.Added) == 0 &&
		len(parametersDiff.Deleted) == 0 &&
		len(parametersDiff.Modified) == 0
}

// ParamNamesByLocation maps param location (path, query, header or cookie) to the params in this location
type ParamNamesByLocation map[string]ParamNames

// ParamDiffByLocation maps param location (path, query, header or cookie) to param diffs in this location
type ParamDiffByLocation map[string]ParamDiffs

func newParametersDiff() *ParametersDiff {
	return &ParametersDiff{
		Added:    ParamNamesByLocation{},
		Deleted:  ParamNamesByLocation{},
		Modified: ParamDiffByLocation{},
	}
}

// ParamNames is a set of parameter names
type ParamNames map[string]struct{}

// ParamDiffs is map of parameter names to their respective diffs
type ParamDiffs map[string]*ParameterDiff

func (parametersDiff *ParametersDiff) addAddedParam(param *openapi3.Parameter) {

	if paramNames, ok := parametersDiff.Added[param.In]; ok {
		paramNames[param.Name] = struct{}{}
	} else {
		parametersDiff.Added[param.In] = ParamNames{param.Name: struct{}{}}
	}
}

func (parametersDiff *ParametersDiff) addDeletedParam(param *openapi3.Parameter) {

	if paramNames, ok := parametersDiff.Deleted[param.In]; ok {
		paramNames[param.Name] = struct{}{}
	} else {
		parametersDiff.Deleted[param.In] = ParamNames{param.Name: struct{}{}}
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
	if diff.empty() {
		return nil
	}
	return diff
}

func getParametersDiffInternal(config *Config, params1, params2 openapi3.Parameters) *ParametersDiff {

	result := newParametersDiff()

	for _, paramRef1 := range params1 {
		if paramRef1 != nil && paramRef1.Value != nil {
			if paramValue2, ok := findParam(paramRef1.Value, params2); ok {
				if diff := getParameterDiff(config, paramRef1.Value, paramValue2); !diff.empty() {
					result.addModifiedParam(paramRef1.Value, diff)
				}
			} else {
				result.addDeletedParam(paramRef1.Value)
			}
		}
	}

	for _, paramRef2 := range params2 {
		if paramRef2 != nil && paramRef2.Value != nil {
			if _, ok := findParam(paramRef2.Value, params1); !ok {
				result.addAddedParam(paramRef2.Value)
			}
		}
	}

	return result
}

func findParam(param1 *openapi3.Parameter, params2 openapi3.Parameters) (*openapi3.Parameter, bool) {
	// TODO: optimize with a map
	for _, paramRef2 := range params2 {
		if paramRef2 == nil || paramRef2.Value == nil {
			continue
		}

		if equalParams(param1, paramRef2.Value) {
			return paramRef2.Value, true
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
