package gomodguard

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/mod/modfile"
	"golang.org/x/mod/module"
)

func TestIsModuleBlocked(t *testing.T) {
	var tests = []struct {
		testName         string
		processor        Processor
		testModule       string
		wantBlockReasons int
	}{
		{
			"previous version blocked",
			Processor{
				blockedModulesFromModFile: map[string][]string{
					"github.com/foo/bar": {blockReasonNotInAllowedList},
				},
			},
			"github.com/foo/bar/v2",
			0,
		},
		{
			"ensure modules with similar prefixes are not blocked",
			Processor{
				blockedModulesFromModFile: map[string][]string{
					"github.com/aws/aws-sdk-go": {blockReasonNotInAllowedList},
				},
			},
			"github.com/aws/aws-sdk-go-v2",
			0,
		},
		{
			"ensure the same module is blocked",
			Processor{
				blockedModulesFromModFile: map[string][]string{
					"github.com/foo/bar/v3": {blockReasonNotInAllowedList},
				},
			},
			"github.com/foo/bar/v3",
			1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			blockReasons := tt.processor.isBlockedPackageFromModFile(tt.testModule)
			if len(blockReasons) != tt.wantBlockReasons {
				t.Logf(
					"Testing %v, expected %v blockReasons, got %v blockReasons",
					tt.testModule,
					tt.wantBlockReasons,
					len(blockReasons),
				)
				t.Fail()
			}
		})
	}
}

func Test_packageInModule(t *testing.T) { //nolint:funlen
	type args struct {
		pkg string
		mod string
	}

	tests := []struct {
		name           string
		args           args
		wantPkgIsInMod bool
	}{
		{
			name: "package with path bar",
			args: args{
				pkg: "github.com/acme/foo/bar",
				mod: "github.com/acme/foo",
			},
			wantPkgIsInMod: true,
		},
		{
			name: "package with path bar/baz",
			args: args{
				pkg: "github.com/acme/foo/bar/baz",
				mod: "github.com/acme/foo",
			},
			wantPkgIsInMod: true,
		},
		{
			name: "aws v1 package should not match v2 module",
			args: args{
				pkg: "github.com/aws/aws-sdk-go",
				mod: "github.com/aws/aws-sdk-go-v2",
			},
			wantPkgIsInMod: false,
		},
		{
			name: "aws v1 package should not match v2 module with a package path",
			args: args{
				pkg: "github.com/aws/aws-sdk-go/foo",
				mod: "github.com/aws/aws-sdk-go-v2",
			},
			wantPkgIsInMod: false,
		},
		{
			name: "aws v2 package should match v2 module",
			args: args{
				pkg: "github.com/aws/aws-sdk-go-v2",
				mod: "github.com/aws/aws-sdk-go-v2",
			},
			wantPkgIsInMod: true,
		},
		{
			name: "aws v2 package should match v2 module with a package path",
			args: args{
				pkg: "github.com/aws/aws-sdk-go-v2/foo",
				mod: "github.com/aws/aws-sdk-go-v2",
			},
			wantPkgIsInMod: true,
		},
		{
			name: "package with different major version",
			args: args{
				pkg: "github.com/foo/bar/v20",
				mod: "github.com/foo/bar",
			},
			wantPkgIsInMod: false,
		},
		{
			name: "package with no major version",
			args: args{
				pkg: "github.com/foo/bar",
				mod: "github.com/foo/bar/v10",
			},
			wantPkgIsInMod: false,
		},
		{
			name: "with different major version",
			args: args{
				pkg: "github.com/foo/bar/v40",
				mod: "github.com/foo/bar/v41",
			},
			wantPkgIsInMod: false,
		},
		{
			name: "with same major version",
			args: args{
				pkg: "github.com/foo/bar/v50",
				mod: "github.com/foo/bar/v50",
			},
			wantPkgIsInMod: true,
		},
		{
			name: "same major version with path",
			args: args{
				pkg: "github.com/foo/bar/v60/baz/taz",
				mod: "github.com/foo/bar/v60",
			},
			wantPkgIsInMod: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gotPkgIsInMod := isPackageInModule(tt.args.pkg, tt.args.mod)
			if gotPkgIsInMod != tt.wantPkgIsInMod {
				t.Errorf("packageInModule() gotPkgIsInMod = %v, want %v", gotPkgIsInMod, tt.wantPkgIsInMod)
			}
		})
	}
}

func TestProcessorSetBlockedModulesWithAllowList(t *testing.T) {
	config := &Configuration{
		Allowed: Allowed{
			Modules: []string{
				"gopkg.in/yaml.v2",
				"github.com/go-xmlfmt/xmlfmt",
				"github.com/Masterminds/semver/v3",
				"github.com/ryancurrah/gomodguard",
			},
			Domains: []string{
				"golang.org",
			},
		},
		Blocked: Blocked{
			Modules: BlockedModules{
				{
					"github.com/uudashr/go-module": BlockedModule{
						Recommendations: []string{"golang.org/x/mod"},
						Reason:          "`mod` is the official go.mod parser library.",
					},
				},
				{
					"github.com/gofrs/uuid": BlockedModule{
						Recommendations: []string{"github.com/ryancurrah/gomodguard"},
						Reason:          "testing if module is not blocked when it is recommended.",
					},
				},
			},
			Versions: BlockedVersions{
				{
					"github.com/mitchellh/go-homedir": BlockedVersion{
						Version: "<= 1.1.0",
						Reason:  "testing if blocked version constraint works.",
					},
				},
			},
			LocalReplaceDirectives: true,
		},
	}

	modfile := &modfile.File{
		Module: &modfile.Module{
			Mod: module.Version{
				Path:    "github.com/ryancurrah/gomodguard",
				Version: "v1.0.0",
			},
		},
		Require: []*modfile.Require{
			{
				Mod: module.Version{
					Path:    "gopkg.in/yaml.v2",
					Version: "v2.4.0",
				},
			},
			{
				Mod: module.Version{
					Path:    "github.com/uudashr/go-module",
					Version: "v1.0.0",
				},
			},
			{
				Mod: module.Version{
					Path:    "github.com/gofrs/uuid",
					Version: "v1.2.3",
				},
			},
			{
				Mod: module.Version{
					Path:    "github.com/mitchellh/go-homedir",
					Version: "v1.0.0",
				},
			},
		},
		Replace: []*modfile.Replace{
			{
				Old: module.Version{
					Path:    "github.com/gofrs/uuid",
					Version: "v1.2.3",
				},
				New: module.Version{
					Path:    "/path/to/local/package",
					Version: "",
				},
			},
		},
	}

	processor := &Processor{
		Config:  config,
		Modfile: modfile,
	}

	processor.SetBlockedModules()

	// Assert number of blocked modules
	assert.Len(t, processor.blockedModulesFromModFile, 3)

	// Assert blocked modules
	assert.Equal(t,
		[]string{"import of package `%s` is blocked because the module is in the blocked modules list. " +
			"`golang.org/x/mod` is a recommended module. `mod` is the official go.mod parser library."},
		processor.blockedModulesFromModFile["github.com/uudashr/go-module"],
	)

	assert.Equal(t,
		[]string{"import of package `%s` is blocked because the module has a local replace directive."},
		processor.blockedModulesFromModFile["github.com/gofrs/uuid"],
	)

	assert.Equal(t,
		[]string{"import of package `%s` is blocked because the module is in the blocked modules list. " +
			"version `v1.0.0` is blocked because it does not meet the version constraint `<= 1.1.0`. testing " +
			"if blocked version constraint works."},
		processor.blockedModulesFromModFile["github.com/mitchellh/go-homedir"],
	)
}

func TestProcessorSetBlockedModulesWithoutAllowList(t *testing.T) {
	config := &Configuration{
		Blocked: Blocked{
			Modules: BlockedModules{
				{
					"gotest.tools/v3": BlockedModule{
						Recommendations: []string{"github.com/stretchr/testify/assert"},
						Reason:          "We have standardized on `github.com/stretchr/testify/assert`.",
					},
				},
			},
		},
	}

	modfile := &modfile.File{
		Module: &modfile.Module{
			Mod: module.Version{
				Path:    "github.com/ryancurrah/gomodguard",
				Version: "v1.0.0",
			},
		},
		Require: []*modfile.Require{
			{
				Mod: module.Version{
					Path:    "gotest.tools/v3",
					Version: "v3.4.0",
				},
			},
		},
	}

	processor := &Processor{
		Config:  config,
		Modfile: modfile,
	}

	processor.SetBlockedModules()

	// Assert number of blocked modules
	assert.Len(t, processor.blockedModulesFromModFile, 1)

	// Assert blocked modules
	assert.Equal(t,
		[]string{"import of package `%s` is blocked because the module is in the blocked modules list. " +
			"`github.com/stretchr/testify/assert` is a recommended module. We have standardized on `github.com/stretchr/testify/assert`."},
		processor.blockedModulesFromModFile["gotest.tools/v3"],
	)
}

func TestProcessorSetBlockedModulesWithInvalidVersionConstraint(t *testing.T) {
	config := &Configuration{
		Blocked: Blocked{
			Versions: BlockedVersions{
				{
					"github.com/gin-gonic/gin": BlockedVersion{
						Version: "== 1.0.0",
					},
				},
			},
		},
	}

	modfile := &modfile.File{
		Module: &modfile.Module{
			Mod: module.Version{
				Path:    "github.com/ryancurrah/gomodguard",
				Version: "v1.0.0",
			},
		},
		Require: []*modfile.Require{
			{
				Mod: module.Version{
					Path:    "github.com/gin-gonic/gin",
					Version: "v1.0.0",
				},
			},
		},
	}

	processor := &Processor{
		Config:  config,
		Modfile: modfile,
	}

	processor.SetBlockedModules()

	// Assert number of blocked modules
	assert.Len(t, processor.blockedModulesFromModFile, 1)

	// Assert blocked modules
	assert.Equal(t,
		[]string{"import of package `%s` is blocked because the version constraint is invalid. improper constraint: == 1.0.0"},
		processor.blockedModulesFromModFile["github.com/gin-gonic/gin"],
	)
}
