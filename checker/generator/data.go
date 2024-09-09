package generator

import "slices"

func GetAll() (MessageGenerator, error) {
	return slices.Concat(
		getEndpoints(),
		getRequest(),
		getResponse(),
		getComponents(),
	), nil
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
		PredicativeAdjective: "%s",
		Names:                []string{"endpoint scheme security"},
		Actions:              []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Hierarchy:            []string{"global security scheme"},
		Names:                []string{"security scope"},
		Actions:              []string{"add", "remove"},
	},
}

var endpointValueSets = ValueSets{
	ValueSetA{
		Names:   []string{"stability"}, // /Paths/PathItem/Operation
		Actions: []string{"decrease"},
	},
	ValueSetA{
		Names:   []string{"api path", "api"},
		Actions: []string{"remove"},
		Adverbs: []string{"without deprecation", "before sunset"},
	},
	ValueSetB{
		Names:   []string{"endpoint"}, // /Paths/PathItem
		Actions: []string{"add", "remove", "deprecate", "reactivate"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Names:                []string{"success response status", "non-success response status"}, // /Paths/PathItem/Operation/Responses/Response/content/media-type/
		Actions:              []string{"add", "remove"},
	},
	ValueSetA{
		Names:   []string{"operation id"},
		Actions: []string{"change"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Names:                []string{"operation id", "tag"},
		Actions:              []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Hierarchy:            []string{"endpoint security scheme"},
		Names:                []string{"security scope"},
		Actions:              []string{"add", "remove"},
	},
}

var operationValueSets = ValueSets{
	ValueSetB{
		Names:   []string{"required request body", "optional request body"},
		Actions: []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		AttributiveAdjective: "%s",
		Names:                []string{"request parameter"},
		Actions:              []string{"add", "remove"},
	},
}

var schemaValueSets = ValueSets{
	ValueSetA{
		PredicativeAdjective: "value",
		Names:                []string{"max", "maxLength", "min", "minLength", "minItems", "maxItems"},
		Actions:              []string{"set", "increase", "decrease"},
	},
	ValueSetA{
		Names:   []string{"type/format"},
		Actions: []string{"change", "generalize"},
	},
	ValueSetA{
		Names:   []string{"discriminator property name"},
		Actions: []string{"change"},
	},
	ValueSetA{
		Names:   []string{"pattern"},
		Actions: []string{"change", "generalize"},
	},
	ValueSetA{
		Names:   []string{"required property", "optional property"},
		Actions: []string{"change"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Names:                []string{"pattern"},
		Actions:              []string{"add", "remove"},
	},
	ValueSetB{
		Names:   []string{"default value"},
		Actions: []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Hierarchy:            []string{"anyOf list"},
		Names:                []string{"schema"},
		Actions:              []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Hierarchy:            []string{"anyOf list"},
		Names:                []string{"schema"},
		Actions:              []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Hierarchy:            []string{"anyOf list"},
		Names:                []string{"schema"},
		Actions:              []string{"add", "remove"},
	},
	ValueSetB{
		PredicativeAdjective: "%s",
		Names:                []string{"discriminator", "mapping keys"},
		Actions:              []string{"add", "remove"},
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
