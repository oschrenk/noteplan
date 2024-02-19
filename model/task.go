package model

import (
	"fmt"
	"strings"

	"github.com/fatih/color"
)

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
