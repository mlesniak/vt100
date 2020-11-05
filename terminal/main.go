// screen defines different commands to allow basic line-independent output
// on a VT100 terminal.
//
// Needs termios support for hiding input characters, disabling line
// buffering and getting screen size.
//
// Based on http://ascii-table.com/ansi-escape-sequences-vt-100.php
package terminal

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

func Initialize() {
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

	fmt.Printf("\x1b[?25l")
}

func Reset() {
	fmt.Printf("\x1b[?25h")
	_ = unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
}

func Clear() {
	fmt.Print("\x1b[2J")
}

func Goto(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

func GetPosition() (int, int) {
	fmt.Printf("\x1b[6n")
	x, y := 0, 0
	_, _ = fmt.Scanf("\x1b[%d;%dR", &x, &y)
	return x, y
}

func Put(c byte) {
	fmt.Printf("%c", c)
}

func Get() string {
	n, _ := os.Stdin.Read(inputBuffer)
	return string(inputBuffer[:n])
}

func Size() (int, int) {
	ws, _ := unix.IoctlGetWinsize(fd, unix.TIOCGWINSZ)
	return int(ws.Col), int(ws.Row)
}
