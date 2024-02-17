package tasklist

import (
	"regexp"
)

type TaskState int64

const (
	Unknown = iota
	Open
	Cancelled
	Done
	Forwarded
)

var TaskStateRegexp = regexp.MustCompile(`^\[([\sx\->])\]\s*`)

func NewTaskState(b byte) TaskState {
	switch b {
	case ' ':
		return Open
	case '-':
		return Cancelled
	case 'x':
		return Done
	case '>':
		return Forwarded
	}

	return Unknown
}

func (s TaskState) Trigger() byte {
	switch s {
	case Open:
		return ' '
	case Cancelled:
		return '-'
	case Done:
		return 'x'
	case Forwarded:
		return '>'
	}

	return '?'
}

func (s TaskState) String() string {
	return string(s.Trigger())
}

func (s TaskState) NotUnknown() bool {
	return s != Unknown
}
