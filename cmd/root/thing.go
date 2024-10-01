package root

import (
	"fmt"
	"math"

	"github.com/spf13/cobra"
)

func createMockCommand() (mockCmd *cobra.Command) {
	mockCmd = &cobra.Command{
		Use:  "fake",
		Long: "fake command",
		Args: cobra.MaximumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("in execute\n")
			fmt.Printf("args %s\n", args)
		},
	}
	// taken from sweet.go:51
	mockCmd.Flags().StringP("language", "l", "", "Language for the typing game")
	mockCmd.Flags().UintP("start", "s", 0, "Language for the typing game")
	mockCmd.Flags().UintP("end", "e", math.MaxUint, "Language for the typing game")
	return
}

// func main() {

// 	mockCmd := createMockCommand()
// 	args := []string{"-s", "10", "-l", "js", "filename"}
// 	mockCmd.SetArgs(args)
// 	mockCmd.Execute()

// 	language, _ := mockCmd.Flags().GetString("language")
// 	start, _ := mockCmd.Flags().GetUint("start")
// 	end, _ := mockCmd.Flags().GetUint("end")

// 	fmt.Printf("language %s\n", language)
// 	fmt.Printf("start    %d\n", start)
// 	fmt.Printf("end      %d\n", end)
// 	fmt.Printf("args     %s\n", args)

// }
