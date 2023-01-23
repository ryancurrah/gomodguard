package gomodguard

import (
	"testing"
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
