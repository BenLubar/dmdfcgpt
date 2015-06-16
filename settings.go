package main

import (
	"sync"

	"github.com/nsf/termbox-go"
)

type Settings struct {
	HighContrast bool
}

var currentSettings Settings
var settingsLock sync.Mutex

func CurrentSettings() Settings {
	settingsLock.Lock()
	settings := currentSettings
	settingsLock.Unlock()

	return settings
}

func makeSettingsMenu() Frame {
	settings := CurrentSettings()

	updateSettings := func() {
		settingsLock.Lock()
		currentSettings = settings
		settingsLock.Unlock()

		select {
		case rerender <- struct{}{}:
		default:
		}
	}

	return SettingsMenu{
		HighContrast: Checkbox{
			Checked: settings.HighContrast,
			Activate: func(ok bool) error {
				settings.HighContrast = ok

				updateSettings()

				return nil
			},
		},
	}
}

const (
	settingsMenuHighContrast = iota
	settingsMenuCount
)

var settingsMenuText [settingsMenuCount][]rune = [...][]rune{
	settingsMenuHighContrast: []rune("High contrast mode"),
}

type SettingsMenu struct {
	Selection    int
	HighContrast Frame
}

func (sm *SettingsMenu) frames() [settingsMenuCount]*Frame {
	return [...]*Frame{
		settingsMenuHighContrast: &sm.HighContrast,
	}
}

func (sm SettingsMenu) Render(x0, y0, x1, y1 int) (Frame, error) {
	var out Frame

	highContrast := sm.HighContrast.(Checkbox).Checked

	center := x0 + (x1-x0)/2

	for i, frame := range sm.frames() {
		y := y0 + i
		if y >= y1 {
			break
		}

		fg, bg := termbox.ColorWhite, termbox.ColorBlack
		if sm.Selection == i {
			if highContrast {
				fg, bg = bg, fg
				for x := x0; x < center; x++ {
					termbox.SetCell(x, y, ' ', fg, bg)
				}
			} else {
				fg |= termbox.AttrBold
			}
		}
		for dx, r := range settingsMenuText[i] {
			x := x0 + dx
			if x >= center {
				break
			}
			termbox.SetCell(x, y, r, fg, bg)
		}
		if newFrame, err := (*frame).Render(center, y, x1, y+1); err != nil {
			return nil, err
		} else if newFrame != nil {
			*frame = newFrame
			out = sm
		}
	}

	return out, nil
}

func (sm SettingsMenu) Key(key termbox.Key, mod termbox.Modifier) (Frame, error) {
	if mod == 0 {
		switch key {
		case termbox.KeyArrowUp:
			sm.Selection = (sm.Selection + settingsMenuCount - 1) % settingsMenuCount
			return sm, nil

		case termbox.KeyArrowDown:
			sm.Selection = (sm.Selection + 1) % settingsMenuCount
			return sm, nil

		case termbox.KeyEsc:
			return freshMainMenu, nil
		}
	}

	frame := sm.frames()[sm.Selection]
	if newFrame, err := (*frame).Key(key, mod); err != nil {
		return nil, err
	} else if newFrame != nil {
		*frame = newFrame
		return sm, nil
	}

	return nil, nil
}

func (sm SettingsMenu) Ch(ch rune, mod termbox.Modifier) (Frame, error) {
	frame := sm.frames()[sm.Selection]
	if newFrame, err := (*frame).Ch(ch, mod); err != nil {
		return nil, err
	} else if newFrame != nil {
		*frame = newFrame
		return sm, nil
	}

	return nil, nil
}

func (sm SettingsMenu) Mouse(mx, my, w, h int) (Frame, error) {
	if mx >= w/2 && my < settingsMenuCount {
		sm.Selection = my
		frame := sm.frames()[my]

		if newFrame, err := (*frame).Mouse(mx-w/2, 0, w-w/2, 1); err != nil {
			return nil, err
		} else if newFrame != nil {
			*frame = newFrame
		}
		return sm, nil
	}

	return nil, nil
}
