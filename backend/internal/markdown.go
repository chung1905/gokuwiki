package internal

import (
	"github.com/gomarkdown/markdown"
	"github.com/gomarkdown/markdown/parser"
)

func Md2html(md []byte) []byte {
	extensions := parser.CommonExtensions | parser.HardLineBreak | parser.FencedCode
	parserModel := parser.NewWithExtensions(extensions)
	content := markdown.NormalizeNewlines(md)
	output := markdown.ToHTML(content, parserModel, nil)

	return output
}

func NormalizeNewlines(md []byte) []byte {
	return markdown.NormalizeNewlines(md)
}
