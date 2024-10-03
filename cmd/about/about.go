/*
about - Prints some details about sweet.
This shows the name of the application, some author contact details,
and repository information.
*/
package about

import (
	"fmt"

	"github.com/spf13/cobra"
)

// used to fill in the version number
// i.e. go build -ldflags "-X github.com/NicksPatties/sweet/cmd/about.version=v0.1.0" .

var version string

var Command = &cobra.Command{
	Use:   "about",
	Short: "Print details about the application",
	Run: func(cmd *cobra.Command, args []string) {
		printAbout(version, "issueLink", "supportLink", "executableName")
	},
}

func printAbout(version, issueLink, supportLink, executableName string) {
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

The Software Engineering Exercise for Typing 

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
	fmt.Printf(msg, version, issueLink, executableName, supportLink)
}
