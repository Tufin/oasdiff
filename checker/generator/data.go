package generator

import "slices"

func getAll() ValueSets {
	return slices.Concat(
		getEndpoints(),
		getRequest(),
		getResponse(),
		getComponents(),
	)
}

func getRequest() ValueSets {
	return slices.Concat(
		NewValueSets([]string{"media-type", "request body"}, schemaValueSets),
		NewValueSets([]string{"property", "media-type", "request body"}, schemaValueSets),
		NewValueSets([]string{"request parameter"}, schemaValueSets),
		NewValueSets(nil, operationValueSets),
	)
}

func getResponse() ValueSets {
	return slices.Concat(
		NewValueSets([]string{"media-type", "response"}, schemaValueSets),
		NewValueSets([]string{"property", "media-type", "response"}, schemaValueSets),
	)
}

func getEndpoints() ValueSets {
	return NewValueSets(nil, endpointValueSets)
}

func getComponents() ValueSets {
	return slices.Concat(
		getSecurity(),
	)
}

func getSecurity() ValueSets {
	return NewValueSets(nil, securityValueSets)
}

var securityValueSets = ValueSets{
	ValueSetB{
		predicativeAdjective: "%s",
		objects:              []string{"endpoint scheme security"},
		actions:              []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		hierarchy:            []string{"global security scheme"},
		objects:              []string{"security scope"},
		actions:              []string{"add", "remove"},
	},
}

var endpointValueSets = ValueSets{
	ValueSetA{
		objects: []string{"stability"}, // /Paths/PathItem/Operation
		actions: []string{"decrease"},
	},
	ValueSetA{
		objects: []string{"api path", "api"},
		actions: []string{"remove"},
		adverb:  []string{"without deprecation", "before sunset"},
	},
	ValueSetB{
		objects: []string{"endpoint"}, // /Paths/PathItem
		actions: []string{"add", "remove", "deprecate", "reactivate"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		objects:              []string{"success response status", "non-success response status"}, // /Paths/PathItem/Operation/Responses/Response/content/media-type/
		actions:              []string{"add", "remove"},
	},
	ValueSetA{
		objects: []string{"operation id"},
		actions: []string{"change"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		objects:              []string{"operation id", "tag"},
		actions:              []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		hierarchy:            []string{"endpoint security scheme"},
		objects:              []string{"security scope"},
		actions:              []string{"add", "remove"},
	},
}

var operationValueSets = ValueSets{
	ValueSetB{
		objects: []string{"required request body", "optional request body"},
		actions: []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		attributiveAdjective: "%s",
		objects:              []string{"request parameter"},
		actions:              []string{"add", "remove"},
	},
}

var schemaValueSets = ValueSets{
	ValueSetA{
		predicativeAdjective: "value",
		objects:              []string{"max", "maxLength", "min", "minLength", "minItems", "maxItems"},
		actions:              []string{"set", "increase", "decrease"},
	},
	ValueSetA{
		objects: []string{"type/format"},
		actions: []string{"change", "generalize"},
	},
	ValueSetA{
		objects: []string{"discriminator property name"},
		actions: []string{"change"},
	},
	ValueSetA{
		objects: []string{"pattern"},
		actions: []string{"change", "generalize"},
	},
	ValueSetA{
		objects: []string{"required property", "optional property"},
		actions: []string{"change"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		objects:              []string{"pattern"},
		actions:              []string{"add", "remove"},
	},
	ValueSetB{
		objects: []string{"default value"},
		actions: []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		hierarchy:            []string{"anyOf list"},
		objects:              []string{"schema"},
		actions:              []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		hierarchy:            []string{"anyOf list"},
		objects:              []string{"schema"},
		actions:              []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		hierarchy:            []string{"anyOf list"},
		objects:              []string{"schema"},
		actions:              []string{"add", "remove"},
	},
	ValueSetB{
		predicativeAdjective: "%s",
		objects:              []string{"discriminator", "mapping keys"},
		actions:              []string{"add", "remove"},
	},
}

/*
missing:
api-deprecated-sunset-parse
api-path-sunset-parse
api-invalid-stability-level
api-deprecated-sunset-missing
api-sunset-date-too-small
*/
