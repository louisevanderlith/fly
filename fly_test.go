package main

import "testing"

func TestLoadConfig_PriorityCorrect_RouterFirst(t *testing.T) {
	expect := "router"
	actual, err := loadConfig()

	if err != nil {
		t.Error(err)
	}

	if actual.Programs[0].Name != expect {
		t.Errorf("Expected %s, got %s", expect, actual.Programs[0].Name)
	}
}
