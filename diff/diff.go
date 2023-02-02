package diff

import (
	"encoding/json"
	"errors"
	"fmt"

	"cloud.google.com/go/civil"
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/load"
)

const SinceDateExtension = "x-since-date"

var (
	DefaultSinceDate = civil.Date{Year: 2000, Month: 1, Day: 1}
)

// Diff describes the changes between a pair of OpenAPI objects: https://swagger.io/specification/#schema
type Diff struct {
	ExtensionsDiff   *ExtensionsDiff           `json:"extensions,omitempty" yaml:"extensions,omitempty"`
	OpenAPIDiff      *ValueDiff                `json:"openAPI,omitempty" yaml:"openAPI,omitempty"`
	InfoDiff         *InfoDiff                 `json:"info,omitempty" yaml:"info,omitempty"`
	PathsDiff        *PathsDiff                `json:"paths,omitempty" yaml:"paths,omitempty"`
	EndpointsDiff    *EndpointsDiff            `json:"endpoints,omitempty" yaml:"endpoints,omitempty"`
	SecurityDiff     *SecurityRequirementsDiff `json:"security,omitempty" yaml:"security,omitempty"`
	ServersDiff      *ServersDiff              `json:"servers,omitempty" yaml:"servers,omitempty"`
	TagsDiff         *TagsDiff                 `json:"tags,omitempty" yaml:"tags,omitempty"`
	ExternalDocsDiff *ExternalDocsDiff         `json:"externalDocs,omitempty" yaml:"externalDocs,omitempty"`

	ComponentsDiff `json:"components,omitempty" yaml:"components,omitempty"`
}

type OperationsSourcesMap map[*openapi3.Operation]string

func newDiff() *Diff {
	return &Diff{}
}

// Empty indicates whether a change was found in this element
func (diff *Diff) Empty() bool {
	return diff == nil || *diff == Diff{}
}

func (diff *Diff) removeNonBreaking() {

	if diff.Empty() {
		return
	}

	diff.ExtensionsDiff = nil
	diff.OpenAPIDiff = nil
	diff.InfoDiff = nil
	diff.ServersDiff = nil
	diff.TagsDiff = nil
	diff.ExternalDocsDiff = nil

	diff.ComponentsDiff.removeNonBreaking()
}

/*
Get calculates the diff between a pair of OpenAPI objects.

Note that Get expects OpenAPI References (https://swagger.io/docs/specification/using-ref/) to be resolved.
References are normally resolved automatically when you load the spec.
In other cases you can resolve refs using https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn.
*/
func Get(config *Config, s1, s2 *openapi3.T) (*Diff, error) {
	diff, err := getDiff(config, newState(), s1, s2)
	if err != nil {
		return nil, err
	}
	return diff, nil
}

/*
GetWithOperationsSourcesMap calculates the diff between a pair of OpenAPI objects.

Note that GetWithOperationsSourcesMap expects OpenAPI References (https://swagger.io/docs/specification/using-ref/) to be resolved.
References are normally resolved automatically when you load the spec.
In other cases you can resolve refs using https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn.
*/
func GetWithOperationsSourcesMap(config *Config, s1, s2 *load.OpenAPISpecInfo) (*Diff, *OperationsSourcesMap, error) {
	diff, err := getDiff(config, newState(), s1.Spec, s2.Spec)
	if err != nil {
		return nil, nil, err
	}

	_, operationsSources1, err := mergedPaths([]load.OpenAPISpecInfo{*s1})
	if err != nil {
		return nil, nil, err
	}
	_, operationsSources2, err := mergedPaths([]load.OpenAPISpecInfo{*s2})
	if err != nil {
		return nil, nil, err
	}

	operationsSources := *operationsSources1
	for k, v := range *operationsSources2 {
		operationsSources[k] = v
	}
	return diff, &operationsSources, nil
}

/*
GetPathsDiff calculates the diff between a pair of slice of OpenAPI objects.
It is helpfull when you want to find diff and check for breaking changes for API divided into multiple files.
If there are same paths in different OpenAPI objects, then function uses version of the path with the last x-since-date extension.
The x-since-date extension should be set on path or operations level. Extension set on the operations level overrides the value set on path level.
If such path doesn't have the x-since-date extension, its value is default "2000-01-01"
If there are same paths with the same x-since-date value, then function returns error.
The format of the x-since-date is the RFC3339 full-date format

Note that Get expects OpenAPI References (https://swagger.io/docs/specification/using-ref/) to be resolved.
References are normally resolved automatically when you load the spec.
In other cases you can resolve refs using https://pkg.go.dev/github.com/getkin/kin-openapi/openapi3#Loader.ResolveRefsIn.
*/
func GetPathsDiff(config *Config, s1, s2 []load.OpenAPISpecInfo) (*Diff, *OperationsSourcesMap, error) {
	state := newState()
	result := newDiff()
	var err error
	paths1, operationsSources1, err := mergedPaths(s1)
	if err != nil {
		return nil, nil, err
	}
	paths2, operationsSources2, err := mergedPaths(s2)
	if err != nil {
		return nil, nil, err
	}

	if result.PathsDiff, err = getPathsDiff(config, state, *paths1, *paths2); err != nil {
		return nil, nil, err
	}

	if result.EndpointsDiff, err = getEndpointsDiff(config, state, *paths1, *paths2); err != nil {
		return nil, nil, err
	}

	operationsSources := *operationsSources1
	for k, v := range *operationsSources2 {
		operationsSources[k] = v
	}
	return result, &operationsSources, nil
}

func mergedPaths(s1 []load.OpenAPISpecInfo) (*openapi3.Paths, *OperationsSourcesMap, error) {
	result := make(openapi3.Paths, 0)
	operationsSources := make(OperationsSourcesMap)
	for _, s := range s1 {
		for path, pathItem := range s.Spec.Paths {
			p := result.Find(path)
			if p == nil {
				result[path] = pathItem
				for _, opItem := range pathItem.Operations() {
					operationsSources[opItem] = s.Url
				}
				continue
			}

			for op, opItem := range pathItem.Operations() {
				oldOperation := p.GetOperation(op)
				if oldOperation == nil {
					p.SetOperation(op, opItem)
					operationsSources[opItem] = s.Url
					continue
				}

				oldSince, err := sinceDateFrom(*p, *oldOperation)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid %s extension value in %s(%s %s), %w", SinceDateExtension, operationsSources[oldOperation], op, path, err)
				}
				newSince, err := sinceDateFrom(*pathItem, *opItem)
				if err != nil {
					return nil, nil, fmt.Errorf("invalid %s extension value in %s(%s %s), %w", SinceDateExtension, s.Url, op, path, err)
				}
				if newSince.After(oldSince) {
					p.SetOperation(op, opItem)
					operationsSources[opItem] = s.Url
				}

				if newSince == oldSince {
					return nil, nil, fmt.Errorf("multiple endpoints found in %s(%s %s) and %s(%s %s). Add the %s extension with ordered values to operations to specify its order", operationsSources[oldOperation], op, path, s.Url, op, path, SinceDateExtension)
				}
			}

		}
	}
	return &result, &operationsSources, nil
}

func sinceDateFrom(pathItem openapi3.PathItem, operation openapi3.Operation) (civil.Date, error) {
	since, _, err := getSinceDate(pathItem.Extensions)
	if err != nil {
		return DefaultSinceDate, err
	}
	opSince, ok, err := getSinceDate(operation.Extensions)
	if err != nil {
		return DefaultSinceDate, err
	}
	if ok {
		since = opSince
	}
	return since, nil
}

func getSinceDate(extensions map[string]interface{}) (civil.Date, bool, error) {
	var since string
	since, ok := extensions[SinceDateExtension].(string)
	if !ok {
		sinceJson, ok := extensions[SinceDateExtension].(json.RawMessage)

		if !ok {
			return DefaultSinceDate, false, nil
		}

		if err := json.Unmarshal(sinceJson, &since); err != nil {
			return civil.Date{}, false, errors.New("unmarshal failed")
		}
	}

	date, err := civil.ParseDate(since)
	if err != nil {
		return civil.Date{}, false, errors.New("failed to parse time")
	}

	return date, true, nil
}

func getDiff(config *Config, state *state, s1, s2 *openapi3.T) (*Diff, error) {

	if s1 == nil || s2 == nil {
		return nil, errors.New("spec is nil")
	}

	diff, err := getDiffInternal(config, state, s1, s2)
	if err != nil {
		return nil, err
	}

	if config.BreakingOnly {
		diff.removeNonBreaking()
	}

	if diff.Empty() {
		return nil, nil
	}

	return diff, nil
}

func getDiffInternal(config *Config, state *state, s1, s2 *openapi3.T) (*Diff, error) {

	result := newDiff()
	var err error

	result.ExtensionsDiff = getExtensionsDiff(config, state, s1.Extensions, s2.Extensions)
	result.OpenAPIDiff = getValueDiff(s1.OpenAPI, s2.OpenAPI)

	if result.InfoDiff, err = getInfoDiff(config, state, s1.Info, s2.Info); err != nil {
		return nil, err
	}

	if result.PathsDiff, err = getPathsDiff(config, state, s1.Paths, s2.Paths); err != nil {
		return nil, err
	}

	if result.EndpointsDiff, err = getEndpointsDiff(config, state, s1.Paths, s2.Paths); err != nil {
		return nil, err
	}

	result.SecurityDiff = getSecurityRequirementsDiff(config, state, &s1.Security, &s2.Security)
	result.ServersDiff = getServersDiff(config, state, &s1.Servers, &s2.Servers)
	result.TagsDiff = getTagsDiff(config, state, s1.Tags, s2.Tags)
	result.ExternalDocsDiff = getExternalDocsDiff(config, state, s1.ExternalDocs, s2.ExternalDocs)

	if result.ComponentsDiff, err = getComponentsDiff(config, state, s1.Components, s2.Components); err != nil {
		return nil, err
	}

	return result, nil
}

// GetSummary returns a summary of the changes in the diff
func (diff *Diff) GetSummary() *Summary {

	summary := newSummary()

	if diff.Empty() {
		return summary
	}

	summary.Diff = true

	// swagger
	summary.add(diff.PathsDiff, PathsDetail)
	summary.add(diff.SecurityDiff, SecurityDetail)
	summary.add(diff.ServersDiff, ServersDetail)
	summary.add(diff.TagsDiff, TagsDetail)

	// components
	summary.add(diff.SchemasDiff, SchemasDetail)
	summary.add(diff.ParametersDiff, ParametersDetail)
	summary.add(diff.HeadersDiff, HeadersDetail)
	summary.add(diff.RequestBodiesDiff, RequestBodiesDetail)
	summary.add(diff.ResponsesDiff, ResponsesDetail)
	summary.add(diff.SecuritySchemesDiff, SecuritySchemesDetail)
	summary.add(diff.ExamplesDiff, ExamplesDetail)
	summary.add(diff.LinksDiff, LinksDetail)
	summary.add(diff.CallbacksDiff, CallbacksDetail)

	// special
	summary.add(diff.EndpointsDiff, EndpointsDetail)
	return summary
}

// Patch applies the patch to a spec
func (diff *Diff) Patch(s *openapi3.T) error {

	if diff.Empty() {
		return nil
	}

	err := diff.PathsDiff.Patch(s.Paths)
	if err != nil {
		return err
	}

	return nil
}
