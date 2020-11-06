// TODO(mlesniak) Write to buffer and show only changed values or write everything at once?
// TODO(mlesniak) Everything else are .
// TODO(mlesniak) Move continuous updates to go routine

package main

import (
	"github.com/mlesniak/rogue/canvas"
)

func main() {
	c := canvas.New()
	c.Clear()

	w, h := c.Size()
	x, y := w/2, h/2
Loop:
	for {
		c.Goto(x, y)
		c.Put('@')
		c.Display()

		b := c.Get()
		c.Goto(x, y)
		c.Put(' ')

		switch b {
		case canvas.KeyUp:
			y--
		case canvas.KeyDown:
			y++
		case canvas.KeyRight:
			x++
		case canvas.KeyLeft:
			x--
		case canvas.KeyEscape:
			break Loop
		}
	}

	c.Clear()
	c.Goto(0, 0)
	c.Reset()
}
