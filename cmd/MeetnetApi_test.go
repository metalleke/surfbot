package main

import (
	"fmt"
	"testing"
)

func Testtranslate(t *testing.T) {
	translations := []Translation{
		{
			Locale:  "en-NL",
			Message: "A message",
		},
	}
	locale := "en-NL"

	got := translate(translations, locale)
	if got != "A message" {
		t.Errorf("Did not get the right message got: %s", got)
	}
}

func TestsafeToKite(t *testing.T) {
	var result bool

	result = safeToKite(0)
	if result == false {
		t.Error("This should be safe to kite")
	}

	result = safeToKite(6)
	fmt.Println(result)
	if result == false {
		t.Error("This should be safe to kite")
	}

	result = safeToKite(7)
	if result == true {
		t.Error("This should not be safe to kite")
	}
}
