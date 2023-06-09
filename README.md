```
Hey! That's

         ,gg,                                                  gg  
        i8""8i                                        I8      ,gg, 
        `8,,8'                                        I8      i88i 
         `88'                                      88888888   i88i 
         dP"8,                                        I8      i88i 
        dP' `8a  gg    gg    gg    ,ggg,    ,ggg,     I8      ,gg, 
       dP'   `Yb I8    I8    88bg i8" "8i  i8" "8i    I8       gg  
   _ ,dP'     I8 I8    I8    8I   I8, ,8I  I8, ,8I   ,I8,          
   "888,,____,dP,d8,  ,d8,  ,8I   `YbadP'  `YbadP'  ,d88b,     aa  
   a8P"Y88888P" P""Y88P""Y88P"   888P"Y888888P"Y88888P""Y88    88  
```

[![Go Reference](https://pkg.go.dev/badge/github.com/NicksPatties/sweet.svg)](https://pkg.go.dev/github.com/NicksPatties/sweet)
[![Go Report Card](https://goreportcard.com/badge/github.com/NicksPatties/sweet)](https://goreportcard.com/report/github.com/NicksPatties/sweet)


# What is Sweet?
**Sweet** is a **S**oft**w**are **E**ngineering **E**xercise for **T**yping. In other words, it's a touch typing exercise command line interface specifically designed for programmers.

# Installation
## Using `go`
Assuming you have `go` installed, you can use the following command:
```
$ go install github.com/NicksPatties/sweet@latest
```
Now you can use `sweet`!

## Downloading an executable
You can also look in the [releases](https://github.com/NicksPatties/sweet/releases) section for an executable which matches your operating system and system architecture.

1. Download the executable with the matching operating system and architecture in the file name. The executables in the release are shown in the format `sweet_<os>_<arch>`.

2. Verify your downloaded executable is the same by comparing its checksum with the provided checksum in the release page. You can use a terminal command like this to generate the checksum for comparison.
    ```
    $ shasum -a 256 <your_sweet_executable_file>
    ```
    If the checksums do not match, **do not run the file**. Please notify us by reporting an issue on our GitHub page.

3. If the executable looks good, then place the executable in a folder that is specified by your `$PATH` variable.

If all is well, then you're ready to use `sweet`!

# Usage
```
Usage:

	sweet [subcommand]

Subcommands

	help                            Opens this help menu
	add [path]                      Adds the file tat this path to the exercise list
	lang [go|js|ts|java...]	        Finds a random exercise with the specified extension
	list                            Lists the available exercises to run
	[exercise name]                 Runs this exercise
```

# Contributions
If you notice any bugs, or have general feedback regarding your experience using `sweet`, please post an [issue](https://github.com/NicksPatties/sweet/issues) in our GitHub repo.

Wanna contribute a change to the code? Please fork the repository, and then submit a pull request!

# License
[MIT](LICENSE)

