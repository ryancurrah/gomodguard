package gomodguard_test

import (
	"strings"
	"testing"

	"github.com/ryancurrah/gomodguard"
)

func TestResultString(t *testing.T) {
	var tests = []struct {
		testName   string
		result     gomodguard.Issue
		wantString string
	}{
		{
			"reason lint failed",
			gomodguard.Issue{FileName: "test.go", LineNumber: 1, Reason: "Some reason."},
			"test.go:1:1 Some reason.",
		},
	}

	for _, tt := range tests {
		t.Run(tt.testName, func(t *testing.T) {
			result := tt.result.String()
			if !strings.EqualFold(result, tt.wantString) {
				t.Errorf("got '%s' want '%s'", result, tt.wantString)
			}
		})
	}
}
