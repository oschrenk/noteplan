package internal

import (
	"testing"
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
