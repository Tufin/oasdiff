[![codecov](https://codecov.io/gh/Tufin/oasdiff/branch/master/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/Tufin/oasdiff)

# OpenAPI Spec Diff
A diff tool for OpenAPI Spec 3.  

## Unique features vs. other OAS3 diff tools
- go module
- deep diff into paths, schemas, parameters, responses, enums etc.

## Build
```
git clone https://github.com/Tufin/oasdiff.git
cd oasdiff
go build
```

## Running from the command-line
```
./oasdiff -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```

## Help
```
./oasdiff --help
```

## Embedding into your Go program
```
package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func main() {
	loader := load.NewOASLoader()

	base, err := loader.From("v1.yaml")
	if err != nil {
		return
	}

	revision, err := loader.From("v2.yaml")
	if err != nil {
		return
	}

	bytes, err := json.MarshalIndent(diff.Get(base, revision, "", ""), "", " ")
	if err != nil {
		log.Errorf("failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
```
