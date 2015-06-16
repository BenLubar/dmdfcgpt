package main

import "github.com/nsf/termbox-go"

const mainMenuMenuOffsetY0 = 5

const (
	mainMenuSettings = iota
	mainMenuExit
	mainMenuCount
)

var mainMenuText [mainMenuCount][]rune = [...][]rune{
	mainMenuSettings: []rune("Settings"),
	mainMenuExit:     []rune("Exit"),
}

var freshMainMenu = MainMenu{
	Menu: Menu{
		Options:  mainMenuText[:],
		Activate: mainMenuActivate,
	},
}

type MainMenu struct {
	Menu Menu
}

func mainMenuActivate(option int) error {
	switch option {
	case mainMenuSettings:
		newFrame := makeSettingsMenu()
		SetFrame(inputFrame, &newFrame)

	case mainMenuExit:
		Tomb.Kill(nil)
	}

	return nil
}

func (mm MainMenu) Render(x0, y0, x1, y1 int) (Frame, error) {
	var out Frame

	if menu, err := mm.Menu.Render(x0, y0+mainMenuMenuOffsetY0, x1, y1); err != nil {
		return nil, err
	} else if menu != nil {
		mm.Menu = menu.(Menu)
		out = mm
	}

	return out, nil
}

func (mm MainMenu) Mouse(mx, my, w, h int) (Frame, error) {
	var out Frame

	if my >= mainMenuMenuOffsetY0 {
		if menu, err := mm.Menu.Mouse(mx, my-mainMenuMenuOffsetY0, w, h-mainMenuMenuOffsetY0); err != nil {
			return nil, err
		} else if menu != nil {
			mm.Menu = menu.(Menu)
			out = mm
		}
	}

	return out, nil
}

func (mm MainMenu) Key(key termbox.Key, mod termbox.Modifier) (Frame, error) {
	var out Frame

	if mod == 0 {
		switch key {
		case termbox.KeyEsc:
			if mm.Menu.Selection == mainMenuExit {
				return nil, mainMenuActivate(mainMenuExit)
			}
			mm.Menu.Selection = mainMenuExit
			out = mm
		}
	}

	if menu, err := mm.Menu.Key(key, mod); err != nil {
		return nil, err
	} else if menu != nil {
		mm.Menu = menu.(Menu)
		out = mm
	}

	return out, nil
}

func (mm MainMenu) Ch(ch rune, mod termbox.Modifier) (Frame, error) {
	var out Frame

	if menu, err := mm.Menu.Ch(ch, mod); err != nil {
		return nil, err
	} else if menu != nil {
		mm.Menu = menu.(Menu)
		out = mm
	}

	return out, nil
}
