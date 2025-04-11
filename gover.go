package gomodguard

import "github.com/Masterminds/semver/v3"

func IsOkGoVersion(c, goVersion string) (bool, error) {
	if goVersion == "" {
		return true, nil
	}

	constraint, err := semver.NewConstraint(c)
	if err != nil {
		return false, err
	}

	version, err := semver.NewVersion(goVersion)
	if err != nil {
		return false, err
	}

	return constraint.Check(version), nil
}
