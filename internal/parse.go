package internal

import (
	"bytes"
	"regexp"
	"strings"

	model "github.com/oschrenk/noteplan/model"

	"github.com/yuin/goldmark/ast"
)

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
// Additionally there is a trigger character for a "checklist" item
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
func (noteplan *Noteplan) parseTask(marker string, task string, depth int) model.Task {
	var category model.TaskCategory

	switch {
	case marker == "*":
		if noteplan.settings.IsAsteriskTodo {
			category = model.Todo
		} else {
			category = model.Bullet
		}
	case marker == "-":
		if noteplan.settings.IsDashTodo {
			category = model.Todo
		} else {
			category = model.Bullet
		}
	case marker == "+":
		category = model.Checklist
	}
	text := strings.TrimSpace(task)

	return noteplan.parseTaskState(category, text, depth)
}

func (noteplan *Noteplan) parseTaskState(category model.TaskCategory, text string, depth int) model.Task {
	var state model.TaskState

	stateRegEx := `^\[([\s-x])\].*`
	re := regexp.MustCompile(stateRegEx)
	matches := re.FindStringSubmatch(text)

	// if any state trigger is found we assume it's a Todo
	if len(matches) == 2 {
		// the user might have changed their settings in the past or might have
		// manually changed a bullet item into a task item by using task markers
		// if it's checklist we keep it
		if category == model.Bullet {
			category = model.Todo
		}
		stateChar := matches[1]
		text = strings.TrimSpace(text[3:])
		switch {
		case stateChar == model.Open.Trigger():
			state = model.Open
		case stateChar == model.Cancelled.Trigger():
			state = model.Cancelled
		case stateChar == model.Done.Trigger():
			state = model.Done
		}

		// otherwise assume Open
	} else {
		state = model.Open
	}
	// if bullet do not return TaskState
	if category == model.Bullet {
		return model.Task{Category: category, Text: text, Depth: depth}
	}

	return model.Task{Category: category, State: state, Text: text, Depth: depth}
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

func (noteplan *Noteplan) parseTasks(data []byte, doc ast.Node) []model.Task {
	var tasks []model.Task
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
