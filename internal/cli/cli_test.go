package cli_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/ryancurrah/gomodguard"
	"github.com/ryancurrah/gomodguard/internal/cli"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestWriteCheckstyle(t *testing.T) {
	outFile, err := os.CreateTemp(t.TempDir(), "checkstyle-*.xml")
	require.NoError(t, err)
	defer outFile.Close()

	issues := []gomodguard.Issue{
		{
			FileName:   "first.go",
			LineNumber: 10,
			Reason:     "first test reason",
		},
		{
			FileName:   "second.go",
			LineNumber: 20,
			Reason:     "second test reason",
		},
	}

	err = cli.WriteCheckstyle(outFile.Name(), issues)
	require.NoError(t, err)

	got, err := os.ReadFile(outFile.Name())
	require.NoError(t, err)

	want := `
<?xml version="1.0" encoding="UTF-8"?>
<checkstyle version="1.0.0">
  <file name="first.go">
    <error line="10" column="1" severity="error" message="first test reason" source="gomodguard"></error>
  </file>
  <file name="second.go">
    <error line="20" column="1" severity="error" message="second test reason" source="gomodguard"></error>
  </file>
</checkstyle>`
	assert.Equal(t, want, string(got))
}
