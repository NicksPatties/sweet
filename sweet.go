/*
sweet - The Software Engineer's Exercise in Typing.
Runs an interactive typing exercise.
Once complete, it displays statistics, including words per minute (WPM), accuracy, and number of mistakes.
*/

package main

import (
	"fmt"
	"github.com/NicksPatties/sweet/about"
	"github.com/NicksPatties/sweet/add"
	"github.com/NicksPatties/sweet/util"

	"github.com/NicksPatties/sweet/exercise"
	"github.com/NicksPatties/sweet/stats"
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:  "sweet",
	Long: "The Software Engineer Exercise for Typing.",
	Run: func(cmd *cobra.Command, args []string) {
		language, _ := cmd.Flags().GetString("language")
		configDir := util.GetConfigDirectory()
		// TODO I should pass in a flags struct for this command.
		exercise.Run(configDir, language)
	},
}

func init() {
	// Add language flag to root command only.
	// The flags for other commands will be defined in their respective modules.
	rootCmd.Flags().StringP("language", "l", "", "Language for the typing game")

	rootCmd.CompletionOptions.DisableDefaultCmd = true

	commands := []*cobra.Command{
		about.Command,
		stats.Command,
		add.Command,
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
