package diff

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
)

// ParametersDiff describes the changes between a pair of lists of parameter objects: https://swagger.io/specification/#parameter-object
type ParametersDiff struct {
	Added    ParamNamesByLocation `json:"added,omitempty" yaml:"added,omitempty"`
	Deleted  ParamNamesByLocation `json:"deleted,omitempty" yaml:"deleted,omitempty"`
	Modified ParamDiffByLocation  `json:"modified,omitempty" yaml:"modified,omitempty"`
}

// Empty indicates whether a change was found in this element
func (parametersDiff *ParametersDiff) Empty() bool {
	if parametersDiff == nil {
		return true
	}

	return len(parametersDiff.Added) == 0 &&
		len(parametersDiff.Deleted) == 0 &&
		len(parametersDiff.Modified) == 0
}

// Breaking indicates whether this element includes a breaking change
func (parametersDiff *ParametersDiff) Breaking() bool {
	if parametersDiff.Empty() {
		return false
	}

	return len(parametersDiff.Deleted) > 0 ||
		parametersDiff.Modified.Breaking()
}

// ParamNamesByLocation maps param location (path, query, header or cookie) to the params in this location
type ParamNamesByLocation map[string]StringList

// ParamDiffByLocation maps param location (path, query, header or cookie) to param diffs in this location
type ParamDiffByLocation map[string]ParamDiffs

func (paramDiffByLocation ParamDiffByLocation) Breaking() bool {
	for _, paramDiffs := range paramDiffByLocation {
		if paramDiffs.Breaking() {
			return true
		}
	}
	return false
}

func newParametersDiff() *ParametersDiff {
	return &ParametersDiff{
		Added:    ParamNamesByLocation{},
		Deleted:  ParamNamesByLocation{},
		Modified: ParamDiffByLocation{},
	}
}

// ParamDiffs is map of parameter names to their respective diffs
type ParamDiffs map[string]*ParameterDiff

func (paramDiffs ParamDiffs) Breaking() bool {
	for _, parameterDiff := range paramDiffs {
		if parameterDiff.Breaking() {
			return true
		}
	}
	return false
}

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

func getParametersDiff(config *Config, params1, params2 openapi3.Parameters) (*ParametersDiff, error) {
	diff, err := getParametersDiffInternal(config, params1, params2)
	if err != nil {
		return nil, err
	}

	if diff.Empty() {
		return nil, nil
	}
	return diff, nil
}

func getParametersDiffInternal(config *Config, params1, params2 openapi3.Parameters) (*ParametersDiff, error) {

	result := newParametersDiff()

	for _, paramRef1 := range params1 {
		param1, err := derefParam(paramRef1)
		if err != nil {
			return nil, err
		}

		param2, err := findParam(param1, params2)
		if err != nil {
			return nil, err
		}

		if param2 != nil {
			diff, err := getParameterDiff(config, param1, param2)
			if err != nil {
				return nil, err
			}

			if !diff.Empty() {
				result.addModifiedParam(param1, diff)
			}
		} else {
			result.addDeletedParam(param1)
		}
	}

	for _, paramRef2 := range params2 {
		param2, err := derefParam(paramRef2)
		if err != nil {
			return nil, err
		}

		param, err := findParam(param2, params1)
		if err != nil {
			return nil, err
		}
		if param == nil {
			result.addAddedParam(param2)
		}
	}

	return result, nil
}

func derefParam(ref *openapi3.ParameterRef) (*openapi3.Parameter, error) {

	if ref == nil || ref.Value == nil {
		return nil, fmt.Errorf("parameter reference is nil")
	}

	return ref.Value, nil
}

func findParam(param1 *openapi3.Parameter, params2 openapi3.Parameters) (*openapi3.Parameter, error) {
	// TODO: optimize with a map
	for _, paramRef2 := range params2 {
		param2, err := derefParam(paramRef2)
		if err != nil {
			return nil, err
		}

		equal, err := equalParams(param1, param2)
		if err != nil {
			return nil, err
		}

		if equal {
			return param2, nil
		}
	}

	return nil, nil
}

func equalParams(param1 *openapi3.Parameter, param2 *openapi3.Parameter) (bool, error) {
	if param1 == nil || param2 == nil {
		return false, fmt.Errorf("param is nil")
	}
	return param1.Name == param2.Name && param1.In == param2.In, nil
}

func (parametersDiff *ParametersDiff) getSummary() *SummaryDetails {
	return &SummaryDetails{
		Added:    len(parametersDiff.Added),
		Deleted:  len(parametersDiff.Deleted),
		Modified: len(parametersDiff.Modified),
	}
}

// Patch applies the patch to parameters
func (parametersDiff *ParametersDiff) Patch(parameters openapi3.Parameters) error {

	if parametersDiff.Empty() {
		return nil
	}

	for location, paramDiffs := range parametersDiff.Modified {
		for name, parameterDiff := range paramDiffs {
			err := parameterDiff.Patch(parameters.GetByInAndName(location, name))
			if err != nil {
				return err
			}
		}
	}

	return nil
}
