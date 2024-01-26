package noteplan

type TaskSummary struct {
	Iso    string `json:"iso"`
	Open   int    `json:"open"`
	Closed int    `json:"closed"`
}

func emptyTaskSummary(iso string) *TaskSummary {
	return &TaskSummary{Iso: iso, Open: 0, Closed: 0}
}

func newTaskSummary(iso string, open int, closed int) *TaskSummary {
	return &TaskSummary{Iso: iso, Open: open, Closed: closed}
}
