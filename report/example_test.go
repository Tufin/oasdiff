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

	s1, err := swaggerLoader.LoadFromFile("../data/openapi-test1.yaml")
	if err != nil {
		fmt.Printf("failed to load spec with %v", err)
		return
	}

	s2, err := swaggerLoader.LoadFromFile("../data/openapi-test3.yaml")
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

	// Output:
	// ### New Endpoints: None
	// -----------------------
	//
	// ### Deleted Endpoints: None
	// ---------------------------
	//
	// ### Modified Endpoints: 4
	// -------------------------
	// GET /api/{domain}/{project}/badges/security-score
	// - Deleted query param: filter
	// - Deleted header param: user
	// - Deleted cookie param: test
	// - Modified query param: image
	//   - Schema changed
	//     - Property 'Not' changed
	//       - Schema added
	//     - Type changed from 'string' to ''
	//     - Format changed from 'general string' to ''
	//     - Description changed from 'alphanumeric' to ''
	//     - Pattern changed from '^(?:[\w-./:]+)$' to ''
	// - Modified query param: token
	//   - Schema changed
	//     - MaxLength changed from 29 to 30
	// - Responses changed
	//   - New response: default
	//   - Deleted response: 200
	//   - Deleted response: 201
	//
	// GET /api/{domain}/{project}/install-command
	// - Deleted header param: network-policies
	// - Modified path param: project
	//   - Schema changed
	//     - New enum values: [test1]
	// - Responses changed
	//   - Modified response: default
	//     - Description changed from 'Tufin1' to 'Tufin'
	//     - Headers changed
	//       - Modified header: X-RateLimit-Limit
	//         - Description changed from 'Request limit per hour.' to 'Request limit per min.'
	// - Servers changed
	//   - New server: https://tufin.io/securecloud
	//   - New server: https://www.tufin.io/securecloud
	//
	// POST /register
	// - Callbacks changed
	// - Security changed
	//   - Deleted security requirements: bearerAuth
	//   - Modified security requirements: OAuth
	//     - Scheme OAuth Added scopes: [write:pets]
	//
	// POST /subscribe
	// - Callbacks changed
	//
	// Security Requirements changed
	// - Deleted security requirements: bearerAuth
	//
	// Servers changed
	// - Deleted server: tufin.com
	//
}

func ExampleGetHTMLReportAsString() {
	loader := openapi3.NewLoader()
	loader.IsExternalRefsAllowed = true

	s1, err := loader.LoadFromFile("../data/openapi-test1.yaml")
	if err != nil {
		fmt.Printf("failed to load spec with %v", err)
		return
	}

	s2, err := loader.LoadFromFile("../data/openapi-test3.yaml")
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

	// Output:
	// <h3 id="new-endpoints-none">New Endpoints: None</h3>
	// <hr>
	// <h3 id="deleted-endpoints-none">Deleted Endpoints: None</h3>
	// <hr>
	// <h3 id="modified-endpoints-4">Modified Endpoints: 4</h3>
	// <hr>
	// <p>GET /api/{domain}/{project}/badges/security-score</p>
	// <ul>
	// <li>Deleted query param: filter</li>
	// <li>Deleted header param: user</li>
	// <li>Deleted cookie param: test</li>
	// <li>Modified query param: image
	// <ul>
	// <li>Schema changed
	// <ul>
	// <li>Property 'Not' changed
	// <ul>
	// <li>Schema added</li>
	// </ul>
	// </li>
	// <li>Type changed from 'string' to ''</li>
	// <li>Format changed from 'general string' to ''</li>
	// <li>Description changed from 'alphanumeric' to ''</li>
	// <li>Pattern changed from '^(?:[\w-./:]+)$' to ''</li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// <li>Modified query param: token
	// <ul>
	// <li>Schema changed
	// <ul>
	// <li>MaxLength changed from 29 to 30</li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// <li>Responses changed
	// <ul>
	// <li>New response: default</li>
	// <li>Deleted response: 200</li>
	// <li>Deleted response: 201</li>
	// </ul>
	// </li>
	// </ul>
	// <p>GET /api/{domain}/{project}/install-command</p>
	// <ul>
	// <li>Deleted header param: network-policies</li>
	// <li>Modified path param: project
	// <ul>
	// <li>Schema changed
	// <ul>
	// <li>New enum values: [test1]</li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// <li>Responses changed
	// <ul>
	// <li>Modified response: default
	// <ul>
	// <li>Description changed from 'Tufin1' to 'Tufin'</li>
	// <li>Headers changed
	// <ul>
	// <li>Modified header: X-RateLimit-Limit
	// <ul>
	// <li>Description changed from 'Request limit per hour.' to 'Request limit per min.'</li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// <li>Servers changed
	// <ul>
	// <li>New server: <a href="https://tufin.io/securecloud">https://tufin.io/securecloud</a></li>
	// <li>New server: <a href="https://www.tufin.io/securecloud">https://www.tufin.io/securecloud</a></li>
	// </ul>
	// </li>
	// </ul>
	// <p>POST /register</p>
	// <ul>
	// <li>Callbacks changed</li>
	// <li>Security changed
	// <ul>
	// <li>Deleted security requirements: bearerAuth</li>
	// <li>Modified security requirements: OAuth
	// <ul>
	// <li>Scheme OAuth Added scopes: [write:pets]</li>
	// </ul>
	// </li>
	// </ul>
	// </li>
	// </ul>
	// <p>POST /subscribe</p>
	// <ul>
	// <li>Callbacks changed</li>
	// </ul>
	// <p>Security Requirements changed</p>
	// <ul>
	// <li>Deleted security requirements: bearerAuth</li>
	// </ul>
	// <p>Servers changed</p>
	// <ul>
	// <li>Deleted server: tufin.com</li>
	// </ul>
}
