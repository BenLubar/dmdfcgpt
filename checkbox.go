package main

import "github.com/nsf/termbox-go"

type Checkbox struct {
	Checked  bool
	Activate func(bool) error
}

func (f Checkbox) Render(x0, y0, x1, y1 int) (Frame, error) {
	if f.Checked {
		termbox.SetCell(x0, y0, checkboxChecked, termbox.ColorWhite, termbox.ColorBlack)
	} else {
		termbox.SetCell(x0, y0, checkboxUnchecked, termbox.ColorWhite, termbox.ColorBlack)
	}
	return nil, nil
}

func (f Checkbox) Key(key termbox.Key, mod termbox.Modifier) (Frame, error) {
	if mod != 0 {
		return nil, nil
	}

	switch key {
	case termbox.KeySpace, termbox.KeyEnter:
		f.Checked = !f.Checked
		return f, f.Activate(f.Checked)
	}

	return nil, nil
}

func (f Checkbox) Ch(ch rune, mod termbox.Modifier) (Frame, error) {
	return nil, nil
}

func (f Checkbox) Mouse(mx, my, w, h int) (Frame, error) {
	if mx == 0 && my == 0 {
		f.Checked = !f.Checked
		return f, f.Activate(f.Checked)
	}
	return nil, nil
}
