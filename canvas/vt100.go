package canvas

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

func showCursor() {
	fmt.Printf("\x1b[?25h")
}

func clear() {
	fmt.Print("\x1b[2J")
}

func setCursor(x, y int) {
	fmt.Printf("\x1b[%d;%dH", y, x)
}

func getPosition() (int, int) {
	fmt.Printf("\x1b[6n")
	x, y := 0, 0
	_, _ = fmt.Scanf("\x1b[%d;%dR", &x, &y)
	return x, y
}

func write(c byte) {
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

func hideCursor() {
	fmt.Printf("\x1b[?25l")
}

func reset() {
	showCursor()
	clear()
	setCursor(0, 0)
	_ = unix.IoctlSetTermios(fd, ioctlWriteTermios, termios)
}

func initVT100() {
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

	hideCursor()
	clear()
}
