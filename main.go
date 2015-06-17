package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/nsf/termbox-go"

	"gopkg.in/tomb.v2"
)

var Tomb tomb.Tomb

func main() {
	flag.Parse()

	// SIGQUIT normally kills the process and prints a stack trace of all
	// goroutines, but termbox prevents that from happening cleanly.
	// We grab SIGQUIT, close termbox, and print the stack trace ourselves,
	// and then immediately exit the process. The deferred termbox close
	// does not run and no cleanup is done. The alternative would be to
	// wait for all processes to cleanly exit, which could be the problem
	// that I'm debugging. We use a buffered channel to make sure only one
	// of the two exits are done.
	exitLimiter := make(chan struct{}, 1)
	exitLimiter <- struct{}{}

	quitch := make(chan os.Signal, 1)
	signal.Notify(quitch, syscall.SIGQUIT)

	go func() {
		<-quitch
		buf := make([]byte, 4096)
		for {
			n := runtime.Stack(buf, true)
			if n < len(buf) {
				buf = buf[:n]
				break
			}
			buf = make([]byte, len(buf)*2)
		}
		<-exitLimiter // hang forever if we're already exiting cleanly
		termbox.Close()
		os.Stderr.Write(buf)
		os.Exit(1)
	}()

	// Catch the terminal going away and exit cleanly.
	hupch := make(chan os.Signal, 1)
	signal.Notify(hupch, syscall.SIGHUP)

	go func() {
		<-hupch
		Tomb.Killf("terminal went away")
	}()

	if err := termbox.Init(); err != nil {
		fmt.Println("cannot start:", err)
	}

	Tomb.Go(startup)

	err := Tomb.Wait()

	<-exitLimiter // hang forever if we're already exiting from SIGQUIT
	termbox.Close()
	if err != nil {
		fmt.Println("fatal error:", err)
	}
}

func startup() error {
	Tomb.Go(renderer)
	Tomb.Go(input)
	Tomb.Go(network)

	frame := makeMainMenu(mainMenuExit)
	SetFrame(CurrentFrame(), &frame)

	return nil
}
