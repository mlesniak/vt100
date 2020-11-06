// Screen defines different commands to allow basic line-independent output
// on a VT100 terminal.
//
// Needs termios support for hiding input characters, disabling line
// buffering and getting Screen size.
//
// Based on http://ascii-table.com/ansi-escape-sequences-vt-100.php
//
// Design step: abstract canvas operations and use interface in virtual buffer.
package canvas

import (
	"os"
)

var buffer *Screen

type Screen struct {
	width  int
	height int
	data   [][]byte
	x      int
	y      int

	ops map[pos]byte
}

type pos struct {
	x int
	y int
}

func New() *Screen {
	initVT100()

	// Initialize buffer.
	buffer := &Screen{}
	buffer.width, buffer.height = size()
	buffer.ops = make(map[pos]byte)
	buffer.data = make([][]byte, buffer.height)
	for i := 0; i < buffer.height; i++ {
		buffer.data[i] = make([]byte, buffer.width)
	}

	return buffer
}

func (b *Screen) Clear() {
	for i := 0; i < b.height; i++ {
		for j := 0; j < b.width; j++ {
			b.data[i][j] = ' '
		}
	}

	b.x, b.y = 0, 0
}

func (b *Screen) Display() {
	for pos, val := range b.ops {
		setCursor(pos.x, pos.y)
		write(val)
		delete(b.ops, pos)
	}
}

func (b *Screen) Goto(x, y int) {
	b.x, b.y = x, y
}

func (b *Screen) Put(c byte) {
	b.ops[pos{b.x, b.y}] = c
}

func (b *Screen) Get() string {
	n, _ := os.Stdin.Read(inputBuffer)
	return string(inputBuffer[:n])
}

func (b *Screen) Size() (int, int) {
	return b.width, b.height
}

func (b *Screen) Reset() {
	reset()
}
