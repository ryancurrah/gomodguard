package gomodguard_test

import (
	"testing"

	"github.com/ryancurrah/gomodguard"
)

func TestCmdRun(t *testing.T) {
	wantExitCode := 2
	exitCode := gomodguard.Run()

	if exitCode != wantExitCode {
		t.Errorf("got exit code '%d' want '%d'", exitCode, wantExitCode)
	}
}
