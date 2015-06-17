package main

import (
	"sort"
	"strings"
)

type Translation struct {
	MainMenu [mainMenuCount][]rune
	Settings [settingsMenuCount][]rune
}

func CurrentTranslation() *Translation {
	return Translations[CurrentSettings().Language]
}

func splitWords(text string) (words [][]rune) {
	for _, word := range strings.SplitAfter(text, " ") {
		words = append(words, []rune(word))
	}
	return
}

var Translations = map[string]*Translation{
	"English":    &EnglishTranslation,
	"la lojban.": &LojbanTranslation,
}

var TranslationIDs = func() (ids []string) {
	for id := range Translations {
		ids = append(ids, id)
	}
	sort.Strings(ids)
	return
}()
