package internal

import (
	"testing"

	model "github.com/oschrenk/noteplan/model"
)

var defaultSettings = Settings{
	// these paths should not be used during testing
	// CalendarDataPath: ...,
	// NoteCachePath:    ...,
	// Extension:      "md",
	IsAsteriskTodo: true,
	IsDashTodo:     false,
}

func TestParseEmptyString(t *testing.T) {
	md := []byte(``)
	data, doc, _ := parseString(md)
	noteplan := NewInstance(defaultSettings)
	tasks := noteplan.parseTasks(data, doc)

	if len(tasks) != 0 {
		t.Fatalf(`len(tasks) should be 0, but was %d`, len(tasks))
	}
}

func TestParseString(t *testing.T) {
	md := []byte(`# Heading`)
	data, doc, _ := parseString(md)
	noteplan := NewInstance(defaultSettings)
	tasks := noteplan.parseTasks(data, doc)

	if len(tasks) != 0 {
		t.Fatalf(`len(tasks) should be 0, but was %d`, len(tasks))
	}
}

func TestParseAsteriskTask(t *testing.T) {
	md := []byte(`* foo`)
	want := model.Task{
		Category: model.Todo,
		State:    model.Open,
		Text:     "foo",
		Depth:    0,
	}
	data, doc, _ := parseString(md)
	noteplan := NewInstance(defaultSettings)
	tasks := noteplan.parseTasks(data, doc)

	if len(tasks) != 1 {
		t.Fatalf(`len(tasks) should be 1, but was %d`, len(tasks))
	}

	if tasks[0] != want {
		t.Fatalf(`Wanted %v but got %v`, want, tasks[0])
	}
}
