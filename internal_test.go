package gomodguard

import "testing"

func TestIsModuleBlocked(t *testing.T) {
	var tests = []struct {
		testName   string
		processor  Processor
		testModule string
	}{
		{
			"previous version blocked",
			Processor{
				blockedModulesFromModFile: map[string][]string{
					"github.com/foo/bar": {blockReasonNotInAllowedList},
				},
			},
			"github.com/foo/bar/v2",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			blockReasons := tt.processor.isBlockedPackageFromModFile(tt.testModule)
			if len(blockReasons) > 0 {
				t.Logf("Testing %v, expected allowed, was blocked: %v", tt.testModule, blockReasons)
				t.Fail()
			}
		})
	}
}
