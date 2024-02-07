package internal

import (
	"log"
)

func logThenEmptyOrErr(err error, iso string, failFast bool) (*TaskSummary, error) {
	if failFast {
		log.Fatal(err)
		return nil, err
	} else {
		return emptyTaskSummary(iso), nil
	}
}
