package gomodguard_test

import (
	"reflect"
	"testing"

	"github.com/ryancurrah/gomodguard"
)

func TestAllowedIsAllowedModule(t *testing.T) {
	var tests = []struct {
		testName            string
		allowedModules      gomodguard.Allowed
		lintedModuleName    string
		wantIsAllowedModule bool
	}{
		{
			"module is allowed",
			gomodguard.Allowed{Modules: []string{"github.com/someallowed/module"}},
			"github.com/someallowed/module",
			true,
		},
		{
			"module not allowed",
			gomodguard.Allowed{},
			"github.com/someblocked/module",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			isAllowedModule := tt.allowedModules.IsAllowedModule(tt.lintedModuleName)
			if !reflect.DeepEqual(isAllowedModule, tt.wantIsAllowedModule) {
				t.Errorf("got '%v' want '%v'", isAllowedModule, tt.wantIsAllowedModule)
			}
		})
	}
}

func TestAllowedIsAllowedModulePrefix(t *testing.T) {
	var tests = []struct {
		testName                  string
		allowedModules            gomodguard.Allowed
		lintedModuleName          string
		wantIsAllowedModulePrefix bool
	}{
		{
			"module is allowed",
			gomodguard.Allowed{Prefixes: []string{"github.com"}},
			"github.com/someallowed/module",
			true,
		},
		{
			"module not allowed",
			gomodguard.Allowed{},
			"github.com/someblocked/module",
			false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			isAllowedModulePrefix := tt.allowedModules.IsAllowedModulePrefix(tt.lintedModuleName)
			if !reflect.DeepEqual(isAllowedModulePrefix, tt.wantIsAllowedModulePrefix) {
				t.Errorf("got '%v' want '%v'", isAllowedModulePrefix, tt.wantIsAllowedModulePrefix)
			}
		})
	}
}
