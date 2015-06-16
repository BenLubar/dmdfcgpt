package main

import "github.com/nsf/termbox-go"

type Menu struct {
	Selection int
	Options   [][]rune
	Activate  func(int) error
}

func (f Menu) layout(option int, x0, y0, x1, y1 int) (text []rune, x, y int) {
	var offset int

	// if y1 == y0, we can have 1 option without scrolling, and so on.
	if len(f.Options) > y1-y0 {
		offset = -f.Selection + (y1-y0)/2
	}

	text = f.Options[option]
	x = (x0 + x1 - len(text)) / 2
	// always display at least the first character.
	if x < x0 {
		x = x0
	}
	y = y0 + option + offset
	return
}

func (f Menu) Render(x0, y0, x1, y1 int) (Frame, error) {
	settings := CurrentSettings()

	for i := range f.Options {
		text, x, y := f.layout(i, x0, y0, x1, y1)
		fg, bg := termbox.ColorWhite, termbox.ColorBlack
		if y < y0 || y >= y1 {
			continue
		}
		if i == f.Selection {
			if settings.HighContrast {
				fg, bg = bg, fg
				for dx := x0; dx < x1; dx++ {
					termbox.SetCell(dx, y, ' ', fg, bg)
				}
			} else {
				fg |= termbox.AttrBold
			}
		}
		for dx, r := range text {
			if x+dx >= x0 && x+dx < x1 {
				termbox.SetCell(x+dx, y, r, fg, bg)
			}
		}
	}

	return nil, nil
}

func (f Menu) Mouse(mx, my, w, h int) (Frame, error) {
	settings := CurrentSettings()

	for i := range f.Options {
		text, x, y := f.layout(i, 0, 0, w, h)
		if y != my {
			continue
		}
		if !settings.HighContrast && (mx < x || mx > x+len(text)) {
			continue
		}
		if i == f.Selection {
			return nil, f.Activate(i)
		}
		f.Selection = i
		return f, nil
	}
	return nil, nil
}

func (f Menu) Key(key termbox.Key, mod termbox.Modifier) (Frame, error) {
	if mod != 0 {
		return nil, nil
	}

	switch key {
	case termbox.KeyArrowUp:
		f.Selection = (f.Selection + len(f.Options) - 1) % len(f.Options)
		return f, nil

	case termbox.KeyArrowDown:
		f.Selection = (f.Selection + 1) % len(f.Options)
		return f, nil

	case termbox.KeyEnter:
		return nil, f.Activate(f.Selection)
	}
	return nil, nil
}

func (f Menu) Ch(ch rune, mod termbox.Modifier) (Frame, error) {
	return nil, nil
}
