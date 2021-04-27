package report_test

import (
	"fmt"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/report"
)

func ExampleGetTextReportAsString() {
	swaggerLoader := openapi3.NewLoader()
	swaggerLoader.IsExternalRefsAllowed = true

	s1, err := swaggerLoader.LoadFromFile("../../data/openapi-test1.yaml")
	if err != nil {
		fmt.Printf("failed to load spec with %v", err)
		return
	}

	s2, err := swaggerLoader.LoadFromFile("../../data/openapi-test3.yaml")
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

func ExampleGetHTMLReportAsString() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := loader.LoadFromFile("../../data/openapi-test1.yaml")
	if err != nil {
		fmt.Printf("failed to load spec with %v", err)
		return
	}

	s2, err := loader.LoadFromFile("../../data/openapi-test3.yaml")
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
