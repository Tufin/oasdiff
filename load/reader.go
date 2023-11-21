package load

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"

	"github.com/getkin/kin-openapi/openapi3"
)

var ReaderFunc = openapi3.URIMapCache(openapi3.ReadFromURIs(ReadFromHTTP(http.DefaultClient), openapi3.ReadFromFile))

func ReadFromHTTP(cl *http.Client) openapi3.ReadFromURIFunc {
	return func(loader *openapi3.Loader, location *url.URL) ([]byte, error) {
		if location.Scheme == "" || location.Host == "" {
			return nil, openapi3.ErrURINotSupported
		}
		req, err := http.NewRequest("GET", location.String(), nil)
		if err != nil {
			return nil, err
		}
		resp, err := cl.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()
		if resp.StatusCode > 399 {
			return nil, fmt.Errorf("error loading %q: request returned status code %d", location.String(), resp.StatusCode)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}

		return replaceExternalLinks(body)
	}
}

func traverse(m map[string]any, ind int, state int) bool {
	for k, v := range m {
		// fmt.Printf(strings.Repeat("  ", ind)+"%s\n", k)
		if m1, ok := v.(map[string]any); ok {
			if k == "schema" {
				state = 1
			}
			if traverse(m1, ind+1, state) {
				m[k] = map[string]string{"type": "string"}
			}
		} else {
			if _, ok := v.(string); ok {
				if state == 1 && k == "$ref" /*&& strings.HasPrefix(s, "../")*/ {
					return true
				}
			}
			// fmt.Printf(strings.Repeat("  ", ind+1)+"%v\n", v)
		}
	}
	return false
}

func replaceExternalLinks(data []byte) ([]byte, error) {
	var objmap map[string]any
	err := json.Unmarshal(data, &objmap)
	if err != nil {
		return nil, err
	}

	traverse(objmap, 0, 0)
	data, err = json.Marshal(objmap)
	if err != nil {
		return nil, err
	}
	return data, nil
}
