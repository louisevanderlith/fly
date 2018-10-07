package main

import (
	"io"
	"log"
)

type logger struct {
	progName string
	pid      string
}

func newLogger(progName, pid string) io.Writer {
	return logger{
		progName: progName,
		pid:      pid,
	}
}

func (l logger) Write(p []byte) (int, error) {
	log.Printf("%s:%s\n %s", l.pid, l.progName, string(p))

	return len(p), nil
}
