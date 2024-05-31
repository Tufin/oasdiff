package diff_test

import (
	"fmt"
	"os"
	"strings"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
)

func ExampleGet() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := loader.LoadFromFile("../data/simple1.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec: %v", err)
		return
	}

	s2, err := loader.LoadFromFile("../data/simple2.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec: %v", err)
		return
	}

	diffReport, err := diff.Get(diff.NewConfig(), s1, s2)

	if err != nil {
		fmt.Fprintf(os.Stderr, "diff failed with %v", err)
		return
	}

	bytes, err := yaml.Marshal(diffReport)
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to marshal result with %v", err)
		return
	}
	fmt.Printf("%s\n", bytes)

	// Output:
	// paths:
	//     modified:
	//         /api/test:
	//             operations:
	//                 added:
	//                     - POST
	//                 deleted:
	//                     - GET
	// endpoints:
	//     added:
	//         - method: POST
	//           path: /api/test
	//     deleted:
	//         - method: GET
	//           path: /api/test
}

func ExampleGetPathsDiff() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := load.NewSpecInfo(loader, load.NewSource("../data/openapi-test1.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec: %v", err)
		return
	}

	s2, err := load.NewSpecInfo(loader, load.NewSource("../data/openapi-test3.yaml"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec: %v", err)
		return
	}

	diffRes, operationsSources, err := diff.GetPathsDiff(diff.NewConfig(),
		[]*load.SpecInfo{s1},
		[]*load.SpecInfo{s2},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "diff failed with %v", err)
		return
	}

	errs := checker.CheckBackwardCompatibility(checker.NewConfig(), diffRes, operationsSources)

	// process configuration file for ignoring errors
	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, "../data/ignore-err-example.txt", checker.NewDefaultLocalizer())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ignore errors failed with %v", err)
		return
	}

	// process configuration file for ignoring warnings
	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.WARN, errs, "../data/ignore-warn-example.txt", checker.NewDefaultLocalizer())
	if err != nil {
		fmt.Fprintf(os.Stderr, "ignore warnings failed with %v", err)
		return
	}

	// pretty print breaking changes errors
	if len(errs) > 0 {
		localizer := checker.NewDefaultLocalizer()
		count := errs.GetLevelCount()
		fmt.Print(localizer("total-errors", len(errs), count[checker.ERR], "error", count[checker.WARN], "warning"))
		for _, bcerr := range errs {
			fmt.Printf("%s\n\n", strings.TrimRight(bcerr.SingleLineError(localizer, checker.ColorNever), " "))
		}
	}

	// Output:
	// 4 breaking changes: 1 error, 3 warning
	// error at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score removed the success response with the status '201' [response-success-status-removed].
	//
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score deleted the 'cookie' request parameter 'test' [request-parameter-removed].
	//
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score deleted the 'header' request parameter 'user' [request-parameter-removed].
	//
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score deleted the 'query' request parameter 'filter' [request-parameter-removed].
	//
}
