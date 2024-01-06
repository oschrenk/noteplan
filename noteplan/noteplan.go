package noteplan

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	"howett.net/plist"
	"log"
	"os"
	"time"
)

type Todos struct {
	Iso    string `json:"iso"`
	Open   int    `json:"open"`
	Closed int    `json:"closed"`
}

const Extension = "md"
const SqlitePath = "Library/Containers/co.noteplan.NotePlan3/Data/Library/Application Support/co.noteplan.NotePlan3/Caches/note-cache.db"
const Query = "select content from metadata where filename = ? LIMIT 1"

func emptyTodos(iso string) *Todos {
	return &Todos{Iso: iso, Open: 0, Closed: 0}
}

func newTodos(iso string, open int, closed int) *Todos {
	return &Todos{Iso: iso, Open: open, Closed: closed}
}

func logThenEmptyOrErr(err error, iso string, failFast bool) (*Todos, error) {
	if failFast {
		log.Fatal(err)
		return nil, err
	} else {
		return emptyTodos(iso), nil
	}
}

func fetch(iso string, entry string, failFast bool) (*Todos, error) {
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
		return nil, err
	}
	noteCachePath := fmt.Sprint(home, "/", SqlitePath)

	db, err := sql.Open("sqlite3", noteCachePath)
	if err != nil {
		return logThenEmptyOrErr(err, iso, failFast)
	}
	defer db.Close()

	row := db.QueryRow(Query, entry)

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

	return newTodos(iso, numOpenTodos, numDoneTodos), nil
}

func Day(dateTime time.Time, failFast bool) (*Todos, error) {
	iso := dateTime.Format("2006-01-02")
	entry := fmt.Sprint(dateTime.Format("20060102"), ".", Extension)
	return fetch(iso, entry, failFast)
}

func Week(dateTime time.Time, failFast bool) (*Todos, error) {
	year, week := dateTime.ISOWeek()
	iso := fmt.Sprint(year, "-W", fmt.Sprintf("%02d", week))
	entry := fmt.Sprint(iso, ".", Extension)
	return fetch(iso, entry, failFast)
}
