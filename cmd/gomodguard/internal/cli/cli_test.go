package cli_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/ryancurrah/gomodguard/cmd/gomodguard/v2/internal/cli"
	"github.com/ryancurrah/gomodguard/v2"
)

func TestCmdRun(t *testing.T) {
	t.Chdir(examplesDir + "alloptions")

	wantExitCode := 2
	exitCode := cli.Run()

	if exitCode != wantExitCode {
		t.Errorf("got exit code '%d' want '%d'", exitCode, wantExitCode)
	}
}

func TestWriteCheckstyle(t *testing.T) {
	outFile, err := os.CreateTemp(t.TempDir(), "checkstyle-*.xml")
	require.NoError(t, err)

	defer func() {
		err := outFile.Close()
		require.NoError(t, err)
	}()

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
