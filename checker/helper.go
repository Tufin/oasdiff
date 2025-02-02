package checker

import (
	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
)

// helper is a struct that can be used as a parameter to the checker functions.
// it is used to simplify the function signature and to provide a more readable API.
type helper struct {
	config            *Config
	operation         *openapi3.Operation
	operationsSources *diff.OperationsSourcesMap
	method            string
	path              string
}

func newHelper(config *Config, operation *openapi3.Operation, operationsSources *diff.OperationsSourcesMap, method, path string) helper {
	return helper{
		config:            config,
		operation:         operation,
		operationsSources: operationsSources,
		method:            method,
		path:              path,
	}
}
