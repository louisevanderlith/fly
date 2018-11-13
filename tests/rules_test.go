package tests

import (
	"testing"

	"github.com/louisevanderlith/fly/patterns"
)

//RULES
//Parent Contains main.go (OneCmd)
//Parent Contains other *.go files (OnePkg)
//Parent Contains no Main.go --> Next
//Scan folders;
// Folder has only Folders --> Next
// Folder has main.go (NCmd)
// Folder has other *.go files (NPkg)

func TestRule_RootHasMain_OneCmd(t *testing.T) {

}

func TestRule_RootHasGoFiles_OnePkg(t *testing.T) {

}

func TestRule_RootHasFolders_WithMain_NCmd(t *testing.T) {

}

func TestRule_RootHasFolders_WithGoFiles_NPkg(t *testing.T) {

}

func TestRule_RootHasFolders_ChildHasFolders_None(t *testing.T) {

}

func TestStructure_NCmd_OnlyValid(t *testing.T) {

}

func TestStrucure_NCmdNPkg_3Valid(t *testing.T) {

}

func TestStructure_NCmdOnePkg_3Valid(t *testing.T) {

}

func TestStructure_NPkg_OnlyValid(t *testing.T) {

}

func TestStructure_OneCmd_OnlyValid(t *testing.T) {

}

func TestStructure_OneCmdNPkg_3Valid(t *testing.T) {

}

func TestStructure_OneCmdOnePkg_3Valid(t *testing.T) {

}

func TestStructure_OnePkg_OnlyPkgType(t *testing.T) {
	conf, err := patterns.DetectConfig("onepkg", "TEST")

	if err != nil {
		t.Error(err)
	}

	if len(conf.Programs) != 1 {
		t.Error("Expected one Program")
		return
	}

	if conf.Programs[0].Type != patterns.Pkg {
		t.Error("Incorrect program type")
	}
}

func TestStructure_OnePkg_LibraryName(t *testing.T) {
	conf, err := patterns.DetectConfig("onepkg", "TEST")

	if err != nil {
		t.Error(err)
	}

	if len(conf.Programs) != 1 {
		t.Errorf("Expected one Program; %+v", conf)
		return
	}

	if conf.Programs[0].Name != "onepkg" {
		t.Errorf("Expected 'library', got '%s'", conf.Programs[0].Name)
	}
}
