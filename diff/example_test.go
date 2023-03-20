package diff_test

import (
	"fmt"
	"os"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/tufin/oasdiff/checker"
	"github.com/tufin/oasdiff/checker/localizations"
	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
	"gopkg.in/yaml.v3"
)

func ExampleGet() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := loader.LoadFromFile("../data/simple1.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec with %v", err)
		return
	}

	s2, err := loader.LoadFromFile("../data/simple2.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec with %v", err)
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

	s1, err := checker.LoadOpenAPISpecInfo("../data/openapi-test1.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec with %v", err)
		return
	}

	s2, err := checker.LoadOpenAPISpecInfo("../data/openapi-test3.yaml")
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to load spec with %v", err)
		return
	}

	diffConfig := diff.NewConfig()
	diffConfig.IncludeExtensions.Add(checker.XStabilityLevelExtension)
	diffConfig.IncludeExtensions.Add(diff.SunsetExtension)
	diffConfig.IncludeExtensions.Add(checker.XExtensibleEnumExtension)

	diffRes, operationsSources, err := diff.GetPathsDiff(diffConfig,
		[]load.OpenAPISpecInfo{*s1},
		[]load.OpenAPISpecInfo{*s2},
	)

	if err != nil {
		fmt.Fprintf(os.Stderr, "diff failed with %v", err)
		return
	}

	c := checker.GetDefaultChecks()
	c.Localizer = *localizations.New("en", "en")
	errs := checker.CheckBackwardCompatibility(c, diffRes, operationsSources)

	// process configuration file for ignoring errors
	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, "../data/ignore-err-example.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ignore errors failed with %v", err)
		return
	}

	// process configuration file for ignoring warnings
	errs, err = checker.ProcessIgnoredBackwardCompatibilityErrors(checker.ERR, errs, "../data/ignore-warn-example.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "ignore warnings failed with %v", err)
		return
	}

	// pretty print breaking changes errors
	if len(errs) > 0 {
		fmt.Printf(c.Localizer.Get("messages.total-errors"), len(errs))
		for _, bcerr := range errs {
			fmt.Printf("%s\n\n", bcerr.PrettyError(c.Localizer))
		}
	}

	// Backward compatibility errors (5):
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score deleted the 'query' request parameter 'filter' [request-parameter-removed].
	//
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score deleted the 'header' request parameter 'user' [request-parameter-removed].
	//
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score deleted the 'cookie' request parameter 'test' [request-parameter-removed].
	//
	// error at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/badges/security-score removed the success response with the status '201' [response-success-status-removed].
	//
	// warning at ../data/openapi-test3.yaml, in API GET /api/{domain}/{project}/install-command deleted the 'header' request parameter 'network-policies' [request-parameter-removed].
	//
}
