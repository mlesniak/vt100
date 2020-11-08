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
	"sync"
	"time"
)

var buffer *Screen

type Screen struct {
	width  int
	height int
	data   [][]byte
	x      int
	y      int

	opsLock sync.Mutex
	ops     map[pos]byte
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

func (s *Screen) Clear() {
	for i := 0; i < s.height; i++ {
		for j := 0; j < s.width; j++ {
			s.data[i][j] = ' '
		}
	}

	s.x, s.y = 0, 0
}

func (s *Screen) Update() {
	s.opsLock.Lock()
	defer s.opsLock.Unlock()
	for pos, val := range s.ops {
		setCursor(pos.x, pos.y)
		write(val)
		delete(s.ops, pos)
	}
}

func (s *Screen) Goto(x, y int) {
	s.x, s.y = x, y
}

func (s *Screen) PutAt(x, y int, c byte) {
	s.Goto(x, y)
	s.Put(c)
}

func (s *Screen) Put(c byte) {
	s.opsLock.Lock()
	s.ops[pos{s.x, s.y}] = c
	s.opsLock.Unlock()
}

func (s *Screen) Get() string {
	n, _ := os.Stdin.Read(inputBuffer)
	return string(inputBuffer[:n])
}

func (s *Screen) Size() (int, int) {
	return s.width, s.height
}

func (s *Screen) Reset() {
	reset()
}

func (s *Screen) StartDisplayLoop() {
	//var fps int64 = 24
	go func() {
		for {
			s.Update()
			//time.Sleep(time.Duration(time.Second.Milliseconds() / fps) * time.Millisecond
			time.Sleep(time.Millisecond / 24)
		}
	}()
}
