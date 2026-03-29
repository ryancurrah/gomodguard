package cli

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/Masterminds/semver/v3"
	"go.yaml.in/yaml/v4"

	"github.com/ryancurrah/gomodguard/v2"
)

type V1Allowed struct {
	Modules  []string `yaml:"modules"`
	Domains  []string `yaml:"domains"`
	Prefixes []string `yaml:"prefixes"`
}

type V1BlockedVersion struct {
	Version string `yaml:"version"`
	Reason  string `yaml:"reason"`
}

type V1BlockedModule struct {
	Recommendations []string `yaml:"recommendations"`
	Reason          string   `yaml:"reason"`
}

type V1Blocked struct {
	Modules                []map[string]V1BlockedModule  `yaml:"modules"`
	Versions               []map[string]V1BlockedVersion `yaml:"versions"`
	LocalReplaceDirectives bool                          `yaml:"local_replace_directives"`
}

type V1Configuration struct {
	Allowed V1Allowed `yaml:"allowed"`
	Blocked V1Blocked `yaml:"blocked"`
}

func MigrateConfig(filename string) int {
	b, err := os.ReadFile(filepath.Clean(filename))
	if err != nil {
		fmt.Fprintf(os.Stderr, "error reading %s: %v\n", filename, err)
		return 1
	}

	var v1 V1Configuration
	if err := yaml.Unmarshal(b, &v1); err != nil {
		fmt.Fprintf(os.Stderr, "error parsing %s as v1 config: %v\n", filename, err)
		return 1
	}

	v2 := gomodguard.Configuration{
		LocalReplaceDirectives: v1.Blocked.LocalReplaceDirectives,
	}

	for _, mod := range v1.Allowed.Modules {
		v2.Allowed = append(v2.Allowed, gomodguard.AllowedModule{
			Module:    mod,
			MatchType: gomodguard.ExactMatch,
		})
	}

	for _, pref := range v1.Allowed.Prefixes {
		v2.Allowed = append(v2.Allowed, gomodguard.AllowedModule{
			Module:    pref,
			MatchType: gomodguard.PrefixMatch,
		})
	}

	for _, dom := range v1.Allowed.Domains {
		v2.Allowed = append(v2.Allowed, gomodguard.AllowedModule{
			Module:    dom,
			MatchType: gomodguard.PrefixMatch,
		})
	}

	// Build a temporary map to merge modules and versions by module name.
	blockedIndex := make(map[string]int)

	for _, modMap := range v1.Blocked.Modules {
		for modName, bm := range modMap {
			v2.Blocked = append(v2.Blocked, gomodguard.BlockedModule{
				Module:          modName,
				MatchType:       gomodguard.ExactMatch,
				Recommendations: bm.Recommendations,
				Reason:          bm.Reason,
			})
			blockedIndex[modName] = len(v2.Blocked) - 1
		}
	}

	for _, versMap := range v1.Blocked.Versions {
		for modName, bv := range versMap {
			var constraint *semver.Constraints

			if bv.Version != "" {
				c, err := semver.NewConstraint(bv.Version)
				if err != nil {
					fmt.Fprintf(os.Stderr, "error parsing version constraint %q for %s: %v\n", bv.Version, modName, err)
					return 1
				}

				constraint = c
			}

			if idx, ok := blockedIndex[modName]; ok {
				// Merge version info into existing entry.
				entry := v2.Blocked[idx]
				entry.Version = constraint

				if bv.Reason != "" {
					entry.Reason = bv.Reason
				}

				v2.Blocked[idx] = entry
			} else {
				v2.Blocked = append(v2.Blocked, gomodguard.BlockedModule{
					Module:    modName,
					MatchType: gomodguard.ExactMatch,
					Version:   constraint,
					Reason:    bv.Reason,
				})
			}
		}
	}

	out, err := yaml.Marshal(v2)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error generating v2 config: %v\n", err)
		return 1
	}

	fmt.Printf("%s\n", out)

	return 0
}
