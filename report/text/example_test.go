package text_test

import (
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report/text"
)

func ExampleReport_Output() {
	swaggerLoader := openapi3.NewSwaggerLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	s1, err := swaggerLoader.LoadSwaggerFromFile("../../data/openapi-test1.yaml")
	if err != nil {
		fmt.Printf("failed to load spec with %v", err)
		return
	}

	s2, err := swaggerLoader.LoadSwaggerFromFile("../../data/openapi-test3.yaml")
	if err != nil {
		fmt.Printf("failed to load spec with %v", err)
		return
	}

	diffReport, err := diff.Get(&diff.Config{}, s1, s2)
	if err != nil {
		fmt.Printf("diff failed with %v", err)
		return
	}

	report := text.Report{
		Writer: os.Stdout,
	}
	report.Output(diffReport)
}
