package main

import "github.com/nsf/termbox-go"

func makeAddressFrame(replace *Frame, wait <-chan struct{}) Frame {
	var address []rune

	// Get the address unless the network goroutine is busy or not yet
	// started. If we can't get the address, we wait until we're sure
	// the frame has been set, then we wait for the address to be
	// available, and replace the frame if we haven't moved to a different
	// one.
	select {
	case addr := <-ExternalAddr:
		address = []rune(addr)

	default:
		go func() {
			<-wait

			select {
			case addr := <-ExternalAddr:
				var frame Frame = AddressFrame{
					Addr: []rune(addr),
				}
				SetFrame(replace, &frame)

			case <-Tomb.Dying():
			}
		}()
	}

	return AddressFrame{
		Addr: address,
	}
}

var (
	addressTextTitle       = []rune("Network Address")
	addressTextWaiting     = []rune("(waiting for networking module)")
	addressTextDescription = splitWords("Give this address to your friends so they can connect to your game. Press any key to go back to the main menu.")
)

type AddressFrame struct {
	Addr []rune
}

func (f AddressFrame) Render(x0, y0, x1, y1 int) (Frame, error) {
	y := y0

	if y1-y0 > 4 {
		// TODO: title
		y += 2
	}

	if f.Addr != nil {
		fg, bg := termbox.ColorWhite, termbox.ColorBlack
		if CurrentSettings().HighContrast {
			fg, bg = bg, fg
			for x := x0; x < x1; x++ {
				termbox.SetCell(x, y, ' ', fg, bg)
			}
		} else {
			fg |= termbox.AttrBold
		}
		x := x0 + (x1-x0)/2 - len(f.Addr)/2
		if x < x0 {
			x = x0
		}
		for dx, r := range f.Addr {
			if x+dx >= x1 {
				break
			}
			termbox.SetCell(x+dx, y, r, fg, bg)
		}
	} else {
		fg, bg := termbox.ColorBlack|termbox.AttrBold, termbox.ColorBlack
		if CurrentSettings().HighContrast {
			fg = termbox.ColorWhite
		}
		x := x0 + (x1-x0)/2 - len(addressTextWaiting)/2
		if x < x0 {
			x = x0
		}
		for dx, r := range addressTextWaiting {
			if x+dx >= x1 {
				break
			}
			termbox.SetCell(x+dx, y, r, fg, bg)
		}
	}
	y += 2

	x := x0 + 1
	for _, word := range addressTextDescription {
		if x != x0+1 && x+len(word) > x1-1 {
			x = x0 + 1
			y++
		}
		if y >= y1 {
			break
		}

		for dx, r := range word {
			if x+dx >= x1-1 {
				break
			}
			termbox.SetCell(x+dx, y, r, termbox.ColorWhite, termbox.ColorBlack)
		}

		x += len(word)
	}

	return nil, nil
}
func (f AddressFrame) Mouse(mx, my, w, h int) (Frame, error) {
	return nil, nil
}
func (f AddressFrame) Key(key termbox.Key, mod termbox.Modifier) (Frame, error) {
	return makeMainMenu(mainMenuAddress), nil
}
func (f AddressFrame) Ch(ch rune, mod termbox.Modifier) (Frame, error) {
	return makeMainMenu(mainMenuAddress), nil
}
