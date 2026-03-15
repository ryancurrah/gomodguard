package gomodguard

import (
	"slices"
	"strings"
)

// Allowed is a list of modules and module
// prefixes that are allowed to be used.
type Allowed struct {
	Modules  []string `yaml:"modules"`
	Domains  []string `yaml:"domains"`
	Prefixes []string `yaml:"prefixes"`
}

// IsAllowedModule returns true if the given module
// name is in the allowed modules list.
func (a *Allowed) IsAllowedModule(moduleName string) bool {
	allowedModules := a.Modules

	for i := range allowedModules {
		if strings.TrimSpace(moduleName) == strings.TrimSpace(allowedModules[i]) {
			return true
		}
	}

	return false
}

// IsAllowedModulePrefix returns true if the given modules prefix is
// in the allowed module prefixes list.
func (a *Allowed) IsAllowedModulePrefix(moduleName string) bool {
	allowedPrefixes := make([]string, 0, len(a.Prefixes)+len(a.Domains))
	allowedPrefixes = append(allowedPrefixes, a.Prefixes...)
	allowedPrefixes = append(allowedPrefixes, a.Domains...)
	allowedPrefixes = slices.Compact(allowedPrefixes)

	for i := range allowedPrefixes {
		if strings.HasPrefix(strings.TrimSpace(strings.ToLower(moduleName)),
			strings.TrimSpace(strings.ToLower(allowedPrefixes[i]))) {
			return true
		}
	}

	return false
}
