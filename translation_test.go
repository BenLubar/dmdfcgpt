package main

import "testing"

func TestTranslations(t *testing.T) {
	for _, id := range TranslationIDs {
		translation := Translations[id]

		for i, text := range translation.MainMenu {
			if text == nil {
				t.Errorf("Missing translation in %q for main menu option #%d.", id, i)
			}
		}

		for i, text := range translation.Settings {
			if text == nil {
				t.Errorf("Missing translation in %q for setting #%d.", id, i)
			}
		}
	}
}
