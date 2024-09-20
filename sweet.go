/*
sweet - The Software Engineer's Exercise in Typing.
Runs an interactive typing exercise.
Once complete, it displays statistics, including words per minute (WPM), accuracy, and number of mistakes.
*/

package main

import (
	"fmt"
	"github.com/NicksPatties/sweet/about"
	"github.com/NicksPatties/sweet/stats"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "sweet",
	Long: "The Software Engineer Exercise for Typing.",
	Run: func(cmd *cobra.Command, args []string) {
		language, _ := cmd.Flags().GetString("language")
		runTypingGame(language)
	},
}

func runTypingGame(language string) {
	fmt.Printf("Running typing game in %s language\n", language)
	// Implement your typing game logic here
}

func init() {
	// Add language flag to root command only
	rootCmd.Flags().StringP("language", "l", "", "Language for the typing game")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	commands := []*cobra.Command{
		about.Command,
		stats.Command,
	}

	for _, c := range commands {
		rootCmd.AddCommand(c)
	}
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
	}
}
