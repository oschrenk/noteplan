package internal

import (
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"time"

	. "github.com/oschrenk/noteplan/model"

	"github.com/yuin/goldmark/ast"
)

type Noteplan struct {
	settings Settings
}

func NewInstance() Noteplan {
	return Noteplan{settings: LoadSettings()}
}

func (noteplan *Noteplan) GetTasks(dateTime time.Time, tp TimePrecision) ([]Task, error) {
	entry := ""
	switch tp {
	case Day:
		entry = fmt.Sprint(dateTime.Format("20060102"), ".", noteplan.settings.Extension)
	case Week:
		year, week := dateTime.ISOWeek()
		entry = fmt.Sprint(year, "-W", fmt.Sprintf("%02d", week), ".", noteplan.settings.Extension)
		fmt.Println(entry)
	default:
		return nil, fmt.Errorf("unsupported precision %s", tp)
	}
	path := noteplan.settings.CalendarDataPath + "/" + entry
	data, doc, err := parseMarkdown(path)
	if err != nil {
		return nil, err
	}

	tasks := noteplan.parseTasks(data, doc)

	return tasks, nil
}

func (noteplan *Noteplan) Day(dateTime time.Time, failFast bool) (*TaskSummary, error) {
	iso := dateTime.Format("2006-01-02")
	entry := fmt.Sprint(dateTime.Format("20060102"), ".", noteplan.settings.Extension)

	return noteplan.fetch(iso, entry, failFast)
}

func (noteplan *Noteplan) Week(dateTime time.Time, failFast bool) (*TaskSummary, error) {
	year, week := dateTime.ISOWeek()
	iso := fmt.Sprint(year, "-W", fmt.Sprintf("%02d", week))
	entry := fmt.Sprint(iso, ".", noteplan.settings.Extension)

	return noteplan.fetch(iso, entry, failFast)
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
