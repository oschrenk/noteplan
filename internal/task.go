package internal

import (
	"fmt"
	"strings"

	tasklist "github.com/oschrenk/noteplan/extension/tasklist"

	"github.com/fatih/color"
	"github.com/yuin/goldmark/ast"
)

type TaskCategory int64

const (
	Bullet TaskCategory = iota
	Checklist
	Todo
)

func (t TaskCategory) String() string {
	switch {
	case t == Bullet:
		return "*"
	case t == Checklist:
		return "+"
	case t == Todo:
		return "-"
	}
	return "?"
}

type Task struct {
	Category TaskCategory
	State    tasklist.TaskState
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
		case t.State == tasklist.Open:
			char = "󰝦" // nf-md-checkbox_blank_circle_outline, \udb81\udf66
		case t.State == tasklist.Cancelled:
			char = "" // nf-oct-x_circle, \uf52f
		case t.State == tasklist.Done:
			char = "󰄴" // nf-md-checkbox_marked_circle_outline, \udb80\udd34
		}
	case t.Category == Checklist:
		switch {
		case t.State == tasklist.Open:
			char = "" // nf-seti-checkbox_unchecked, \ue640
		case t.State == tasklist.Cancelled:
			char = "󱋭" // nf-md-checkbox_blank_off_outline, \uf52f
		case t.State == tasklist.Done:
			char = "󰄵" // nf-md-checkbox_marked_outline, \udb80\udd35
		}
	}
	indent := strings.Repeat(" ", t.Depth*2)

	grey := color.New(color.Faint, color.FgWhite).SprintFunc()
	none := color.New().SprintFunc()

	colorize := none
	if t.State == tasklist.Done {
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
func (noteplan *Noteplan) BuildTask(marker string, depth int, state tasklist.TaskState, text string) Task {
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

	if state.NotUnknown() && category == Bullet {
		category = Todo
	}

	if category == Bullet {
		return Task{Category: category, Text: text, Depth: depth}
	}
	if state == tasklist.Unknown {
		state = tasklist.Open
	}

	return Task{Category: category, State: state, Text: text, Depth: depth}
}

func (noteplan *Noteplan) parseTasks(data []byte, doc ast.Node) []Task {
	var tasks []Task
	var depth = -1
	markerMap := make(map[int]string)
	stack := NewStack[Task]()

	ast.Walk(doc, func(node ast.Node, enter bool) (ast.WalkStatus, error) {
		// found a new (sub)list
		if n, ok := node.(*ast.List); ok {
			if enter {
				// increasing depth counter when we enter a new list
				depth = depth + 1
				// TODO make it a char
				markerMap[depth] = string((*n).Marker)
			} else {
				depth = depth - 1
			}
		}

		// WRONG: ot every ListItem is closed. We can hvae nested elements
		if _, ok := node.(*ast.ListItem); ok {
			if enter {
				stack.Push(Task{})
			} else {
				taskStub := stack.Pop()
				actualTask := noteplan.BuildTask(markerMap[depth], depth, taskStub.State, taskStub.Text)
				tasks = append(tasks, actualTask)
			}
		}

		if n, ok := node.(*ast.Text); ok && stack.Size() > 0 {
			// `* Blog <2024-02-10` results in two Text nodes
			if enter {
				taskStub := stack.Pop()
				text := string(n.Text(data))
				fmt.Println(text)
				taskStub.Text = taskStub.Text + text
				stack.Push(taskStub)
			}
		}

		if n, ok := node.(*tasklist.TaskCheckBox); ok && enter {
			taskStub := stack.Peek()
			taskStub.State = n.State
		}

		return ast.WalkContinue, nil
	})

	return tasks
}
