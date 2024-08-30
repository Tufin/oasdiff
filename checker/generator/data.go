package generator

import "slices"

func getAll() ValueSets {
	return slices.Concat(
		getPaths(),
		getRequest(),
		getResponse(),
	)
}

func getRequest() ValueSets {
	return slices.Concat(
		schemaValueSets([]string{"media-type", "request body"}, nil),
		schemaValueSets([]string{"property", "media-type", "request body"}, nil),
		schemaValueSets([]string{"request parameter"}, []bool{true}),
	)
}

func getResponse() ValueSets {
	return slices.Concat(
		schemaValueSets([]string{"media-type", "response"}, nil),
		schemaValueSets([]string{"property", "media-type", "response"}, nil),
	)
}

func getPaths() ValueSets {
	return pathValueSets(nil, nil)
}

func pathValueSets(hierarchy []string, attributed []bool) ValueSets {
	return ValueSets{
		ValueSetB{
			adjective:     "%s",
			adjectiveType: PREDICATIVE,
			hierarchy:     hierarchy,
			nouns:         []string{"success response status", "non-success response status"},
			actions:       []string{"add", "remove"},
		},
	}
}

func schemaValueSets(hierarchy []string, attributed []bool) ValueSets {
	return ValueSets{
		ValueSetA{
			adjective:     "value",
			adjectiveType: PREDICATIVE,
			hierarchy:     hierarchy,
			attributed:    attributed,
			nouns:         []string{"max", "maxLength", "min", "minLength", "minItems", "maxItems"},
			actions:       []string{"set", "increase", "decrease"},
		},
		ValueSetA{
			hierarchy: hierarchy,
			nouns:     []string{"type/format"},
			actions:   []string{"change", "generalize"},
		},
		ValueSetA{
			hierarchy: hierarchy,
			nouns:     []string{"discriminator property name"},
			actions:   []string{"change"},
		},
		ValueSetA{
			hierarchy: hierarchy,
			nouns:     []string{"pattern"},
			actions:   []string{"change"},
		},
		ValueSetB{
			adjective:     "%s",
			adjectiveType: PREDICATIVE,
			hierarchy:     hierarchy,
			nouns:         []string{"pattern"},
			actions:       []string{"add", "remove"},
		},
		ValueSetB{
			hierarchy: hierarchy,
			nouns:     []string{"default value"},
			actions:   []string{"add", "remove"},
		},
		ValueSetB{
			adjective:     "%s",
			adjectiveType: PREDICATIVE,
			hierarchy:     append([]string{"anyOf list"}, hierarchy...),
			nouns:         []string{"schema"},
			actions:       []string{"add", "remove"},
		},
		ValueSetB{
			adjective:     "%s",
			adjectiveType: PREDICATIVE,
			hierarchy:     append([]string{"oneOf list"}, hierarchy...),
			nouns:         []string{"schema"},
			actions:       []string{"add", "remove"},
		},
		ValueSetB{
			adjective:     "%s",
			adjectiveType: PREDICATIVE,
			hierarchy:     append([]string{"allOf list"}, hierarchy...),
			nouns:         []string{"schema"},
			actions:       []string{"add", "remove"},
		},
		ValueSetB{
			adjective:     "%s",
			adjectiveType: PREDICATIVE,
			hierarchy:     hierarchy,
			nouns:         []string{"discriminator", "mapping keys"},
			actions:       []string{"add", "remove"},
		},
	}
}
