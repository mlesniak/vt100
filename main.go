package main

import (
	"github.com/mlesniak/rogue/terminal"
)

func main() {
	terminal.Initialize()

	// TODO(mlesniak) Write to buffer and show only changed values or write everything at once?
	// TODO(mlesniak) Everything else are .

	terminal.Clear()
	w, h := terminal.Size()

	x, y := w/2, h/2
Loop:
	for {
		terminal.Goto(x, y)
		terminal.Put('@')

		b := terminal.Get()

		terminal.Goto(x, y)
		terminal.Put(' ')

		switch b {
		case "\x1b[A":
			y--
		case "\x1b[B":
			y++
		case "\x1b[C":
			x++
		case "\x1b[D":
			x--
		case "\x1b":
			break Loop
		}

	}

	terminal.Clear()
	terminal.Goto(0, 0)
	terminal.Reset()
}
