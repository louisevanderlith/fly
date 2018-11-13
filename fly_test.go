package main

import (
	"testing"

	"github.com/louisevanderlith/fly/patterns"
)

func TestLoadConfig_PriorityCorrect_RouterFirst(t *testing.T) {
	expect := "router"
	actual, err := patterns.DetectConfig(".", "TEST")

	if err != nil {
		t.Error(err)
	}

	if actual.Programs[0].Name != expect {
		t.Errorf("Expected %s, got %s", expect, actual.Programs[0].Name)
	}
}
