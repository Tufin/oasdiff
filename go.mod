module github.com/tufin/oasdiff

go 1.15

require (
	github.com/getkin/kin-openapi v0.39.0
	github.com/sirupsen/logrus v1.7.0
	github.com/stretchr/testify v1.5.1
)

replace "github.com/tufin/oasdiff/diff" => ./diff