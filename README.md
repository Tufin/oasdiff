[![codecov](https://codecov.io/gh/Tufin/oasdiff/branch/master/graph/badge.svg?token=Y8BM6X77JY)](https://codecov.io/gh/Tufin/oasdiff)

# OpenAPI Spec Diff
A go module for reporting changes between versions of OpenAPI (Swagger) files.    
The diff report is a go struct which can also be marshalled like this:
```json
{
 "diff": {
  "endpoints": {
   "modified": {
    "/api/{domain}/{project}/badges/security-score": {
     "operations": {
      "modified": {
       "GET": {
        "parameters": {
         "modified": {
          "cookie": {
           "test": {
            "content": {
             "mediaTypeDiff": true
            }
           }
          },
          "header": {
           "user": {
            "schema": {
             "schemaDeleted": true
            },
            "content": {
             "mediaTypeAdded": true
            }
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
  "schemas": {
   "deleted": [
    "requests"
   ],
   "modified": {
    "network-policies": {
     "additionalPropertiesAllowed": {
      "oldValue": false,
      "newValue": true
     }
    },
    "rules": {
     "additionalPropertiesAllowed": {
      "oldValue": false,
      "newValue": null
     }
    }
   }
  }
 },
 "summary": {
  "diff": true,
  "paths": {
   "modified": 1
  },
  "schemas": {
   "deleted": 1,
   "modified": 2
  }
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
	loader := load.NewSwaggerLoader()

	base, err := loader.From("v1.yaml")
	if err != nil {
		return
	}

	revision, err := loader.From("v2.yaml")
	if err != nil {
		return
	}

	bytes, err := json.MarshalIndent(diff.Run(base, revision, "", ""), "", " ")
	if err != nil {
		log.Errorf("failed to marshal result with '%v'", err)
		return
	}

	fmt.Printf("%s\n", bytes)
}
```
