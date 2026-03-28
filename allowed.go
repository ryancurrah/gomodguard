package gomodguard

import (
	"fmt"

	"github.com/Masterminds/semver/v3"
)

// Allowed is a map of modules or prefixes that are allowed to be used.
type Allowed map[string]AllowedRule

// AllowedRule defines the options for an allowed module.
type AllowedRule struct {
	MatchType MatchType           `yaml:"match_type"`
	Version   *semver.Constraints `yaml:"version"`
	Matcher   Matcher             `yaml:"-"`
}

// CheckVersion returns true if the module version matches the allowed constraint,
// or if no version constraint is specified.
func (r *AllowedRule) CheckVersion(moduleVersion string) (bool, error) {
	if r.Version == nil {
		return true, nil
	}

	version, err := semver.NewVersion(moduleVersion)
	if err != nil {
		return false, err
	}

	return r.Version.Check(version), nil
}

// NotAllowedReason returns the reason why the module version is not allowed.
func (r *AllowedRule) NotAllowedReason(moduleVersion string) string {
	if r == nil || r.Version == nil {
		return "the module is not in the allowed modules list."
	}

	return fmt.Sprintf("version `%s` does not meet the allowed version constraint `%s`.", moduleVersion, r.Version)
}

func (r *AllowedRule) ruleMatchType() MatchType {
	return r.MatchType
}

func (r *AllowedRule) ruleMatcher() Matcher { //nolint:ireturn
	return r.Matcher
}
