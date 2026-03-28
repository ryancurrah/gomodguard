# gomodguard
[![License](https://img.shields.io/github/license/ryancurrah/gomodguard?style=flat-square)](/LICENSE)
[![Codecov](https://img.shields.io/codecov/c/gh/ryancurrah/gomodguard?style=flat-square)](https://codecov.io/gh/ryancurrah/gomodguard)
[![GitHub Workflow Status](https://img.shields.io/github/actions/workflow/status/ryancurrah/gomodguard/go.yml?branch=main&logo=Go&style=flat-square)](https://github.com/ryancurrah/gomodguard/actions?query=workflow%3AGo)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/ryancurrah/gomodguard?style=flat-square)](https://github.com/ryancurrah/gomodguard/releases/latest)
[![Docker](https://img.shields.io/docker/pulls/ryancurrah/gomodguard?style=flat-square)](https://hub.docker.com/r/ryancurrah/gomodguard)
[![Github Releases Stats of golangci-lint](https://img.shields.io/github/downloads/ryancurrah/gomodguard/total.svg?logo=github&style=flat-square)](https://somsubhra.com/github-release-stats/?username=ryancurrah&repository=gomodguard)

<img src="https://storage.googleapis.com/gopherizeme.appspot.com/gophers/9afcc208898c763be95f046eb2f6080146607209.png" width="30%">

Allow and block list linter for direct Go module dependencies. This is useful for organizations where they want to standardize on the modules used and be able to recommend alternative modules.

## Description

Allowed and blocked modules are defined in a `./.gomodguard.yaml` or `~/.gomodguard.yaml` file. 

Modules can be allowed by module or prefix name. When allowed modules are specified any modules not in the allowed configuration are blocked.

If no allowed modules or module prefixes are specified then all modules are allowed except for blocked ones.

The linter looks for blocked modules in `go.mod` and searches for imported packages where the imported packages module is blocked. Indirect modules are not considered.

Alternative modules can be optionally recommended in the blocked modules list.

If the linted module imports a blocked module but the linted module is in the recommended modules list the blocked module is ignored. Usually, this means the linted module wraps that blocked module for use by other modules, therefore the import of the blocked module should not be blocked.

Version constraints can be specified for modules as well which lets you block new or old versions of modules or specific versions.

When multiple rules can match the same module (e.g., overlapping exact, prefix, and regex rules), they are evaluated using a layered strategy for deterministic results:

1. **Exact match** — highest priority; wins over prefix and regex.
2. **Prefix match** — next priority; longest matching prefix wins.
3. **Regex match** — lowest priority; evaluated in alphabetical key order; first match wins.

Results are printed to `stdout`.

Logging statements are printed to `stderr`.

Results can be exported to different report formats. Which can be imported into CI tools. See the help section for more information.

# Configuration

```yaml
allowed:
  go.yaml.in/yaml/v4:
  github.com/go-xmlfmt/xmlfmt:
  github.com/confluentinc/confluent-kafka-go/v2:
    version: "== 2.5.0"
  github.com/kubernetes:
    match_type: prefix
  github.com/apache/arrow-go:
    match_type: prefix
  "github.com/somecompany/.*":
    match_type: regex

blocked:
  github.com/uudashr/go-module:
    match_type: exact # or regex, prefix
    recommendations:
      - golang.org/x/mod
    reason: "`mod` is the official go.mod parser library."
  github.com/mitchellh/go-homedir:
    version: "<= 1.1.0"
    reason: "testing if blocked version constraint works."
  "github.com/badcompany/.*":
    match_type: regex
    reason: "No badcompany packages are permitted."
```

## Example .gomodguard.yaml Files

The following example configuration files are available:

- [examples/alloptions/.gomodguard.yaml](examples/alloptions/.gomodguard.yaml)
- [examples/allowedversion/.gomodguard.yaml](examples/allowedversion/.gomodguard.yaml)
- [examples/emptyallowlist/.gomodguard.yaml](examples/emptyallowlist/.gomodguard.yaml)
- [examples/indirectdep/.gomodguard.yaml](examples/indirectdep/.gomodguard.yaml)
- [examples/majorversion/.gomodguard.yaml](examples/majorversion/.gomodguard.yaml)
- [examples/regexversion/.gomodguard.yaml](examples/regexversion/.gomodguard.yaml)
- [examples/regextest/.gomodguard.yaml](examples/regextest/.gomodguard.yaml)

### Migrating from v1

If you have a v1 `.gomodguard.yaml` file, you can automatically migrate it to the new v2 schema by running:

```
gomodguard migrate > .gomodguard-v2.yaml
mv .gomodguard-v2.yaml .gomodguard.yaml
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

  -h	Show this help text
  -help

  -i int
    	Exit code when issues were found (default 2)
  -issues-exit-code int 
      (default 2)
  
  -n	Don't lint test files
  -no-test

  -r string
    	Report results to one of the following formats: checkstyle. A report file destination must also be specified
  -report string
```

## Example

```
╰─ cd examples/alloptions
╰─ gomodguard -r checkstyle -f gomodguard-checkstyle.xml ./...

info: allowed modules, [github.com/Masterminds/semver/v3 github.com/go-xmlfmt/xmlfmt golang.org gopkg.in/yaml.v3]
info: blocked modules, [github.com/gofrs/uuid github.com/mitchellh/go-homedir github.com/uudashr/go-module]
blocked_example.go:6:1 import of package `github.com/gofrs/uuid` is blocked because the module is in the blocked modules list. `github.com/ryancurrah/gomodguard` is a recommended module. testing if module is not blocked when it is recommended.
blocked_example.go:7:1 import of package `github.com/mitchellh/go-homedir` is blocked because the module is in the blocked modules list. version `v1.1.0` is blocked because it does not meet the version constraint `<=1.1.0`. testing if blocked version constraint works.
blocked_example.go:8:1 import of package `github.com/uudashr/go-module` is blocked because the module is in the blocked modules list. `golang.org/x/mod` is a recommended module. `mod` is the official go.mod parser library.
```

Resulting checkstyle file

```
╰─ cat gomodguard-checkstyle.xml

<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="1.0.0">
  <file name="blocked_example.go">
    <error line="6" column="1" severity="error" message="import of package `github.com/gofrs/uuid` is blocked because the module is in the blocked modules list. `github.com/ryancurrah/gomodguard` is a recommended module. testing if module is not blocked when it is recommended." source="gomodguard"></error>
    <error line="7" column="1" severity="error" message="import of package `github.com/mitchellh/go-homedir` is blocked because the module is in the blocked modules list. version `v1.1.0` is blocked because it does not meet the version constraint `&lt;=1.1.0`. testing if blocked version constraint works." source="gomodguard"></error>
    <error line="8" column="1" severity="error" message="import of package `github.com/uudashr/go-module` is blocked because the module is in the blocked modules list. `golang.org/x/mod` is a recommended module. `mod` is the official go.mod parser library." source="gomodguard"></error>
  </file>
</checkstyle>
```

## Install

```
go install github.com/ryancurrah/gomodguard/v2/cmd/gomodguard@latest
```

## Develop

```
git clone https://github.com/ryancurrah/gomodguard.git && cd gomodguard/cmd/gomodguard

go build -o gomodguard main.go
```

## License

**MIT**
