package tasklist

import (
	"fmt"
	gast "github.com/yuin/goldmark/ast"
)

// A TaskCheckBox struct represents a checkbox of a task list.
type TaskCheckBox struct {
	gast.BaseInline
	State TaskState
}

// Dump implements Node.Dump.
func (n *TaskCheckBox) Dump(source []byte, level int) {
	m := map[string]string{
		"State": fmt.Sprintf("%v", n.State),
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
