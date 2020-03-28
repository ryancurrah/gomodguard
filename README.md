# gomodguard

<img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/9afcc208898c763be95f046eb2f6080146607209.png" width="30%">

Allow list linter for direct Go module dependencies. This is useful for organizations where they want to standardize on the modules used and be able to recommend alternative modules.

## Description

Allowed modules are defined in a `.gomodguard.yaml` or `~/.gomodguard.yaml` file. Modules can be allowed/permitted by module or domain name.

Any modules or domains not listed in the configuration are blocked.

The linter looks for blocked modules in `go.mod` and searches for imported packages where the packages module is blocked. Indirect modules are not considered.

Replacement modules can be suggested in the configuration.

Results are reported to `stdout` and `gomodguard-checkstyle.xml` (If specified) which will allow the results to be imported into CI tools that read checkstyle format.

Logging statements are reported to `stderr`.

## Configuration

```yaml
allow:
  modules:                                                  # List of allowed modules
    - gopkg.in/yaml.v2
    - github.com/go-xmlfmt/xmlfmt
    - github.com/phayes/checkstyle
    - github.com/mitchellh/go-homedir
  domains:                                                  # List of allowed module domains
    - golang.org

replacements:
  - modules:                                                # List of modules that should be replaced
      - github.com/uudashr/go-module
    replacement: golang.org/x/mod                           # Module that should be used instead
    reason: "`mod` is the official go.mod parser library."  # Reason why the module should be used
```

## Usage

```
╰─ ./gomodguard -h
Usage: gomodguard <file> [files...]
Also supports package syntax but will use it in relative path, i.e. ./pkg/...
Flags:
  -f string
        Report results to the specified file. A report type must also be specified
  -file string

  -h    Show this help text
  -help

  -n    Don't lint test files
  -no-test

  -r string
        Report results to one of the following formats: checkstyle. A report file destination must also be specified
  -report string
```

## Example

```
╰─ ./gomodguard -r checkstyle -f gomodguard-checkstyle.xml ./...

info: allowed modules, [gopkg.in/yaml.v2 github.com/go-xmlfmt/xmlfmt github.com/phayes/checkstyle github.com/mitchellh/go-homedir]
info: allowed module domains, [golang.org]
info: go.mod file has '1' blocked module(s), [github.com/uudashr/go-module]
example/blocked_example.go:6: import of package `github.com/uudashr/go-module` is blocked because the module is not in the allowed modules list. `golang.org/x/mod` should be used instead. reason: `mod` is the official go.mod parser library.
```

Resulting checkstyle file

```
╰─ cat gomodguard-checkstyle.xml

<checkstyle version="1.0.0">
  <file name="example/blocked_example.go">
    <error line="6" column="1" severity="error" message="import" source="import of package `github.com/uudashr/go-module` is blocked because the module is not in the allowed modules list. `golang.org/x/mod` should be used instead. reason: `mod` is the official go.mod parser library.">
    </error>
  </file>
</checkstyle>
```

## Install

```
go get -u github.com/ryancurrah/gomodguard/cmd/gomodguard
```

## Develop

```
git clone https://github.com/ryancurrah/gomodguard.git && cd gomodguard

go build -o gomodguard cmd/gomodguard/main.go
```

## License

**MIT**
