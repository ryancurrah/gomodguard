module github.com/ryancurrah/gomodguard/v2

go 1.25.0

require (
	github.com/Masterminds/semver/v3 v3.5.0
	github.com/stretchr/testify v1.12.0
	golang.org/x/mod v0.36.0
)

require gopkg.in/yaml.v3 v3.0.1 // indirect

retract (
	v2.0.1
	v2.0.0
)
