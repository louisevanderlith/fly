package main

import (
	"fmt"
	"io"

	"github.com/logrusorgru/aurora"
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
	gName := aurora.Green(l.progName)
	msg := aurora.Blue(p)
	fmt.Printf(fmt.Sprintf(aurora.Magenta("%s:%s\n %s").String(), l.pid, gName, msg))

	return len(p), nil
}
