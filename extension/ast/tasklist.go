package ast

import (
	"fmt"
	gast "github.com/yuin/goldmark/ast"
)

type TaskState int64

const (
	Open = iota
	Cancelled
	Done
	Incomplete
	Forwarded
	Scheduling
	Question
	Important
	Star
	Quote
	Location
	Bookmark
)

func NewTaskState(b byte) TaskState {
	return Open
}

func (s TaskState) Char() string {
	switch s {
	case Open:
		return " "
	case Cancelled:
		return "-"
	case Done:
		return "x"
	case Forwarded:
		return ">"
	}

	// TOOD panic
	return "u"
}

// A TaskCheckBox struct represents a checkbox of a task list.
type TaskCheckBox struct {
	gast.BaseInline
	State TaskState
}

// Dump implements Node.Dump.
func (n *TaskCheckBox) Dump(source []byte, level int) {
	m := map[string]string{
		"State": fmt.Sprintf("%v", n.State.Char()),
	}
	gast.DumpHelper(n, source, level, m, nil)
}

// KindTaskCheckBox is a NodeKind of the TaskCheckBox node.
var KindTaskCheckBox = gast.NewNodeKind("TaskCheckBox")

// Kind implements Node.Kind.
func (n *TaskCheckBox) Kind() gast.NodeKind {
	return KindTaskCheckBox
}

// NewTaskCheckBox returns a new TaskCheckBox node.
func NewTaskCheckBox(state TaskState) *TaskCheckBox {
	return &TaskCheckBox{
		State: state,
	}
}
