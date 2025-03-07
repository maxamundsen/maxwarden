package ui

import (
	"github.com/gomarkdown/markdown"
	. "maragu.dev/gomponents"
)

func Markdown(input string) Node {
	html := markdown.ToHTML([]byte(input), nil, nil)
	return SafeRaw(string(html))
}
