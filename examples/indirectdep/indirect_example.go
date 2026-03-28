package indirectdep

import (
	"os"

	"github.com/gofrs/uuid"
	"github.com/mitchellh/go-homedir"

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

	_, _ = homedir.Expand("~/something")
}
