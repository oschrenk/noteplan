package extension

import (
	"regexp"

	"github.com/oschrenk/noteplan/extension/ast"

	"github.com/yuin/goldmark"
	gast "github.com/yuin/goldmark/ast"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/text"
	"github.com/yuin/goldmark/util"
)

var taskListRegexp = regexp.MustCompile(`^\[([\sxX])\]\s*`)

type taskCheckBoxParser struct {
}

var defaultTaskCheckBoxParser = &taskCheckBoxParser{}

// NewTaskCheckBoxParser returns a new  InlineParser that can parse
// checkboxes in list items.
// This parser must take precedence over the parser.LinkParser.
func NewTaskCheckBoxParser() parser.InlineParser {
	return defaultTaskCheckBoxParser
}

func (s *taskCheckBoxParser) Trigger() []byte {
	return []byte{'['}
}

func (s *taskCheckBoxParser) Parse(parent gast.Node, block text.Reader, pc parser.Context) gast.Node {
	// Given AST structure must be like
	// - List
	//   - ListItem         : parent.Parent
	//     - TextBlock      : parent
	//       (current line)
	if parent.Parent() == nil || parent.Parent().FirstChild() != parent {
		return nil
	}

	if parent.HasChildren() {
		return nil
	}
	if _, ok := parent.Parent().(*gast.ListItem); !ok {
		return nil
	}
	line, _ := block.PeekLine()
	m := taskListRegexp.FindSubmatchIndex(line)
	if m == nil {
		return nil
	}
	value := line[m[2]:m[3]][0]
	block.Advance(m[1])
	state := ast.NewTaskState(value)
	return ast.NewTaskCheckBox(state)
}

func (s *taskCheckBoxParser) CloseBlock(parent gast.Node, pc parser.Context) {
	// nothing to do
}

type taskList struct {
}

// TaskList is an extension that allow you to use GFM task lists.
var TaskList = &taskList{}

func (e *taskList) Extend(m goldmark.Markdown) {
	m.Parser().AddOptions(parser.WithInlineParsers(
		util.Prioritized(NewTaskCheckBoxParser(), 0),
	))
}
