package generator

import "slices"

func getAll() ValueSets {
	return slices.Concat(
		getPaths(),
		getRequest(),
		getResponse(),
		getComponents(),
	)
}

func getRequest() ValueSets {
	return slices.Concat(
		NewValueSets([]string{"media-type", "request body"}, nil, schemaValueSets),
		NewValueSets([]string{"property", "media-type", "request body"}, nil, schemaValueSets),
		NewValueSets([]string{"request parameter"}, []bool{true}, schemaValueSets),
		NewValueSets(nil, nil, operationValueSets),
	)
}

func getResponse() ValueSets {
	return slices.Concat(
		NewValueSets([]string{"media-type", "response"}, nil, schemaValueSets),
		NewValueSets([]string{"property", "media-type", "response"}, nil, schemaValueSets),
	)
}

func getPaths() ValueSets {
	return NewValueSets(nil, nil, pathsValueSets)
}

func getComponents() ValueSets {
	return slices.Concat(
		getSecurity(),
	)
}

func getSecurity() ValueSets {
	return NewValueSets(nil, nil, securityValueSets)
}

var securityValueSets = ValueSets{
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		nouns:         []string{"endpoint scheme security"},
		actions:       []string{"add", "remove"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		hierarchy:     []string{"global security scheme"},
		attributed:    []bool{false},
		nouns:         []string{"security scope"},
		actions:       []string{"add", "remove"},
	},
}

var pathsValueSets = ValueSets{
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		nouns:         []string{"success response status", "non-success response status"},
		actions:       []string{"add", "remove"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		hierarchy:     []string{"endpoint security scheme"},
		attributed:    []bool{false},
		nouns:         []string{"security scope"},
		actions:       []string{"add", "remove"},
	},
}

var operationValueSets = ValueSets{
	ValueSetB{
		nouns:   []string{"required request body", "optional request body"},
		actions: []string{"add", "remove"},
	},
}

var schemaValueSets = ValueSets{
	ValueSetA{
		adjective:     "value",
		adjectiveType: PREDICATIVE,
		nouns:         []string{"max", "maxLength", "min", "minLength", "minItems", "maxItems"},
		actions:       []string{"set", "increase", "decrease"},
	},
	ValueSetA{
		nouns:   []string{"type/format"},
		actions: []string{"change", "generalize"},
	},
	ValueSetA{
		nouns:   []string{"discriminator property name"},
		actions: []string{"change"},
	},
	ValueSetA{
		nouns:   []string{"pattern"},
		actions: []string{"change"},
	},
	ValueSetA{
		nouns:   []string{"required property", "optional property"},
		actions: []string{"change"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		nouns:         []string{"pattern"},
		actions:       []string{"add", "remove"},
	},
	ValueSetB{
		nouns:   []string{"default value"},
		actions: []string{"add", "remove"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		hierarchy:     []string{"anyOf list"},
		attributed:    []bool{false},
		nouns:         []string{"schema"},
		actions:       []string{"add", "remove"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		hierarchy:     []string{"anyOf list"},
		attributed:    []bool{false},
		nouns:         []string{"schema"},
		actions:       []string{"add", "remove"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		hierarchy:     []string{"anyOf list"},
		attributed:    []bool{false},
		nouns:         []string{"schema"},
		actions:       []string{"add", "remove"},
	},
	ValueSetB{
		adjective:     "%s",
		adjectiveType: PREDICATIVE,
		nouns:         []string{"discriminator", "mapping keys"},
		actions:       []string{"add", "remove"},
	},
}
