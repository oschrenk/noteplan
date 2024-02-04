package cmd

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ijt/go-anytime"
	"github.com/spf13/cobra"

	np "github.com/oschrenk/noteplan/noteplan"
)

func init() {
	rootCmd.AddCommand(summaryCmd)
}

var summaryCmd = &cobra.Command{
	Use:   "summary",
	Short: "Print the task summmary for the day and week",
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

		noteplan := np.NewInstance()

		// collect data
		taskSummaries := make(map[string]*np.TaskSummary)
		if !weekOnly {
			daySummary, _ := noteplan.Day(dateTime, failFast)
			taskSummaries["day"] = daySummary
		}
		if !dayOnly {
			weekSummary, _ := noteplan.Week(dateTime, failFast)
			taskSummaries["week"] = weekSummary
		}

		// print as json
		if asJson {
			s, _ := json.MarshalIndent(taskSummaries, "", "  ")
			fmt.Print(string(s))

			// print as text
		} else {
			for mode, taskSummary := range taskSummaries {
				titledMode := strings.Title(mode)
				iso := taskSummary.Iso
				fmt.Println(fmt.Sprint(titledMode, ", ", iso, ", Open: ", taskSummary.Open))
				fmt.Println(fmt.Sprint(titledMode, ", ", iso, ", Closed: ", taskSummary.Closed))
			}
		}
	},
}

func init() {
	summaryCmd.Flags().BoolP("day-only", "d", false, "Show summary for the day only")
	summaryCmd.Flags().BoolP("week-only", "w", false, "Show summary for the week only")
	summaryCmd.Flags().BoolP("fail-fast", "f", false, "Fail if entry not found. If false, return 0 counts.")
	summaryCmd.Flags().BoolP("json", "j", false, "Output as json")
}
