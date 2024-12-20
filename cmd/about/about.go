/*
about - Prints some details about sweet.
This shows the name of the application, some author contact details,
and repository information.
*/
package about

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

// used to fill in the version number
// i.e. go build -ldflags "-X github.com/NicksPatties/sweet/cmd/about.version=v0.1.0" .
var version string

var Cmd = &cobra.Command{
	Use:   "about",
	Short: "Print details about the application",
	Run: func(cmd *cobra.Command, args []string) {
		issueLink := "https://github.com/NicksPatties/sweet/issues"
		// This should be a support link on my personal website,
		// but for now this will do.
		supportLink := "https://liberapay.com/NicksPatties/"
		executableName := os.Args[0]
		if version == "" {
			version = "dev"
		}
		printAbout(version, issueLink, supportLink, executableName)
	},
}

func bold(text string) string {
	return fmt.Sprintf("\033[1m%s\033[0m", text)
}

func printAbout(version, issueLink, supportLink, executableName string) {
	tagline := fmt.Sprintf(
		"The %soft%sare %sngineer's %sxercise for %syping",
		bold("S"),
		bold("w"),
		bold("E"),
		bold("E"),
		bold("T"),
	)
	msg := `Hey! That's

      ,gg,                                                  gg  
     i8""8i                                        I8      ,gg, 
     ` + "`8,,8'" + `                                        I8      i88i 
      ` + "`88'" + `                                      88888888   i88i 
      dP"8,                                        I8      i88i 
     dP' ` + "`8a" + `  gg    gg    gg    ,ggg,    ,ggg,     I8      ,gg, 
    dP'   ` + "`Yb" + ` I8    I8    88bg i8" "8i  i8" "8i    I8       gg  
_ ,dP'     I8 I8    I8    8I   I8, ,8I  I8, ,8I   ,I8,          
"888,,____,dP,d8,  ,d8,  ,8I   ` + "`YbadP'" + `  ` + "`YbadP'" + `  ,d88b,     aa  
a8P"Y88888P" P""Y88P""Y88P"   888P"Y888888P"Y88888P""Y88    88  

%s 

Written by:
	NicksPatties

Version:
	%s

Having issues? Report them here:
	%s

Interested in supporting %s? You can do so here!
	%s

Copyright (c) 2023-2024 NicksPatties
`
	fmt.Printf(msg, tagline, version, issueLink, executableName, supportLink)
}
