package internal

import (
	"bytes"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/fatih/color"
	"github.com/yuin/goldmark/ast"
)

type TaskCategory int64
type TaskState int64

const (
	Bullet TaskCategory = iota
	Checklist
	Todo
)

const (
	Open TaskState = iota
	Cancelled
	Done
)

func (s TaskState) Trigger() string {
	switch s {
	case Open:
		return " "
	case Cancelled:
		return "-"
	case Done:
		return "x"
	}

	log.Panicf("failed getting trigger for state: %s", s)
	return "unknown"
}

type Task struct {
	Category TaskCategory
	State    TaskState
	Text     string
	Depth    int
}

func (t Task) String() string {
	var char string
	switch {
	case t.Category == Bullet:
		char = "" // nf-oct-dot_fill, \uf444
	case t.Category == Todo:
		switch {
		case t.State == Open:
			char = "󰝦" // nf-md-checkbox_blank_circle_outline, \udb81\udf66
		case t.State == Cancelled:
			char = "" // nf-oct-x_circle, \uf52f
		case t.State == Done:
			char = "󰄴" // nf-md-checkbox_marked_circle_outline, \udb80\udd34
		}
	case t.Category == Checklist:
		switch {
		case t.State == Open:
			char = "" // nf-seti-checkbox_unchecked, \ue640
		case t.State == Cancelled:
			char = "󱋭" // nf-md-checkbox_blank_off_outline, \uf52f
		case t.State == Done:
			char = "󰄵" // nf-md-checkbox_marked_outline, \udb80\udd35
		}
	}
	indent := strings.Repeat(" ", t.Depth*2)

	grey := color.New(color.Faint, color.FgWhite).SprintFunc()
	none := color.New().SprintFunc()

	colorize := none
	if t.State == Done {
		colorize = grey
	}

	return fmt.Sprintf("%s%s %s", indent, colorize(char), colorize(t.Text))
}

// There are two different "trigger" characters for bulleted lists:
//
// - Asterisk `*`
// - Dash `-`
//
// Starting a line with one of those characters followed by a whitespace will
// trigger a bullet list item.
//
// If the whitespace is followed by `[ ]` (bracketed whitespace), the bullet
// list item is transformed into an open todo.
//
// Noteplan saves the user some time by interpreting a trigger character with
// whitespace as an open task without the bracketed whitespace.
//
// # There are two settings
//
// - Recognize `*` as Todo
// - Recognize `-` as Todo
//
// Additionally there is a trigger character for a "checlist" item
//
// - Plus `+`
//
// Examples:
// + [ ] open checklist
// + open checklist
// + [-] cancelled checklist
// + [x] done checklist
// * open todo
// - [ ] open alternative tod
// * [-] cancelled todo
// - [x] done todo
// - bullet
// *  double spaced todo
// -  double spaced bullet
func (noteplan *Noteplan) parseTask(marker string, task string, depth int) Task {
	var category TaskCategory

	switch {
	case marker == "*":
		if noteplan.settings.IsAsteriskTodo {
			category = Todo
		} else {
			category = Bullet
		}
	case marker == "-":
		if noteplan.settings.IsDashTodo {
			category = Todo
		} else {
			category = Bullet
		}
	case marker == "+":
		category = Checklist
	}
	text := strings.TrimSpace(task)

	return noteplan.parseTaskState(category, text, depth)
}

func (noteplan *Noteplan) parseTaskState(category TaskCategory, text string, depth int) Task {
	var state TaskState

	stateRegEx := `^\[([\s-x])\].*`
	re := regexp.MustCompile(stateRegEx)
	matches := re.FindStringSubmatch(text)

	// if any state trigger is found we assume it's a Todo
	if len(matches) == 2 {
		// the user might have changed their settings in the past or might have
		// manually changed a bullet item into a task item by using task markers
		// if it's checklist we keep it
		if category == Bullet {
			category = Todo
		}
		stateChar := matches[1]
		text = strings.TrimSpace(text[3:])
		switch {
		case stateChar == Open.Trigger():
			state = Open
		case stateChar == Cancelled.Trigger():
			state = Cancelled
		case stateChar == Done.Trigger():
			state = Done
		}

		// otherwise assume Open
	} else {
		state = Open
	}
	// if bullet do not return TaskState
	if category == Bullet {
		return Task{Category: category, Text: text, Depth: depth}
	}

	return Task{Category: category, State: state, Text: text, Depth: depth}
}

func getText(n ast.Node, source []byte) string {
	if n.Type() == ast.TypeBlock {
		var text bytes.Buffer
		for i := 0; i < n.Lines().Len(); i++ {
			line := n.Lines().At(i)
			text.Write(line.Value(source))
		}
		return text.String()
	}
	return ""
}

func (noteplan *Noteplan) parseTasks(data []byte, doc ast.Node) []Task {
	var tasks []Task
	var depth = -1
	markerMap := make(map[int]string)

	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		if n, ok := node.(*ast.List); ok {
			if enter {
				depth = depth + 1
				markerMap[depth] = string((*n).Marker)
			} else {
				depth = depth - 1
			}
		}

		if n, ok := node.(*ast.ListItem); ok && enter {
			item := n.FirstChild()
			text := getText(item, data)
			task := noteplan.parseTask(markerMap[depth], text, depth)
			tasks = append(tasks, task)
		}

		return ast.WalkContinue, nil
	})

	return tasks
}
