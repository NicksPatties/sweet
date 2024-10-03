# Sweet

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

## What is Sweet?

**Sweet** is a **S**oft**w**are **E**ngineering **E**xercise for **T**yping. In other words, it's a touch typing exercise command line interface specifically designed for programmers.

## Installation

### Using `go`

Assuming you have `go` installed, you can use the following command:

```sh
go install github.com/NicksPatties/sweet@latest
```

### Downloading an executable

1. Go to the [releases](https://github.com/NicksPatties/sweet/releases) page.

2. Download the executable with the matching operating system and architecture in the file name. The executables in the release are shown in the format `sweet-<os>-<arch>`.

3. Verify your downloaded executable is the same by comparing its checksum with the provided checksum in the release page. You can use a terminal command like this to generate the checksum for comparison.

    ```sh
    sha256sum <your_sweet_executable_file>
    ```
    If the checksums do not match, **do not run the file**. Please [report an issue](https://github.com/NicksPatties/sweet/issues) if something is wrong.

4. Move the executable to a directory that is included in your `$PATH` variable.

If all is well, then you're ready to use `sweet`!

## Usage

```
The Software Engineer Exercise for Typing.

Usage:
  sweet [flags]
  sweet [command]

Available Commands:
  about       Print details about the application
  add         Add an exercise
  help        Help about any command
  stats       Print statistics about typing exercises

Flags:
  -e, --end uint          Language for the typing game (default 18446744073709551615)
  -h, --help              help for sweet
  -l, --language string   Language for the typing game
  -s, --start uint        Language for the typing game

Use "sweet [command] --help" for more information about a command.
```

## Contributions

If you notice any bugs, or have general feedback regarding your experience using `sweet`, please post an [issue](https://github.com/NicksPatties/sweet/issues) in our GitHub repo. You may also email me at [nickspatties@proton.me](mailto:nickspatties@proton.me?subject=Sweet%20Issue%3A%20%3CYour%20issue%20title%20here%3E&body=Sweet%20version%3A%20%3Csweet%20version%3E%0D%0ADetails%3A%20%3Cadd%20details%20here%3E).

Wanna contribute a change to the code? Please fork the repository, and then submit a pull request!

## License

[MIT](LICENSE)

# Contributor instructions

## Running the application

```sh
go build .
./sweet
```

## Testing

### Running unit tests for a module

```sh
go test ./{{module-name}}
```

### Running unit tests for all modules

```sh
go test ./...
```

## Building a release version

```sh
# Assuming you're currently on the commit you'd like to release
git checkout main
git tag {{version}}
git push origin {{version}}
./release
```

## Creating a new command

1. Create new file called `{{command}}/{{command}}.go`
  - `package {{command}}`
  - `const CommandName = "{{command}}"`
2. Create a function in new file called `Run`
  - should return an int
  - should accept `[]string` as first parameter for args
  - should accept any other inputs it needs
3. (optional) Add `{{command}}Cmd` variable of type `*flags.FlagSet`
4. Create a test file `{{command}}/{{command}}_test.go`
5. In `sweet.go`, add func signature from step 2 to `Commands` struct
6. In `Run` in `sweet.go`, add a case to the `switch subCommand` statement
7. In `Main` in `sweet.go`, add the `Run` function from your new `{{command}}` module to the `defaultCommands` struct.
8. Valiate flags in the `&cobra.Command`'s struct, then pass valid params to the `Run` function

By now you should have a new command that you can run and test like its own standalone application.

## Writing e2e tests

See `e2e/exercise_test.sh` and `e2e/e2e_test_template.sh` for examples.

