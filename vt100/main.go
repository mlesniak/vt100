// Screen defines different commands to allow basic line-independent output
// on a VT100 terminal.
//
// Needs termios support for hiding input characters, disabling line
// buffering and getting Screen size.
//
// Based on http://ascii-table.com/ansi-escape-sequences-vt-100.php
package vt100

import (
	"fmt"
	"golang.org/x/sys/unix"
	"os"
)

const ioctlReadTermios = unix.TIOCGETA
const ioctlWriteTermios = unix.TIOCSETA

var fd int
var termios *unix.Termios
var inputBuffer = make([]byte, 3)
var buffer *Screen

type Screen struct {
	width  int
	height int
	data   [][]byte
	x      int
	y      int
}

func New() *Screen {
	fd = int(os.Stdout.Fd())
	termios, err := unix.IoctlGetTermios(fd, ioctlReadTermios)
	if err != nil {
		panic(err)
	}

	newState := *termios
	newState.Lflag &^= unix.ECHO   // Disable echo
	newState.Lflag &^= unix.ICANON // Disable buffering
	if err := unix.IoctlSetTermios(fd, ioctlWriteTermios, &newState); err != nil {
		panic(err)
	}

	// Hide cursor.
	fmt.Printf("\x1b[?25l")

	// Initialize buffer.
	buffer := &Screen{}
	buffer.width, buffer.height = size()
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
	clear()
	for i := 0; i < b.height; i++ {
		fmt.Println(string(b.data[i]))
	}
}

func (b *Screen) Goto(x, y int) {
	b.x, b.y = x, y
}

func (b *Screen) Put(c byte) {
	b.data[b.y][b.x] = c
}

func (b *Screen) Get() string {
	n, _ := os.Stdin.Read(inputBuffer)
	return string(inputBuffer[:n])
}

func (b *Screen) Size() (int, int) {
	return b.width, b.height
}

func Reset() {
	// Show cursor.
	fmt.Printf("\x1b[?25h")
	clear()
	setCursor(0, 0)
	_ = unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
}

func clear() {
	fmt.Print("\x1b[2J")
}

func setCursor(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

func GetPosition() (int, int) {
	fmt.Printf("\x1b[6n")
	x, y := 0, 0
	_, _ = fmt.Scanf("\x1b[%d;%dR", &x, &y)
	return x, y
}

func put(c byte) {
	fmt.Printf("%c", c)
}

func get() string {
	n, _ := os.Stdin.Read(inputBuffer)
	return string(inputBuffer[:n])
}

func size() (int, int) {
	ws, _ := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	return int(ws.Col), int(ws.Row)
}
