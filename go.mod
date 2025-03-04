module github.com/ryancurrah/gomodguard

go 1.23.0

require (
	github.com/Masterminds/semver/v3 v3.3.1
	github.com/mitchellh/go-homedir v1.1.0
	github.com/phayes/checkstyle v0.0.0-20170904204023-bfd46e6a821d
	github.com/stretchr/testify v1.10.0
	golang.org/x/mod v0.23.0
	gopkg.in/yaml.v3 v3.0.1
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
)

retract v1.2.1 // Originally tagged for commit hash that was subsequently removed, and replaced by another commit hash
