package fineline

import (
	"fmt"
	"strings"
	"syscall"
)

type LineReader struct {
	lineReader

	origTerm termios
}

// enable raw mode and gather metrics, like number of columns
func (l *LineReader) raw() {
	tcgetattr(0, &l.origTerm)

	// Modify the original mode
	raw := l.origTerm
	// Input modes - no break, no CR to NL, no parity check, no strip char,
	//               no start/stop output control
	raw.Iflag &^= BRKINT | ICRNL | INPCK | ISTRIP | IXON
	// Output modes - Disable post processing
	raw.Oflag &^= OPOST
	// Control modes - set 8 bit chars
	raw.Cflag |= CS8
	// Local modes - Echo off, canonical off, no extended functions, no signal chars
	raw.Lflag &^= ECHO | ICANON | IEXTEN | ISIG
	raw.Cc[VMIN] = 1
	raw.Cc[VTIME] = 0

	tcsetattr(0, TCSAFLUSH, &raw)

	var win winsize
	winIoctl(1, syscall.TIOCGWINSZ, &win)
	l.cols = int(win.Col)
}

func (l *LineReader) restore() {
	tcsetattr(0, TCSAFLUSH, &l.origTerm)
}

// x is absolute, y is relative
func (l *LineReader) setCursor(x, y int) {
	fmt.Printf("\x1b[%dG", x + 1)
	// positive is down, negative is up
	if y > 0 {
		fmt.Printf("\x1b[%dB", y)
	} else if y < 0 {
		fmt.Printf("\x1b[%dA", -y)
	}
}

// erase everything from the cursor to the end of the screen
func (l *LineReader) eraseToEnd() {
	// erase to right
	fmt.Print("\x1b[0J")
}

func (l *LineReader) printCandidates() {
	str := "\n\x1b[0G" + strings.Join(l.candidates, "\n\x1b[0G") + "\n"
	fmt.Print(str)
	l.refreshLine()
}

func (l *LineReader) clearScreen() {
	// move to upper left corner, then clear entire screen
	fmt.Print("\x1b[H\x1b[2J")
}
