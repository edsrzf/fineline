// Copyright (c) 2011, Evan Shaw.

package fineline

import (
	"bufio"
	"os"
	"strings"
)

type termCommon struct {
	prompt string
	buf *buffer
	// number of lines we last wrote
	lines int
	pos, cols int
	c Completer
	// candidates from last tab completion
	candidates []string
	display bool
	y int
}

// A circular array of history
var history []string

// The most recent line in the history
var lastEntry int

// The index of the line in history currently being shown
// -1 if we're not showing something in history.
var currentEntry int

func SetMaxHistoryLen(len int) {
	history = make([]string, len)
	lastEntry = 0
}

func AddHistory(line string) {
	if len(history) > 0 {
		lastEntry++
		if lastEntry >= len(history) {
			lastEntry = 0
		}
		history[lastEntry] = line
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

func Read(prompt string) (line string, err os.Error) {
	// TODO: Check if STDIN is a TTY
	if unsupportedTerm() {
		// Fall back to plain old stdin reading
		r := bufio.NewReader(os.Stdin)
		line, err = r.ReadString('\n')
	} else {
		t := &term{}
		t.c = NewSimpleCompleter([]string{"car", "cat", "dog"})
		t.prompt = prompt
		t.init()
		line, err = t.getLine()
		t.disableRawMode()
	}
	return
}

func (t *term) getLine() (string, os.Error) {
	r := bufio.NewReader(os.Stdin)
	t.refreshLine()
	var err os.Error
	cont := true
	for cont && err == nil {
		c, _, err := r.ReadRune()
		if err != nil {
			return "", err
		}
		var op int
		if c < len(keyMap) {
			op = keyMap[c]
		} else {
			op = opPutc
		}
		cont, err = t.exec(r, op, c)
	}
	if err == cancelled {
		return "", nil
	}

	return t.buf.String(), err
}
