package util

import (
	"bytes"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"
)

// MarkdownToHTML converts Markdown content to HTML
// Supports common extensions like tables, strikethrough, etc.
func MarkdownToHTML(markdown string) (string, error) {
	// Configure goldmark with extensions
	md := goldmark.New(
		goldmark.WithExtensions(
			extension.GFM,        // GitHub Flavored Markdown (tables, strikethrough, etc.)
			extension.Typographer, // Smart quotes, dashes, etc.
		),
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(), // Auto-generate heading IDs
		),
		goldmark.WithRendererOptions(
			html.WithHardWraps(), // Convert line breaks to <br>
			html.WithXHTML(),     // Use XHTML-style tags
		),
	)

	var buf bytes.Buffer
	if err := md.Convert([]byte(markdown), &buf); err != nil {
		return "", err
	}

	return buf.String(), nil
}
