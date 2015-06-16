package main

import (
	"os"
	"syscall"

	"github.com/nsf/termbox-go"
)

var inputFrame *Frame

func input() error {
	termbox.SetInputMode(termbox.InputEsc | termbox.InputMouse)

	eventCh := make(chan termbox.Event)
	Tomb.Go(func() error {
		for {
			switch event := termbox.PollEvent(); event.Type {
			case termbox.EventInterrupt:
				return nil

			case termbox.EventError:
				return event.Err

			default:
				select {
				case <-Tomb.Dying():
					// don't return because we still need to get the EventInterrupt.

				case eventCh <- event:
				}
			}
		}
	})
	defer termbox.Interrupt()

	w, h := termbox.Size()

	for {
		select {
		case <-Tomb.Dying():
			return nil

		case event := <-eventCh:
			switch event.Type {
			case termbox.EventResize:
				w, h = event.Width, event.Height

				select {
				case rerender <- struct{}{}:
				default:
				}

			case termbox.EventMouse:
				for {
					inputFrame = CurrentFrame()

					if *inputFrame == nil {
						break
					}

					if newFrame, err := (*inputFrame).Mouse(event.MouseX, event.MouseY, w, h); err != nil {
						return err
					} else if newFrame != nil {
						if !SetFrame(inputFrame, &newFrame) {
							continue
						}
					}
					break
				}

			case termbox.EventKey:
				for {
					inputFrame = CurrentFrame()

					if event.Ch == 0 && event.Key == termbox.KeyCtrlBackslash && event.Mod == 0 {
						// special case: send ourselves SIGQUIT.
						if p, err := os.FindProcess(os.Getpid()); err != nil {
							return err
						} else if err = p.Signal(syscall.SIGQUIT); err != nil {
							return err
						} else {
							break
						}
					}

					if *inputFrame == nil {
						// special case: allow exiting even before the UI initializes.
						if event.Ch == 0 && event.Key == termbox.KeyEsc && event.Mod == 0 {
							Tomb.Kill(nil)
							return nil
						}

						break
					}

					var newFrame Frame
					var err error
					if event.Ch == 0 {
						newFrame, err = (*inputFrame).Key(event.Key, event.Mod)
					} else {
						newFrame, err = (*inputFrame).Ch(event.Ch, event.Mod)
					}
					if err != nil {
						return err
					} else if newFrame != nil {
						if !SetFrame(inputFrame, &newFrame) {
							continue
						}
					}
					break
				}

			default:
				Tomb.Killf("unhandled termbox event type: %v", event.Type)
			}
		}
	}
}
