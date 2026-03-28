module github.com/ryancurrah/gomodguard/cmd/gomodguard/v2

go 1.25.0

replace github.com/ryancurrah/gomodguard/v2 => ../../

require (
	github.com/Masterminds/semver/v3 v3.4.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/phayes/checkstyle v0.0.0-20170904204023-bfd46e6a821d
	github.com/ryancurrah/gomodguard/v2 v2.0.0-00010101000000-000000000000
	github.com/stretchr/testify v1.11.1
	go.yaml.in/yaml/v4 v4.0.0-rc.4
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/mod v0.34.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
