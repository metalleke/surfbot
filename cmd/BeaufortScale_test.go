package main

import "testing"

func TestMeterePerSecondToBeaufortScale(t *testing.T) {
	var result string

	// Unknown
	result = MeterePerSecondToBeaufortScale(-1)
	if result != "Unknown value" {
		t.Errorf("Wanted unknown value got: %s", result)
	}

	// 0 bft
	result = MeterePerSecondToBeaufortScale(0)
	if result != "0 bft" {
		t.Errorf("Wanted 0 bft got: %s", result)
	}

	result = MeterePerSecondToBeaufortScale(0.2)
	if result != "0 bft" {
		t.Errorf("Wanted 0 bft got: %s", result)
	}

	result = MeterePerSecondToBeaufortScale(0.4)
	if result != "0 bft" {
		t.Errorf("Wanted 0 bft got: %s", result)
	}

	// 1 bft
	result = MeterePerSecondToBeaufortScale(0.5)
	if result != "1 bft" {
		t.Errorf("Wanted 1 bft bft got: %s", result)
	}

	// Unknown
	result = MeterePerSecondToBeaufortScale(100)
	if result != "Unknown value" {
		t.Errorf("Wanted unknown value got: %s", result)
	}
}