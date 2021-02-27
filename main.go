package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	log "github.com/sirupsen/logrus"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

var base, revision, prefix, filter string

func init() {
	flag.StringVar(&base, "base", "", "original OpenAPI spec")
	flag.StringVar(&revision, "revision", "", "revised OpenAPI spec")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base spec but not the revision")
	flag.StringVar(&filter, "filter", "", "regex to filter result paths")
}

func main() {
	flag.Parse()

	swaggerLoader := openapi3.NewSwaggerLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	loader := load.NewOASLoader(swaggerLoader)

	base, err := loader.From(base)
	if err != nil {
		return
	}

	revision, err := loader.From(revision)
	if err != nil {
		return
	}

	bytes, err := json.MarshalIndent(diff.Get(base, revision, prefix, filter), "", " ")
	if err != nil {
		log.Errorf("failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
