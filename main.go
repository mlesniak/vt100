// TODO(mlesniak) Everything else are .

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
	// TODO(mlesniak) Add game loop with 1/10s timer, i.e. 100ms per frame.
	for {
		c.PutAt(x, y, '@')
		c.Update()

		b := c.Get()
		c.PutAt(x, y, ' ')

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
