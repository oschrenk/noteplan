package internal

import (
	"log"

	. "github.com/oschrenk/noteplan/model"
)

func logThenEmptyOrErr(err error, iso string, failFast bool) (*TaskSummary, error) {
	if failFast {
		log.Fatal(err)
		return nil, err
	} else {
		return EmptyTaskSummary(iso), nil
	}
}
