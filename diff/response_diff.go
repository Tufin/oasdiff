package diff

import "github.com/getkin/kin-openapi/openapi3"

type ResponseDiff struct {
}

func (responseDiff ResponseDiff) empty() bool {
	return responseDiff == ResponseDiff{}
}

func diffResponseValues(response1, response2 *openapi3.Response) ResponseDiff {
	return ResponseDiff{}
}
