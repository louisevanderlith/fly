package tests

import (
	"testing"

	"github.com/louisevanderlith/fly/patterns"
)

func TestStructure_Std_LibraryName(t *testing.T) {
	conf, err := patterns.DetectConfig("thestandard", "TEST")

	if err != nil {
		t.Error(err)
	}

	if conf.Programs[0].Name != "apiWConf" {
		t.Errorf("Expected 'library', got '%s'", conf.Programs[0].Name)
	}
}
