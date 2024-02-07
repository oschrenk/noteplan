package internal

import (
	"fmt"
	"io"
	"os"

	extension "github.com/oschrenk/noteplan/extension"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func parseMarkdown(path string) ([]byte, ast.Node, error) {

	markdown := goldmark.New(
		goldmark.WithExtensions(extension.TaskList))

	file, err := os.Open(path)
	if err != nil {
		Logger.Log(fmt.Sprintf("Failed reading \"%s\"", path))
		return nil, nil, err
	}
	defer file.Close()

	data, err := io.ReadAll(file)
	if err != nil {
		Logger.Log(fmt.Sprintf("Failed reading \"%s\"", path))
		return nil, nil, err
	}

	doc := markdown.Parser().Parse(text.NewReader(data))

	return data, doc, nil
}
