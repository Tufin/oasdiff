package report_test

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

func ExampleText() {
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

	fmt.Print(report.GetTextReportAsString(diffReport))
}

func ExampleHTML() {
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

	html, err := report.GetHTMLReportAsString(diffReport)
	if err != nil {
		fmt.Printf("failed to generate HTML with %v", err)
		return
	}

	fmt.Print(html)
}
