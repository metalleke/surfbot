package main

import "testing"

func TestMeterePerSecondToBeaufortScale(t *testing.T) {
	var result int

	// Unknown
	result = MeterePerSecondToBeaufortScale(-1)
	if result != -1 {
		t.Errorf("Wanted -1 got: %d", result)
	}

	// 0 bft
	result = MeterePerSecondToBeaufortScale(0)
	if result != 0 {
		t.Errorf("Wanted 0 got: %d", result)
	}

	result = MeterePerSecondToBeaufortScale(0.2)
	if result != 0 {
		t.Errorf("Wanted 0 got: %d", result)
	}

	result = MeterePerSecondToBeaufortScale(0.4)
	if result != 0 {
		t.Errorf("Wanted 0 got: %d", result)
	}

	// 1 bft
	result = MeterePerSecondToBeaufortScale(0.5)
	if result != 1 {
		t.Errorf("Wanted 1 got: %d", result)
	}

	// Unknown
	result = MeterePerSecondToBeaufortScale(100)
	if result != -1 {
		t.Errorf("Wanted -1 got: %d", result)
	}
}
