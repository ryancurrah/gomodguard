package cli_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ryancurrah/gomodguard/internal/cli"
)

func TestMain(m *testing.M) {
	err := os.Chdir("../../_example/allOptions")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestCmdRun(t *testing.T) {
	wantExitCode := 2
	exitCode := cli.Run()

	if exitCode != wantExitCode {
		t.Errorf("got exit code '%d' want '%d'", exitCode, wantExitCode)
	}
}
