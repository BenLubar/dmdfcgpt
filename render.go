package main

import "github.com/nsf/termbox-go"

var rerender = make(chan struct{}, 1)

func renderer() error {
	for {
		select {
		case <-Tomb.Dying():
			return nil

		case <-rerender:
		}

		for {
			frame := CurrentFrame()

			if *frame == nil {
				break
			}

			if err := termbox.Clear(termbox.ColorBlack, termbox.ColorBlack); err != nil {
				return err
			}
			w, h := termbox.Size()
			if newFrame, err := (*frame).Render(0, 0, w, h); err != nil {
				return err
			} else if newFrame != nil {
				if !SetFrame(frame, &newFrame) {
					continue
				}
			}
			break
		}
		if err := termbox.Flush(); err != nil {
			return err
		}
	}
}
