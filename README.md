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

