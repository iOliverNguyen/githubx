package github

import (
	"regexp"
	"strings"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
)

var md = goldmark.New(
	goldmark.WithExtensions(extension.GFM),
)

var reHTML = regexp.MustCompile(`<[^>]*>`)

func ConvertMD(body string) (bodyHTML, bodyText string, _ error) {
	var b strings.Builder
	err := md.Convert([]byte(body), &b)
	bodyHTML = b.String()
	bodyText = reHTML.ReplaceAllString(bodyHTML, "")
	return bodyHTML, bodyText, err
}
