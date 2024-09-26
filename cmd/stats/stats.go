package stats

import (
	"fmt"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "stats",
	Short: "Print statistics about typing exercises",
	Run: func(cmd *cobra.Command, args []string) {
		since, _ := cmd.Flags().GetString("since")
		printStats(since)
	},
}

func printStats(since string) {
	if since != "" {
		fmt.Printf("Printing stats since %s\n", since)
	} else {
		fmt.Println("Printing all-time stats")
	}
	// Implement your stats logic here
}

func init() {
	// Add since flag to stats command
	Command.Flags().StringP("since", "s", "", "Start date for statistics")
}
