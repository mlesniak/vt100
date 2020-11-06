// TODO(mlesniak) Write to buffer and show only changed values or write everything at once?
// TODO(mlesniak) Everything else are .
// TODO(mlesniak) Move continuous updates to go routine

package main

import (
	"github.com/mlesniak/rogue/vt100"
)

func main() {
	screen := vt100.New()
	screen.Clear()

	w, h := screen.Size()
	x, y := w/2, h/2
Loop:
	for {
		screen.Goto(x, y)
		screen.Put('@')
		screen.Display()

		b := screen.Get()

		screen.Goto(x, y)
		screen.Put(' ')

		switch b {
		case vt100.KeyUp:
			y--
		case vt100.KeyDown:
			y++
		case vt100.KeyRight:
			x++
		case vt100.KeyLeft:
			x--
		case vt100.KeyEscape:
			break Loop
		}
	}

	screen.Clear()
	screen.Goto(0, 0)
	vt100.Reset()
}
