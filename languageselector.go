package main

import "github.com/nsf/termbox-go"

const (
	languageSelectorPrev = '◄'
	languageSelectorNext = '►'
)

type LanguageSelector struct {
	Index    int
	Activate func(string) error
}

func (f LanguageSelector) Render(x0, y0, x1, y1 int) (Frame, error) {
	termbox.SetCell(x0, y0, languageSelectorPrev, termbox.ColorWhite, termbox.ColorBlack)
	text := []rune(TranslationIDs[f.Index])
	x := x0 - (x1-x0)/2
	if x < x0+2 {
		x = x0 + 2
	}
	for dx, r := range text {
		if x+dx >= x1-2 {
			break
		}
		termbox.SetCell(x+dx, y0, r, termbox.ColorWhite, termbox.ColorBlack)
	}
	termbox.SetCell(x1-1, y0, languageSelectorNext, termbox.ColorWhite, termbox.ColorBlack)

	return nil, nil
}

func (f LanguageSelector) next(offset int) (Frame, error) {
	f.Index = (f.Index + len(TranslationIDs) + offset) % len(TranslationIDs)

	return f, f.Activate(TranslationIDs[f.Index])
}

func (f LanguageSelector) Mouse(mx, my, w, h int) (Frame, error) {
	if mx > w/2 {
		return f.next(1)
	}
	return f.next(-1)
}

func (f LanguageSelector) Key(key termbox.Key, mod termbox.Modifier) (Frame, error) {
	if mod != 0 {
		return nil, nil
	}
	switch key {
	case termbox.KeyArrowLeft:
		return f.next(-1)

	case termbox.KeyArrowRight:
		return f.next(1)
	}

	return nil, nil
}

func (f LanguageSelector) Ch(ch rune, mod termbox.Modifier) (Frame, error) {
	return nil, nil
}
