package noteplan

import (
	"fmt"
	"io/ioutil"
	"os"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/text"
)

func parseMarkdown(path string) ([]byte, ast.Node, error) {

	markdown := goldmark.New(
		goldmark.WithExtensions(),
	)

	file, err := os.Open(path)
	if err != nil {
		Logger.Log(fmt.Sprintf("Failed reading \"%s\"", path))
		return nil, nil, err
	}
	defer file.Close()

	data, err := ioutil.ReadAll(file)
	if err != nil {
		Logger.Log(fmt.Sprintf("Failed reading \"%s\"", path))
		return nil, nil, err
	}

	doc := markdown.Parser().Parse(text.NewReader(data))

	return data, doc, nil
}
