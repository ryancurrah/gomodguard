package test

import (
	"io/ioutil"

	module "github.com/uudashr/go-module"
)

func aBlockedImport() { // nolint: deadcode,unused
	b, err := ioutil.ReadFile("go.mod")
	if err != nil {
		panic(err)
	}

	mod, err := module.Parse(b)
	if err != nil {
		panic(err)
	}

	_ = mod
}
