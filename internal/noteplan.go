package internal

import (
	"fmt"
	"time"

	model "github.com/oschrenk/noteplan/model"
)

type Noteplan struct {
	settings Settings
}

func NewInstance() Noteplan {
	return Noteplan{settings: LoadSettings()}
}

func (noteplan *Noteplan) GetTasks(dateTime time.Time, tp TimePrecision) ([]model.Task, error) {
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
	data, doc, err := parseFile(path)
	if err != nil {
		return nil, err
	}

	tasks := noteplan.parseTasks(data, doc)

	return tasks, nil
}

func (noteplan *Noteplan) Day(dateTime time.Time, failFast bool) (*model.TaskSummary, error) {
	iso := dateTime.Format("2006-01-02")
	entry := fmt.Sprint(dateTime.Format("20060102"), ".", noteplan.settings.Extension)

	return noteplan.fetch(iso, entry, failFast)
}

func (noteplan *Noteplan) Week(dateTime time.Time, failFast bool) (*model.TaskSummary, error) {
	year, week := dateTime.ISOWeek()
	iso := fmt.Sprint(year, "-W", fmt.Sprintf("%02d", week))
	entry := fmt.Sprint(iso, ".", noteplan.settings.Extension)

	return noteplan.fetch(iso, entry, failFast)
}
