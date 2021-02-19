package main

import (
	"encoding/json"
	"flag"
	"fmt"

	log "github.com/sirupsen/logrus"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

var base, revision, prefix, filter string

func init() {
	flag.StringVar(&base, "base", "", "base swagger file")
	flag.StringVar(&revision, "revision", "", "revised swagger file")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base swagger but not in revision swagger")
	flag.StringVar(&filter, "filter", "", "regex to filter results by endpoints")
}

func main() {
	flag.Parse()

	loader := load.NewSwaggerLoader()

	base, err := loader.From(base)
	if err != nil {
		return
	}

	revision, err := loader.From(revision)
	if err != nil {
		return
	}

	bytes, err := json.MarshalIndent(diff.Run(base, revision, prefix, filter), "", " ")
	if err != nil {
		log.Errorf("failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
