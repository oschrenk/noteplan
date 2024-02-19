package cmd

import (
	"fmt"
	"time"

	"github.com/ijt/go-anytime"
	"github.com/spf13/cobra"

	np "github.com/oschrenk/noteplan/internal"
	model "github.com/oschrenk/noteplan/model"
)

func init() {
	rootCmd.AddCommand(todoCmd)
}

var todoCmd = &cobra.Command{
	Use:   "todo",
	Short: "Show todos",
	Args:  cobra.MaximumNArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		verbose, _ := cmd.Flags().GetBool("verbose")
		np.Logger.Enabled = verbose
		WithSummary, _ := cmd.Flags().GetBool("summary")
		ShowBullet, _ := cmd.Flags().GetBool("show-bullet")
		ShowCancelled, _ := cmd.Flags().GetBool("show-cancelled")
		ShowDone, _ := cmd.Flags().GetBool("show-done")

		now := time.Now()
		dateTime := now
		precision := np.Day
		if len(args) > 0 {
			rawDate := args[0]
			parsedRange, err := anytime.ParseRange(rawDate, now)
			if err != nil {
				np.Logger.Log(fmt.Sprintf("Failed parsing date \"%s\"", rawDate))
			}
			parsedDay := parsedRange.Time
			precision, err = np.BuildTimePrecision(parsedRange.Duration)
			if err != nil {
				np.Logger.Log(fmt.Sprintf("Can't parse precision \"%s\"", rawDate))
			}
			dateTime = parsedDay
		}

		noteplan := np.NewInstance()
		tasks, err := noteplan.GetTasks(dateTime, precision)
		open := 0
		if err == nil {
			for _, task := range tasks {
				switch task.State {
				case model.Cancelled:
					if ShowCancelled {
						fmt.Println(task.String())
					}
				case model.Done:
					if ShowDone {
						fmt.Println(task.String())
					}
				case model.Open:
					if (task.Category == model.Bullet && ShowBullet) ||
						(task.Category == model.Todo) || (task.Category == model.Checklist) {
						fmt.Println(task.String())
					}
					if task.Category != model.Bullet {
						open = open + 1
					}
				}
			}
		}

		if WithSummary {
			fmt.Println(open, "open tasks")
		}
	},
}

func init() {
	todoCmd.Flags().BoolP("verbose", "v", false, "Log verbose")
	todoCmd.Flags().BoolP("summary", "s", true, "Print summary")
	todoCmd.Flags().BoolP("show-cancelled", "c", true, "Show Cancelled")
	todoCmd.Flags().BoolP("show-done", "d", true, "Show Done")
	todoCmd.Flags().BoolP("show-bullet", "b", false, "Show Bullet")
}
