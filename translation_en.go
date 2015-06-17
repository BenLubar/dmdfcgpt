package main

var EnglishTranslation = Translation{
	MainMenu: [...][]rune{
		mainMenuAddress:  []rune("Address"),
		mainMenuSettings: []rune("Settings"),
		mainMenuExit:     []rune("Exit"),
	},
	Settings: [...][]rune{
		settingsMenuHighContrast: []rune("High contrast mode"),
		settingsMenuLanguage:     []rune("Language"),
	},
}
