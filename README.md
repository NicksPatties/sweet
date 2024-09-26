# Sweet

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

## Building

## A prerelease version

```sh
go build -ldflags "-X github.com/NicksPatties/sweet/version.version=v0.0.2-`date -u +%Y%m%d%H%M%S`" .
```
- The version is hard coded for now, so I need to change that to actually look at the git tag.

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

By now you should have a new command that you can run and test like its own standalone application.

## Writing e2e tests

See `e2e/exercise_test.sh` and `e2e/e2e_test_template.sh` for examples.

