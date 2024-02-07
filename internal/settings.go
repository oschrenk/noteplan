package internal

import (
	"log"
	"os"

	"howett.net/plist"
)

type Settings struct {
	CalendarDataPath string
	NoteCachePath    string
	IsAsteriskTodo   bool
	IsDashTodo       bool
	Extension        string `default:"md"`
}

func LoadSettings() Settings {

	GroupSettingsPath := os.ExpandEnv("$HOME/Library/Group Containers/group.co.noteplan.noteplan/Library/Preferences/group.co.noteplan.noteplan.plist")
	CalendarDataPath := os.ExpandEnv("$HOME/Library/Containers/co.noteplan.NotePlan3/Data/Library/Application Support/co.noteplan.NotePlan3/Calendar")
	NoteCachePath := os.ExpandEnv("$HOME/Library/Containers/co.noteplan.NotePlan3/Data/Library/Application Support/co.noteplan.NotePlan3/Caches/note-cache.db")

	file, err := os.ReadFile(GroupSettingsPath)
	if err != nil {
		log.Panicf("failed reading file: %s", err)
	}
	data := make(map[string]interface{})
	if _, err := plist.Unmarshal(file, &data); err != nil {
		log.Panicf("failed reading file: %s", err)
	}

	settings := Settings{
		CalendarDataPath: CalendarDataPath,
		NoteCachePath:    NoteCachePath,
		IsAsteriskTodo:   data["isAsteriskTodo"] == true,
		IsDashTodo:       data["isDashTodo"] == true,
		Extension:        "md",
	}

	return settings
}
