module github.com/ryancurrah/gomodguard/cmd/gomodguard/v2

go 1.25.0

require (
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/mitchellh/go-homedir v1.1.0
	github.com/phayes/checkstyle v0.0.0-20170904204023-bfd46e6a821d
	github.com/ryancurrah/gomodguard/v2 v2.1.1
	github.com/stretchr/testify v1.12.0
	go.yaml.in/yaml/v4 v4.0.0-rc.4
)

require (
	golang.org/x/mod v0.36.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

retract (
	v2.0.1
	v2.0.0
)
