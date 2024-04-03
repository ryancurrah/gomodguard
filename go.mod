module github.com/ryancurrah/gomodguard

go 1.21

require (
	github.com/Masterminds/semver/v3 v3.2.1
	github.com/go-xmlfmt/xmlfmt v1.1.2
	github.com/mitchellh/go-homedir v1.1.0
	github.com/phayes/checkstyle v0.0.0-20170904204023-bfd46e6a821d
	github.com/stretchr/testify v1.9.0
	github.com/t-yuki/gocover-cobertura v0.0.0-20180217150009-aaee18c8195c
	golang.org/x/mod v0.16.0
	gopkg.in/yaml.v2 v2.4.0
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract v1.2.1 // Originally tagged for commit hash that was subsequently removed, and replaced by another commit hash
