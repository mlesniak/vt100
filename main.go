// TODO(mlesniak) Everything else are .

package main

import (
	"github.com/mlesniak/rogue/canvas"
	"time"
)

func main() {
	c := canvas.New()
	c.Clear()

	w, h := c.Size()
	x, y := w/2, h/2

	input := make(chan string)

	// Input loop.
	go func() {
		for {
			b := c.Get()
			input <- b
		}
	}()

	passedTime := time.Now()
	c.PutAt(x, y, '@')
Loop:
	for {
		// Show current time.
		if passedTime.Add(time.Second).Before(time.Now()) {
			passedTime = time.Now()
			c.PrintAt(1, 0, time.Now().String())
		}

		select {
		case b := <-input:
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
		default:
			// Do nothing.
		}

		c.PutAt(x, y, '@')

		c.Update()
		time.Sleep(time.Millisecond * 10)
	}

	c.Clear()
	c.Goto(0, 0)
	c.Reset()
}
