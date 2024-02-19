package internal

import (
	"database/sql"

	. "github.com/oschrenk/noteplan/model"

	_ "github.com/mattn/go-sqlite3"
	"howett.net/plist"
)

func (noteplan *Noteplan) fetch(iso string, entry string, failFast bool) (*TaskSummary, error) {
	db, err := sql.Open("sqlite3", noteplan.settings.NoteCachePath)
	if err != nil {
		return logThenEmptyOrErr(err, iso, failFast)
	}
	defer db.Close()

	const metadataQuery = "select content from metadata where filename = ? LIMIT 1"
	row := db.QueryRow(metadataQuery, entry)

	var content []byte
	err = row.Scan(&content)
	if err != nil {
		return logThenEmptyOrErr(err, iso, failFast)
	}

	// the plist structure is very flexible, and uses a heterogenous
	// array to store values in a compact way
	// we leave it, and cast types as needed
	data := make(map[string]interface{})
	if _, err := plist.Unmarshal(content, &data); err != nil {
		return logThenEmptyOrErr(err, iso, failFast)
	}

	// [
	//   "$null",
	//   {
	//     "$class": 13,
	//     "atTags": 8,
	//     "calendarItemIDs": 12,
	//     "datedTodos": 11,
	//     "filename": 2,
	//     "hashTags": 9,
	//     "linkedItems": 10,
	//     "numCancelledTodos": 4,
	//     "numClosedTodos": 5,
	//     "numCompletedReminders": 4,
	//     "numDoneNoteChecklists": 4,
	//     "numDoneNoteTodos": 4,
	//     "numDoneTodos": 5,
	//     "numEvents": 6,
	//     "numOpenNoteChecklists": 4,
	//     "numOpenNoteTodos": 4,
	//     "numOpenTodos": 3,
	//     "numReminders": 4,
	//     "numScheduledTodos": 4,
	//     "timeframe": 7
	//   },
	//   "20240105.md",
	//   4,
	//   0,
	//   2,
	//   8,
	//   "day",
	//   "",
	//   "#personal",
	//   "W10=",
	//   "W10=",
	//   "W10=",
	//   {
	//     "$classes": [
	//             "NotePlan.NoteModel",
	//             "NSObject"
	//     ],
	//     "$classname": "NotePlan.NoteModel"
	//   }
	// ]
	objects := data["$objects"].([]interface{})
	indexes := objects[1].(map[string]interface{})
	numOpenTodos := int(objects[indexes["numOpenTodos"].(plist.UID)].(uint64))
	numDoneTodos := int(objects[indexes["numDoneTodos"].(plist.UID)].(uint64))

	return NewTaskSummary(iso, numOpenTodos, numDoneTodos), nil
}
