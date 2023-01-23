package gomodguard

import (
	"os"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/go-homedir"
	"github.com/ryancurrah/gomodguard"
	module "github.com/uudashr/go-module"
)

func aBlockedImport() { //nolint: deadcode,unused
	b, err := os.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}

	mod, err := module.Parse(b)
	if err != nil {
		panic(err)
	}

	_ = mod

	_ = uuid.Must(uuid.NewV4())

	var blockedModule gomodguard.BlockedModule
	blockedModule.Set("some.com/module/name")

	_, _ = homedir.Expand("~/something")
}
