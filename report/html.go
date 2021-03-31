package report

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/tufin/oasdiff/diff"
)

// GetHTMLReportAsString returns an HTML diff report as a string
func GetHTMLReportAsString(d *diff.Diff) (string, error) {

	return markdownToHTML(GetTextReportAsBytes(d))
}

func markdownToHTML(source []byte) (string, error) {
	md := goldmark.New(
		goldmark.WithExtensions(extension.GFM),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(),
		),
	)

	var buf bytes.Buffer
	if err := md.Convert(source, &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
