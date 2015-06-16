package main

import (
	"log"
	"os"
	"sync"
)

var Log = log.New(&deferredFileWriter{name: "err.log"}, "", log.Ldate|log.Ltime|log.Lmicroseconds|log.Lshortfile)

type deferredFileWriter struct {
	name string
	file *os.File
	open sync.Once
	err  error
}

func (w *deferredFileWriter) Write(p []byte) (int, error) {
	w.open.Do(func() {
		w.file, w.err = os.OpenFile(w.name, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	})

	if w.err != nil {
		return 0, w.err
	}

	return w.file.Write(p)
}
