package cmd

import (
	"fmt"
	"time"

	"github.com/ijt/go-anytime"
	"github.com/spf13/cobra"

	np "github.com/oschrenk/noteplan/noteplan"
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
		if len(args) > 0 {
			rawDate := args[0]
			parsedDay, err := anytime.Parse(rawDate, now)
			if err != nil {
				np.Logger.Log(fmt.Sprintf("Failed parsing date \"%s\"", rawDate))
			}
			dateTime = parsedDay
		}

		noteplan := np.NewInstance()
		tasks, err := noteplan.GetTasks(dateTime)
		open := 0
		if err == nil {
			for _, task := range tasks {
				switch task.State {
				case np.Cancelled:
					if ShowCancelled {
						fmt.Println(task.String())
					}
				case np.Done:
					if ShowDone {
						fmt.Println(task.String())
					}
				case np.Open:
					if (task.Category == np.Bullet && ShowBullet) ||
						(task.Category == np.Todo) || (task.Category == np.Checklist) {
						fmt.Println(task.String())
					}
					if task.Category != np.Bullet {
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
