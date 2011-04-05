package fineline

import (
	"fmt"
	"strings"
)

type term struct {
	termCommon
}

var origTerm termios

// enable raw mode and gather metrics, like number of columns
func (t *term) init() {
	ttyIoctl(0, TCGETS, &origTerm)

	// Modify the original mode
	raw := origTerm
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

	ttyIoctl(0, TCSETSF, &raw)

	var win winsize
	winIoctl(1, TIOCGWINSZ, &win)
	t.cols = int(win.Col)

	t.buf = new(buffer)
}

func (t *term) disableRawMode() {
	ttyIoctl(0, TCSETSF, &origTerm)
}

// x is absolute, y is relative
func (t *term) setCursor(x, y int) {
	fmt.Printf("\x1b[%dG", x + 1)
	// positive is down, negative is up
	if y > 0 {
		fmt.Printf("\x1b[%dB", y)
	} else if y < 0 {
		fmt.Printf("\x1b[%dA", -y)
	}
}

// erase everything from the cursor to the end of the screen
func (t *term) eraseToEnd() {
	// erase to right
	fmt.Print("\x1b[0J")
}

func (t *term) printCandidates() {
	str := "\n\x1b[0G" + strings.Join(t.candidates, "\n\x1b[0G") + "\n"
	fmt.Print(str)
	t.refreshLine()
}

func (t *term) clearScreen() {
	// move to upper left corner, then clear entire screen
	fmt.Print("\x1b[H\x1b[2J")
}
