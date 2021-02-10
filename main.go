package main

import (
	"encoding/json"
	"flag"
	"fmt"

	"github.com/apex/log"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

var manualSwagger, autoSwagger, prefix, filter string

func init() {
	flag.StringVar(&manualSwagger, "base", "", "base swagger file")
	flag.StringVar(&autoSwagger, "revision", "", "revised swagger file")
	flag.StringVar(&prefix, "prefix", "", "path prefix that exists in base swagger but not in revision swagger")
	flag.StringVar(&filter, "filter", "", "regex to filter results")
}

func main() {
	flag.Parse()

	manual, err := load.Load(manualSwagger)
	if err != nil {
		return
	}

	auto, err := load.Load(autoSwagger)
	if err != nil {
		return
	}

	result := diff.Diff(auto, manual, prefix)
	result.FilterByRegex(filter)

	bytes, err := json.MarshalIndent(result, "", " ")
	if err != nil {
		log.Errorf("failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
