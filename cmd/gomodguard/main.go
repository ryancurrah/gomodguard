package main

import (
	"os"

	"github.com/ryancurrah/gomodguard/cmd/gomodguard/v2/internal/cli"
)

func main() {
	os.Exit(cli.Run())
}
