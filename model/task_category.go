package model

import ()

type TaskCategory int64

const (
	Bullet TaskCategory = iota
	Checklist
	Todo
)
