package cmd

import (
	"encoding/json"
	"fmt"
	"github.com/ijt/go-anytime"
	"github.com/oschrenk/noteplan/noteplan"
	"github.com/spf13/cobra"
	"strings"
	"time"
)

func init() {
	rootCmd.AddCommand(todoCmd)
}

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Print the todo count",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		dayOnly, _ := cmd.Flags().GetBool("day-only")
		weekOnly, _ := cmd.Flags().GetBool("week-only")
		failFast, _ := cmd.Flags().GetBool("fail-fast")
		asJson, _ := cmd.Flags().GetBool("json")

		var dateTime time.Time
		if len(args) == 0 {
			dateTime = time.Now()
		} else {
			parsedDay, err := anytime.Parse(args[0], time.Now())
			if err != nil {
				fmt.Println("Invalid argument")
				return
			}
			dateTime = parsedDay
		}

		// collect data
		todosMap := make(map[string]*noteplan.Todos)
		if !weekOnly {
			dayTodos, _ := noteplan.Day(dateTime, failFast)
			todosMap["day"] = dayTodos
		}
		if !dayOnly {
			weekTodos, _ := noteplan.Week(dateTime, failFast)
			todosMap["week"] = weekTodos
		}

		// print as json
		if asJson {
			s, _ := json.MarshalIndent(todosMap, "", "  ")
			fmt.Print(string(s))

			// print as text
		} else {
			for mode, todos := range todosMap {
				titledMode := strings.Title(mode)
				iso := todos.Iso
				fmt.Println(fmt.Sprint(titledMode, ", ", iso, ", Open: ", todos.Open))
				fmt.Println(fmt.Sprint(titledMode, ", ", iso, ", Closed: ", todos.Closed))
			}
		}
	},
}

func init() {
	todoCmd.Flags().BoolP("day-only", "d", false, "Show count for the day only")
	todoCmd.Flags().BoolP("week-only", "w", false, "Show count for the week only")
	todoCmd.Flags().BoolP("fail-fast", "f", false, "Fail if entry not found. If false, return 0 counts.")
	todoCmd.Flags().BoolP("json", "j", false, "Output as json")
}
