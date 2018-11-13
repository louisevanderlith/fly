package tests

import (
	"testing"

	"github.com/louisevanderlith/fly/patterns"
)

func TestStructure_NCmd_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("ncmd")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewNCmd("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}

func TestStructure_NCmdNPkg_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("ncmdnpkg")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewNCmdNPkg("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}

func TestStructure_NCmdOnePkg_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("ncmdonepkg")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewNCmdOnePkg("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}

func TestStructure_NPkg_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("npkg")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewNPkg("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}

func TestStructure_OneCmd_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("onecmd")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewOneCmd("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}

func TestStructure_OneCmdNPkg_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("onecmdnpkg")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewOneCmdNPkg("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}

func TestStructure_OneCmdOnePkg_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("onecmdonepkg")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewOneCmdOnePkg("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}
func TestStructure_OnePkg_LoadsFull(t *testing.T) {
	structure, err := patterns.ScanFolder("onepkg")

	if err != nil {
		t.Error(err)
	}

	pattrn := patterns.NewOnePkg("TEST", structure)
	pass := pattrn.Test()

	if !pass {
		t.Error("Pattern didn't Pass")
	}
}
