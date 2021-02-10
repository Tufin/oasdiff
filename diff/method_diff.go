package diff

import (
	"github.com/getkin/kin-openapi/openapi3"
)

type MethodDiff struct {
	DeletedParams  map[string]ParamNames `json:"deletedParams,omitempty"`  // key is param location
	ModifiedParams map[string]ParamDiffs `json:"modifiedParams,omitempty"` // key is param location
}

type ParamNames map[string]struct{}  // key is param name
type ParamDiffs map[string]ParamDiff // key is param name

func newMethodDiff() *MethodDiff {
	return &MethodDiff{
		DeletedParams:  map[string]ParamNames{},
		ModifiedParams: map[string]ParamDiffs{},
	}
}

func (methodDiff *MethodDiff) empty() bool {
	return len(methodDiff.DeletedParams) == 0 && len(methodDiff.ModifiedParams) == 0
}

func (methodDiff *MethodDiff) addDeletedParam(param *openapi3.Parameter) {

	if paramNames, ok := methodDiff.DeletedParams[param.In]; ok {
		paramNames[param.Name] = struct{}{}
	} else {
		methodDiff.DeletedParams[param.In] = ParamNames{param.Name: struct{}{}}
	}
}

func (methodDiff *MethodDiff) addModifiedParam(param *openapi3.Parameter, diff ParamDiff) {

	if paramDiffs, ok := methodDiff.ModifiedParams[param.In]; ok {
		paramDiffs[param.Name] = diff
	} else {
		methodDiff.ModifiedParams[param.In] = ParamDiffs{param.Name: diff}
	}
}

func diffParameters(params1 openapi3.Parameters, params2 openapi3.Parameters) *MethodDiff {

	result := newMethodDiff()

	for _, paramRef1 := range params1 {
		if paramRef1 != nil && paramRef1.Value != nil {
			if paramValue2, ok := findParam(paramRef1.Value, params2); ok {
				if diff := diffParamValues(paramRef1.Value, paramValue2); !diff.empty() {
					result.addModifiedParam(paramRef1.Value, diff)
				}
			} else {
				result.addDeletedParam(paramRef1.Value)
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
