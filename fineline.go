// Copyright (c) 2011, Evan Shaw.

package fineline

import (
	"bufio"
	"os"
	"strings"
)

// Common, platform-independent components
type lineReader struct {
	// The prompt that precedes any line entry
	Prompt string
	input  *bufio.Reader

	// A circular array of history
	history []string
	// The most recent line in the history
	lastEntry int
	// The index of the line in history currently being shown
	// -1 if we're not showing something in history.
	currentEntry int

	buf    buffer
	// number of lines we last wrote
	lines     int
	pos, cols int
	c         Completer
	// candidates from last tab completion
	candidates []string
	display    bool
	y          int
}

// NewLineReader creates a new LineReader that reads from stdin.
func NewLineReader() *LineReader {
	var l LineReader
	l.input = bufio.NewReader(os.Stdin)
	l.Prompt = "$ "
	return &l
}

func (l *LineReader) SetMaxHistoryLen(len int) {
	l.history = make([]string, len)
	l.lastEntry = 0
}

func (l *LineReader) AddHistory(line string) {
	if len(l.history) > 0 {
		l.lastEntry++
		if l.lastEntry >= len(l.history) {
			l.lastEntry = 0
		}
		l.history[l.lastEntry] = line
	}
}

var unsupportedTerms = [...]string{"dumb", "cons25"}

func unsupportedTerm() bool {
	term := strings.ToLower(os.Getenv("TERM"))
	for i := range unsupportedTerms {
		if term == unsupportedTerms[i] {
			return true
		}
	}
	return false
}

func (l *LineReader) Read(c Completer) (line string, err error) {
	// TODO: Move these checks to NewLineReader()
	// TODO: Check if STDIN is a TTY
	if unsupportedTerm() {
		// Fall back to plain old stdin reading
		r := bufio.NewReader(os.Stdin)
		line, err = r.ReadString('\n')
	} else {
		l.raw()
		line, err = l.getLine()
		l.restore()
		l.buf.reset()
		l.pos = 0
	}
	return
}

func (l *LineReader) getLine() (string, error) {
	r := bufio.NewReader(os.Stdin)
	l.refreshLine()
	var err error
	cont := true
	for cont && err == nil {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		var op int
		if int(c) < len(keyMap) {
			op = keyMap[c]
		} else {
			op = opPutc
		}
		cont, err = l.exec(r, op, c)
	}
	if err == cancelled {
		return "", nil
	}

	return l.buf.String(), err
}
