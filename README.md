# OpenAPI Spec Diff
A go module for reporting changes between versions of OpenAPI (Swagger) files.    
The diff report is a go struct which can also be marshalled like this:
```json
{
 "diffResult": {
  "pathDiff": {
   "modifiedEndpoints": {
    "/prefix/api/{domain}/{project}/badges/security-score/": {
     "modifiedMethods": {
      "GET": {
       "addedParams": {
        "query": {
         "filter": {}
        }
       },
       "modifiedParams": {
        "path": {
         "domain": {
          "schemaDiff": {
           "minDiff": {
            "oldValue": null,
            "newValue": 7
           }
          }
         }
        },
        "query": {
         "token": {
          "schemaDiff": {
           "anyOfDiff": true
          }
         }
        }
       }
      }
     }
    }
   }
  }
 },
 "diffSummary": {
  "diff": true,
  "addedEndpoints": 0,
  "deletedEndpoints": 0,
  "modifiedEndpoints": 1
 }
}
```

## Build
```
git clone https://github.com/Tufin/oasdiff.git
cd oasdiff
go build
```

## Running from the command-line
```
./oasdiff -base data/openapi-test1.yaml -revision data/openapi-test2.yaml
```

## Help
```
./oasdiff --help
```

## Embedding into your Go program
```
package main

import (
	"encoding/json"
	"fmt"

	log "github.com/sirupsen/logrus"

	"github.com/tufin/oasdiff/diff"
	"github.com/tufin/oasdiff/load"
)

func main() {
	base, err := load.Load("../oasdiff/data/openapi-test1.yaml")
	if err != nil {
		return
	}

	revision, err := load.Load("../oasdiff/data/openapi-test2.yaml")
	if err != nil {
		return
	}

	bytes, err := json.MarshalIndent(diff.GetDiffResponse(base, revision, "", ""), "", " ")
	if err != nil {
		log.Errorf("failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
```
