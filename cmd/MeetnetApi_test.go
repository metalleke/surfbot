package main

import "testing"

func TestTranslate(t *testing.T) {
	translations := []Translation {
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
