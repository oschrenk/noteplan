package model

import (
	"log"
)

type TaskState int64

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
