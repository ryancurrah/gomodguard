# gomodguard

Allow list linter for direct Go module dependencies. This is useful for large organizations where they want to standardize on the modules used and be able to recommend alternative modules.

## Description

Allowed modules are defined in a `.gomodguard.yaml` or `~/.gomodguard.yaml` YAML file. Modules can be allowed/permitted by module or domain name.

Any modules or domains not listed in the configuration are blocked.

The linter looks for blocked modules in `go.mod` and searches for imported packages where the packages module is blocked. Indirect modules are not considered.

Replacement modules can be suggested in the configuration.

Results are reported to `stdout` and `gomodguard-checkstyle.xml` which will allow the results to be imported into CI tools that read checkstyle format.

Logging statements are reported to `stderr`.

## Configuration

```yaml
allow:
  modules:
    - gopkg.in/yaml.v2
    - github.com/go-xmlfmt/xmlfmt
    - github.com/phayes/checkstyle
    - github.com/mitchellh/go-homedir
  domains:
    - golang.org

replacements:
  - modules: 
      - github.com/uudashr/go-module
    replacement: golang.org/x/mod
    reason: "`mod` is the official go.mod parser library."
```

## Example

```
╰─ ./gomodguard ./...

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

## License

**MIT**
