package flatten

import (
	"errors"
	"fmt"
	"math"
	"regexp"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/utils"
)

const (
	FormatErrorMessage  = "unable to resolve Format conflict using default resolver: all Format values must be identical"
	TypeErrorMessage    = "unable to resolve Type conflict: all Type values must be identical"
	FormatResolverError = ""

	FormatInt32  = "int32"
	FormatInt64  = "int64"
	FormatFloat  = "float"
	FormatDouble = "double"
)

type SchemaCollection struct {
	Not                  []*openapi3.SchemaRef
	OneOf                []openapi3.SchemaRefs
	AnyOf                []openapi3.SchemaRefs
	Title                []string
	Type                 []string
	Format               []string
	Description          []string
	Enum                 [][]interface{}
	UniqueItems          []bool
	ExclusiveMin         []bool
	ExclusiveMax         []bool
	Min                  []*float64
	Max                  []*float64
	MultipleOf           []*float64
	MinLength            []uint64
	MaxLength            []*uint64
	Pattern              []string
	MinItems             []uint64
	MaxItems             []*uint64
	Items                []*openapi3.SchemaRef
	Required             [][]string
	Properties           []openapi3.Schemas
	MinProps             []uint64
	MaxProps             []*uint64
	AdditionalProperties []openapi3.AdditionalProperties
	Nullable             []bool
	ReadOnly             []bool
	WriteOnly            []bool
}

type state struct {
	visitedSchemas utils.VisitedRefs
}

func NewState() *state {
	return &state{
		visitedSchemas: utils.VisitedRefs{},
	}
}

// Merge replaces objects under AllOf with a flattened equivalent
func Merge(state *state, baseSchemaRef openapi3.SchemaRef) (*openapi3.Schema, error) {

	/*	if state.visitedSchemas.IsVisited(baseSchemaRef.Ref) {
			return baseSchemaRef.Value, nil // TODO: check
		}
	*/
	baseSchema := baseSchemaRef.Value
	/*
		if baseSchemaRef.Ref != "" {
			state.visitedSchemas.Add(baseSchemaRef.Ref)
			defer state.visitedSchemas.Remove(baseSchemaRef.Ref)
		}
	*/
	allOfSchemas, err := getAllOfSchemas(state, *baseSchema)
	if err != nil {
		return &openapi3.Schema{}, err
	}
	schemas := []*openapi3.Schema{baseSchema}
	schemas = append(schemas, allOfSchemas...)
	result, err := flattenSchemas(state, schemas)
	if err != nil {
		return &openapi3.Schema{}, err
	}
	return result, nil
}

func getAllOfSchemas(state *state, schema openapi3.Schema) ([]*openapi3.Schema, error) {
	schemas := []*openapi3.Schema{}
	if schema.AllOf == nil {
		return schemas, nil
	}
	for _, schema := range schema.AllOf {
		merged, err := Merge(state, *schema)
		if err != nil {
			return schemas, err
		}
		schemas = append(schemas, merged)
	}
	return schemas, nil
}

func flattenSchemas(state *state, schemas []*openapi3.Schema) (*openapi3.Schema, error) {
	result := openapi3.NewSchema()
	collection := collect(schemas)

	result.Title = collection.Title[0]
	result.Description = collection.Description[0]
	result, err := resolveFormat(result, &collection)
	if err != nil {
		return result, err
	}
	result, err = resolveType(result, &collection)
	if err != nil {
		return result, err
	}
	result = resolveNumberRange(result, &collection)
	result.MinLength = findMaxValue(collection.MinLength)
	result.MaxLength = findMinValue(collection.MaxLength)
	result.MinItems = findMaxValue(collection.MinItems)
	result.MaxItems = findMinValue(collection.MaxItems)
	result.MinProps = findMaxValue(collection.MinProps)
	result.MaxProps = findMinValue(collection.MaxProps)
	result.Pattern = resolvePattern(collection.Pattern)
	result.Nullable = !hasFalse(collection.Nullable)
	result.ReadOnly = hasTrue(collection.ReadOnly)
	result.WriteOnly = hasTrue(collection.WriteOnly)

	enums, err := resolveEnum(collection.Enum)
	if err != nil {
		return result, err
	}
	result.Enum = enums
	result = resolveMultipleOf(result, &collection)
	result.Required = resolveRequired(collection.Required)
	result, err = resolveItems(state, result, &collection)
	if err != nil {
		return result, err
	}
	result.UniqueItems = resolveUniqueItems(collection.UniqueItems)
	result, err = resolveProperties(state, result, &collection)
	if err != nil {
		return result, err
	}

	result, err = resolveOneOf(state, result, &collection)
	if err != nil {
		return result, err
	}

	result, err = resolveAnyOf(state, result, &collection)
	if err != nil {
		return result, err
	}

	result, err = resolveNot(state, result, &collection)
	if err != nil {
		return result, err
	}

	return result, nil
}

func hasTrue(values []bool) bool {
	for _, val := range values {
		if val {
			return true
		}
	}
	return false
}

func hasFalse(values []bool) bool {
	for _, val := range values {
		if !val {
			return true
		}
	}
	return false
}

func resolveNumberRange(schema *openapi3.Schema, collection *SchemaCollection) *openapi3.Schema {

	//resolve minimum
	max := math.Inf(-1)
	isExcluded := false
	var value *float64
	for i, s := range collection.Min {
		if s != nil {
			if *s > max {
				max = *s
				value = s
				isExcluded = collection.ExclusiveMin[i]
			}
		}
	}

	schema.Min = value
	schema.ExclusiveMin = isExcluded
	//resolve maximum
	min := math.Inf(1)
	isExcluded = false
	// var value *float64
	for i, s := range collection.Max {
		if s != nil {
			if *s < min {
				min = *s
				value = s
				isExcluded = collection.ExclusiveMax[i]
			}
		}
	}

	schema.Max = value
	schema.ExclusiveMax = isExcluded
	return schema
}

func resolveItems(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	items := []*openapi3.Schema{}
	for _, s := range collection.Items {
		if s != nil {
			items = append(items, s.Value)
		}
	}
	if len(items) == 0 {
		schema.Items = nil
		return schema, nil
	}

	res, err := flattenSchemas(state, items)
	if err != nil {
		return schema, err
	}
	ref := openapi3.SchemaRef{
		Value: res,
	}
	schema.Items = &ref
	return schema, nil
}

func resolveUniqueItems(values []bool) bool {
	for _, v := range values {
		if v {
			return true
		}
	}
	return false
}

/* MultipleOf */
func gcd(a, b uint64) uint64 {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b uint64) uint64 {
	return a * b / gcd(a, b)
}

func containsNonInteger(arr []float64) bool {
	for _, num := range arr {
		if num != math.Trunc(num) {
			return true
		}
	}
	return false
}

func resolveMultipleOf(schema *openapi3.Schema, collection *SchemaCollection) *openapi3.Schema {
	values := []float64{}
	for _, v := range collection.MultipleOf {
		if v == nil {
			continue
		}
		values = append(values, *v)
	}
	if len(values) == 0 {
		schema.MultipleOf = nil
		return schema
	}

	factor := 1.0
	for containsNonInteger(values) {
		factor *= 10.0
		for i := range values {
			values[i] *= factor
		}
	}
	uintValues := make([]uint64, len(values))
	for i, val := range values {
		uintValues[i] = uint64(val)
	}
	lcmValue := uintValues[0]
	for _, v := range uintValues {
		lcmValue = lcm(lcmValue, v)
	}
	schema.MultipleOf = openapi3.Float64Ptr(float64(lcmValue) / factor)
	return schema
}

func hasFalseValue(ap []openapi3.AdditionalProperties) bool {
	for _, v := range ap {
		if v.Has != nil && !*v.Has {
			return true
		}
	}
	return false
}

// resolve properties which have additionalProperties that are set to false.
func resolveFalseProps(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	schema = resolveFalseAdditionalProps(schema, collection)
	propsToMerge := getFalsePropsKeys(collection)
	return mergeProps(state, schema, collection, propsToMerge)
}

func resolveFalseAdditionalProps(schema *openapi3.Schema, collection *SchemaCollection) *openapi3.Schema {
	has := false
	schema.AdditionalProperties.Has = &has
	return schema
}

// if there are additionalProperties which are Schemas, they are merged to a single Schema.
func resolveNonFalseAdditionalProps(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	additionalSchemas := []*openapi3.Schema{}
	for _, ap := range collection.AdditionalProperties {
		if ap.Schema != nil && ap.Schema.Value != nil {
			additionalSchemas = append(additionalSchemas, ap.Schema.Value)
		}
	}

	var schemaRef *openapi3.SchemaRef
	if len(additionalSchemas) > 0 {
		result, err := flattenSchemas(state, additionalSchemas)
		if err != nil {
			return schema, err
		}
		schemaRef = &openapi3.SchemaRef{
			Value: result,
		}
	}
	schema.AdditionalProperties.Has = nil
	schema.AdditionalProperties.Schema = schemaRef
	return schema, nil
}

func resolveNonFalseProps(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	result, err := resolveNonFalseAdditionalProps(state, schema, collection)
	if err != nil {
		return &openapi3.Schema{}, err
	}
	propsToMerge := getNonFalsePropsKeys(collection)
	return mergeProps(state, result, collection, propsToMerge)
}

// the output is the intersection of all properties keys of schemas which have additionalProperties set to false.
func getFalsePropsKeys(collection *SchemaCollection) []string {
	properties := [][]string{}
	for i, schema := range collection.Properties {
		additionalProps := collection.AdditionalProperties[i].Has
		if additionalProps != nil && !*additionalProps {
			keys := []string{}
			for key := range schema {
				keys = append(keys, key)
			}
			properties = append(properties, keys)
		}
	}
	return findIntersection(properties...)
}

// the output is a list of unique properties keys of all schemas.
func getNonFalsePropsKeys(collection *SchemaCollection) []string {
	keys := []string{}
	for _, schema := range collection.Properties {
		for key := range schema {
			keys = append(keys, key)
		}
	}
	return getUniqueStrings(keys)
}

func getUniqueStrings(input []string) []string {
	uniqueStrings := make(map[string]bool)

	for _, str := range input {
		uniqueStrings[str] = true
	}

	result := []string{}
	for str := range uniqueStrings {
		result = append(result, str)
	}

	return result
}

func resolveProperties(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	if hasFalseValue(collection.AdditionalProperties) {
		return resolveFalseProps(state, schema, collection)
	} else {
		return resolveNonFalseProps(state, schema, collection)
	}
}

func mergeProps(state *state, schema *openapi3.Schema, collection *SchemaCollection, propsToMerge []string) (*openapi3.Schema, error) {
	propsToSchemasMap := map[string][]*openapi3.Schema{}
	for _, schema := range collection.Properties {
		for propKey, schemaRef := range schema {
			if containsString(propsToMerge, propKey) {
				propMergedSchema, err := Merge(state, *schemaRef)
				if err != nil {
					return &openapi3.Schema{}, err
				}
				propsToSchemasMap[propKey] = append(propsToSchemasMap[propKey], propMergedSchema)
			}
		}
	}

	result := make(openapi3.Schemas)
	for prop, schemas := range propsToSchemasMap {
		mergedProp, err := flattenSchemas(state, schemas)
		if err != nil {
			return &openapi3.Schema{}, err
		}
		ref := openapi3.SchemaRef{
			Value: mergedProp,
		}
		result[prop] = &ref
	}

	if len(result) == 0 {
		result = nil
	}

	schema.Properties = result
	return schema, nil
}

func resolveEnum(values [][]interface{}) ([]interface{}, error) {
	var nonEmptyEnum [][]interface{}
	for _, enum := range values {
		if len(enum) > 0 {
			nonEmptyEnum = append(nonEmptyEnum, enum)
		}
	}
	var intersection []interface{}
	if len(nonEmptyEnum) == 0 {
		return intersection, nil
	}
	intersection = findIntersectionOfArrays(nonEmptyEnum)
	if len(intersection) == 0 {
		return intersection, errors.New("unable to resolve Enum conflict: intersection of values must be non-empty")
	}
	return intersection, nil
}

func resolvePattern(values []string) string {
	var pattern strings.Builder
	for _, p := range values {
		if len(p) > 0 {
			if !isPatternResolved(p) {
				pattern.WriteString(fmt.Sprintf("(?=%s)", p))
			} else {
				pattern.WriteString(p)
			}
		}
	}
	return pattern.String()
}

func isPatternResolved(pattern string) bool {
	match, _ := regexp.MatchString(`^\(\?=.+\)$`, pattern)
	return match
}

func findMaxValue(values []uint64) uint64 {
	max := uint64(0)
	for _, num := range values {
		if num > max {
			max = num
		}
	}
	return max
}

func findMinValue(values []*uint64) *uint64 {
	dvalues := []uint64{}
	for _, v := range values {
		if v != nil {
			dvalues = append(dvalues, *v)
		}
	}
	if len(dvalues) == 0 {
		return nil
	}
	min := uint64(math.MaxUint64)
	for _, num := range dvalues {
		if num < min {
			min = num
		}
	}
	return openapi3.Uint64Ptr(min)
}

func resolveType(schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	types := filterEmptyStrings(collection.Type)
	if len(types) == 0 {
		schema.Type = ""
		return schema, nil
	}
	if areTypesNumeric(types) {
		for _, t := range types {
			if t == "integer" {
				schema.Type = "integer"
				return schema, nil
			}
		}
		schema.Type = "number"
		return schema, nil
	}
	if allStringsEqual(types) {
		schema.Type = types[0]
		return schema, nil
	}
	return schema, errors.New(TypeErrorMessage)
}

func areTypesNumeric(types []string) bool {
	for _, t := range types {
		if t != "integer" && t != "number" {
			return false
		}
	}
	return true
}

func resolveFormat(schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	formats := filterEmptyStrings(collection.Format)
	if len(formats) == 0 {
		schema.Format = ""
		return schema, nil
	}
	if areFormatsNumeric(formats) {
		schema.Format = resolveNumericFormat(formats)
		return schema, nil
	}
	return defaultFormatResolver(schema, formats)
}

func resolveNumericFormat(formats []string) string {
	orderMap := make(map[string]int)
	orderMap[FormatInt32] = 1
	orderMap[FormatInt64] = 2
	orderMap[FormatFloat] = 3
	orderMap[FormatDouble] = 4
	result := FormatDouble
	for _, format := range formats {
		if orderMap[format] < orderMap[result] {
			result = format
		}
	}
	return result
}

func defaultFormatResolver(schema *openapi3.Schema, formats []string) (*openapi3.Schema, error) {
	if allStringsEqual(formats) {
		schema.Format = formats[0]
		return schema, nil
	}
	return &openapi3.Schema{}, errors.New(FormatErrorMessage)
}

func areFormatsNumeric(values []string) bool {
	for _, val := range values {
		if val != FormatInt32 && val != FormatInt64 && val != FormatFloat && val != FormatDouble {
			return false
		}
	}
	return true
}

func containsString(list []string, search string) bool {
	for _, item := range list {
		if item == search {
			return true
		}
	}
	return false
}

func filterEmptyStrings(input []string) []string {
	var result []string

	for _, s := range input {
		if s != "" {
			result = append(result, s)
		}
	}

	return result
}

func allStringsEqual(values []string) bool {
	first := values[0]
	for _, value := range values {
		if first != value {
			return false
		}
	}
	return true
}

func getIntersection(arr1, arr2 []interface{}) []interface{} {
	intersectionMap := make(map[interface{}]bool)
	result := []interface{}{}

	// Mark elements in the first array
	for _, val := range arr1 {
		intersectionMap[val] = true
	}

	// Check if elements in the second array exist in the intersection map
	for _, val := range arr2 {
		if intersectionMap[val] {
			result = append(result, val)
		}
	}

	return result
}

func findIntersectionOfArrays(arrays [][]interface{}) []interface{} {
	if len(arrays) == 0 {
		return nil
	}

	intersection := arrays[0]

	for i := 1; i < len(arrays); i++ {
		intersection = getIntersection(intersection, arrays[i])
	}
	if len(intersection) == 0 {
		return nil
	}
	return intersection
}

func flattenArray(arrays [][]string) []string {
	var result []string

	for i := 0; i < len(arrays); i++ {
		for j := 0; j < len(arrays[i]); j++ {
			result = append(result, arrays[i][j])
		}
	}

	return result
}

func resolveRequired(values [][]string) []string {
	flatValues := flattenArray(values)
	uniqueMap := make(map[string]bool)
	var uniqueValues []string
	for _, str := range flatValues {
		if _, found := uniqueMap[str]; !found {
			uniqueMap[str] = true
			uniqueValues = append(uniqueValues, str)
		}
	}
	return uniqueValues
}

func collect(schemas []*openapi3.Schema) SchemaCollection {
	collection := SchemaCollection{}
	for _, s := range schemas {
		collection.Not = append(collection.Not, s.Not)
		collection.AnyOf = append(collection.AnyOf, s.AnyOf)
		collection.OneOf = append(collection.OneOf, s.OneOf)
		collection.Title = append(collection.Title, s.Title)
		collection.Type = append(collection.Type, s.Type)
		collection.Format = append(collection.Format, s.Format)
		collection.Description = append(collection.Description, s.Description)
		collection.Enum = append(collection.Enum, s.Enum)
		collection.UniqueItems = append(collection.UniqueItems, s.UniqueItems)
		collection.ExclusiveMin = append(collection.ExclusiveMin, s.ExclusiveMin)
		collection.ExclusiveMax = append(collection.ExclusiveMax, s.ExclusiveMax)
		collection.Min = append(collection.Min, s.Min)
		collection.Max = append(collection.Max, s.Max)
		collection.MultipleOf = append(collection.MultipleOf, s.MultipleOf)
		collection.MinLength = append(collection.MinLength, s.MinLength)
		collection.MaxLength = append(collection.MaxLength, s.MaxLength)
		collection.Pattern = append(collection.Pattern, s.Pattern)
		collection.MinItems = append(collection.MinItems, s.MinItems)
		collection.MaxItems = append(collection.MaxItems, s.MaxItems)
		collection.Items = append(collection.Items, s.Items)
		collection.Required = append(collection.Required, s.Required)
		collection.Properties = append(collection.Properties, s.Properties)
		collection.MinProps = append(collection.MinProps, s.MinProps)
		collection.MaxProps = append(collection.MaxProps, s.MaxProps)
		collection.AdditionalProperties = append(collection.AdditionalProperties, s.AdditionalProperties)
		collection.Nullable = append(collection.Nullable, s.Nullable)
		collection.ReadOnly = append(collection.ReadOnly, s.ReadOnly)
		collection.WriteOnly = append(collection.WriteOnly, s.WriteOnly)
	}
	return collection
}

// getCombinations calculates the cartesian product of groups of SchemaRefs.
func getCombinations(groups []openapi3.SchemaRefs) []openapi3.SchemaRefs {
	if len(groups) == 0 {
		return []openapi3.SchemaRefs{}
	}
	result := []openapi3.SchemaRefs{{}}
	for _, group := range groups {
		var newResult []openapi3.SchemaRefs
		for _, resultItem := range result {
			for _, ref := range group {
				combination := append(openapi3.SchemaRefs{}, resultItem...)
				combination = append(combination, ref)
				newResult = append(newResult, combination)
			}
		}
		result = newResult
	}
	return result
}

func mergeCombinations(state *state, combinations []openapi3.SchemaRefs) ([]*openapi3.Schema, error) {
	merged := []*openapi3.Schema{}
	for _, combination := range combinations {
		schemas := []*openapi3.Schema{}
		for _, ref := range combination {
			schemas = append(schemas, ref.Value)
		}
		schema, err := flattenSchemas(state, schemas)
		if err != nil {
			continue
		}
		merged = append(merged, schema)
	}
	if len(merged) == 0 {
		return merged, errors.New("unable to resolve combined schema")
	}
	return merged, nil
}

func resolveNot(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	refs := filterNilSchemaRef(collection.Not)
	if len(refs) == 0 {
		return schema, nil
	}
	for _, ref := range refs {
		merged, err := Merge(state, *ref)
		if err != nil {
			return &openapi3.Schema{}, err
		}
		ref.Value = merged
	}

	if len(refs) == 1 {
		schema.Not = refs[0]
		return schema, nil
	}
	schema.Not = &openapi3.SchemaRef{
		Value: &openapi3.Schema{
			AnyOf: refs,
		},
	}
	return schema, nil
}

func filterNilSchemaRef(refs []*openapi3.SchemaRef) []*openapi3.SchemaRef {
	result := []*openapi3.SchemaRef{}
	for _, v := range refs {
		if v != nil {
			result = append(result, v)
		}
	}
	return result
}

func filterEmptySchemaRefs(groups []openapi3.SchemaRefs) []openapi3.SchemaRefs {
	result := []openapi3.SchemaRefs{}
	for _, group := range groups {
		if len(group) > 0 {
			result = append(result, group)
		}
	}
	return result
}

func resolveAnyOf(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	refs, err := resolveCombinations(state, collection.AnyOf)
	schema.AnyOf = refs
	return schema, err
}

func resolveOneOf(state *state, schema *openapi3.Schema, collection *SchemaCollection) (*openapi3.Schema, error) {
	refs, err := resolveCombinations(state, collection.OneOf)
	schema.OneOf = refs
	return schema, err
}

func mergeSchemaRefs(state *state, sr []openapi3.SchemaRefs) ([]openapi3.SchemaRefs, error) {
	result := []openapi3.SchemaRefs{}
	for _, refs := range sr {
		r := openapi3.SchemaRefs{}
		for _, ref := range refs {
			merged, err := Merge(state, *ref)
			if err != nil {
				return result, err
			}
			r = append(r, &openapi3.SchemaRef{Value: merged})
		}
		result = append(result, r)
	}
	return result, nil
}

func resolveCombinations(state *state, collection []openapi3.SchemaRefs) (openapi3.SchemaRefs, error) {
	groups := filterEmptySchemaRefs(collection)
	if len(groups) == 0 {
		return nil, nil
	}
	groups, err := mergeSchemaRefs(state, groups)
	if err != nil {
		return openapi3.SchemaRefs{}, err
	}
	// there is only one schema, no need for calculating combinations.
	if len(groups) == 1 {
		return groups[0], nil
	}

	combinations := getCombinations(groups)
	mergedCombinations, err := mergeCombinations(state, combinations)
	if err != nil {
		return nil, err
	}
	var refs openapi3.SchemaRefs
	for _, merged := range mergedCombinations {
		refs = append(refs, &openapi3.SchemaRef{
			Value: merged,
		})
	}
	return refs, err
}

func findIntersection(arrays ...[]string) []string {
	if len(arrays) == 0 {
		return nil
	}

	// Create a map to store the elements of the first array
	elementsMap := make(map[string]bool)
	for _, element := range arrays[0] {
		elementsMap[element] = true
	}

	// Iterate through the remaining arrays and update the map
	for _, arr := range arrays[1:] {
		tempMap := make(map[string]bool)
		for _, element := range arr {
			if elementsMap[element] {
				tempMap[element] = true
			}
		}
		elementsMap = tempMap
	}

	intersection := []string{}
	for element := range elementsMap {
		intersection = append(intersection, element)
	}

	return intersection
}
