module github.com/ryancurrah/gomodguard/v2

go 1.26.0

require (
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/stretchr/testify v1.12.1
	golang.org/x/mod v0.41.0
)

require go.yaml.in/yaml/v3 v3.0.5 // indirect

retract (
	v2.0.1
	v2.0.0
)
