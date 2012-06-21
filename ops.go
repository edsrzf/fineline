package fineline

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strings"
)

const (
	// this is intentionally the zero value
	// it saves us a lot of typing
	opPutc = iota
	opPuts
	opBackspace
	opDelete
	opDeleteToBeginning
	opDeleteToEnd
	opClear
	opHome
	opEnd
	opRight
	opLeft
	opUp
	opDown
	opComplete
	opCancel
	opEof
	opSubmit
	opTranspose
	opEscape
	noop
)

var keyMap = [...]int{
	1:    opHome,      // ctrl-a
	2:    opLeft,      // ctrl-b
	3:    opCancel,    // ctrl-c
	4:    opEof,       // ctrl-d
	5:    opEnd,       // ctrl-e
	6:    opRight,     // ctrl-f
	7:    noop,        // ctrl-g
	8:    opBackspace, // ctrl-h
	'\t': opComplete,
	'\r': opSubmit,
	11:   opDeleteToEnd, // ctrl-k
	12:   opClear,       // ctrl-l
	'\n': opSubmit,
	14:   noop,                // ctrl-n
	15:   opSubmit,            // ctrl-o
	16:   opUp,                // ctrl-p
	17:   noop,                // ctrl-q; should be quoted insert
	18:   noop,                // ctrl-r; should be reverse history search
	19:   noop,                // ctrl-s; should be forward history search?
	20:   opTranspose,         // ctrl-t
	21:   opDeleteToBeginning, // ctrl-u
	22:   noop,                // ctrl-v; should be quoted insert
	23:   noop,                // ctrl-w; should be kill word
	24:   noop,                // ctrl-x; should be something?
	25:   noop,                // ctrl-y; should be yank
	26:   noop,                // ctrl-z
	27:   opEscape,
	127:  opBackspace,
}

var cancelled = errors.New("line cancelled")

func (l *LineReader) exec(r *bufio.Reader, op int, c rune) (bool, error) {
	switch op {
	case opPutc:
		l.putc(c)
	case opHome:
		l.home()
	case opLeft:
		l.left()
	case opCancel:
		return false, cancelled
	case opEof:
		if l.pos < l.buf.len()-1 {
			l.delete()
			break
		}
		return false, io.EOF
	case opEnd:
		l.end()
	case opRight:
		l.right()
	case opBackspace:
		l.backspace()
	case opComplete:
		if l.c != nil {
			l.complete()
		} else {
			l.putc(c)
		}
	case opDeleteToEnd:
		l.deleteToEnd()
	case opClear:
		l.clearScreen()
	case opSubmit:
		l.putc('\n')
		l.setCursor(0, -l.y)
		return false, nil
	case opTranspose:
		l.transpose()
	case opDeleteToBeginning:
		l.deleteToBeginning()
	case opEscape:
		var seq [2]byte
		r.Read(seq[:])
		if seq[0] == 91 {
			if 48 < seq[1] && seq[1] < 55 {
				// extended escape
				var seq2 [2]byte
				_, err := r.Read(seq2[:])
				if err != nil {
					break
				}
				if seq[1] == 51 && seq2[0] == 126 {
					// delete
					l.delete()
				}
				break
			}
			switch seq[1] {
			case 65:
				// down arrow
				l.currentEntry++
				if l.currentEntry > len(l.history) {
					l.currentEntry = 0
				}
			case 66:
				// up arrow
				l.currentEntry--
				if l.currentEntry < 0 {
					l.currentEntry = len(l.history) - 1
				}
			case 67:
				// right arrow
				l.right()
			case 68:
				// left arrow
				l.left()
			case 70:
				// end
				l.end()
			case 72:
				// home
				l.home()
			}
		} else if seq[0] == 79 {
			switch seq[1] {
			case 70:
				// end
				l.end()
			case 72:
				// home
				l.home()
			}
		}
	}
	return true, nil
}

func (l *LineReader) putc(c rune) {
	l.buf.WriteRune(c, l.pos)
	l.pos++
	l.refreshLine()
}

func (l *LineReader) puts(s string) {
	l.buf.WriteString(s, l.pos)
	l.pos += len(s)
	l.refreshLine()
}

// finds the longest common string between the end of the first string and the
// beginning of the second
func findIntersect(head, tail string) string {
	n, m := len(head), len(tail)
	if n == 0 || m == 0 {
		return ""
	}
	for len(tail) > 0 {
		pos := strings.LastIndex(head, tail)
		if pos >= 0 {
			return head[pos:]
		}
		tail = tail[:len(tail)-1]
	}
	return ""
}

// finds the longest common prefix of two strings
func commonPrefix(x, y string) string {
	n, m := len(x), len(y)
	if m < n {
		n = m
	}
	i := 0
	for i < n && x[i] == y[i] {
		i++
	}
	return x[:i]
}

func (l *LineReader) complete() {
	if l.display {
		l.printCandidates()
		return
	}
	str := l.buf.String()
	candidates := l.c.Complete(str[:l.pos])
	n := len(candidates)
	if n == 0 {
		return
	}
	var complete string
	if n == 1 {
		complete = candidates[0]
	} else {
		// look for a common prefix to see if we can fill in anything
		prefix := commonPrefix(candidates[0], candidates[1])
		for i := 2; i < n && prefix != str; i++ {
			prefix = commonPrefix(prefix, candidates[i])
		}
		complete = prefix
		l.display = true
		l.candidates = candidates
	}
	inter := findIntersect(str, complete)
	if inter == complete {
		return
	}
	l.puts(complete[len(inter):])
	l.refreshLine()
}

func (l *LineReader) backspace() {
	if l.pos > 0 {
		l.pos--
		l.buf.remove(l.pos)
		l.refreshLine()
	}
}

// delete the character in front of the cursor, like the delete key
func (l *LineReader) delete() {
	l.buf.remove(l.pos)
	l.refreshLine()
}

func (l *LineReader) deleteToBeginning() {
	l.buf.pretruncate(l.pos)
	l.pos = 0
	l.refreshLine()
}

func (l *LineReader) deleteToEnd() {
	l.buf.truncate(l.pos)
	l.refreshLine()
}

func (l *LineReader) home() {
	l.pos = 0
	l.refreshLine()
}

func (l *LineReader) end() {
	l.pos = l.buf.len()
	l.refreshLine()
}

func (l *LineReader) transpose() {
	l.buf.transpose(l.pos)
	l.refreshLine()
}

// move the cursor left
func (l *LineReader) left() {
	if l.pos > 0 {
		l.pos--
		l.refreshLine()
	}
}

// move the cursor right
func (l *LineReader) right() {
	if l.pos < l.buf.len() {
		l.pos++
		l.refreshLine()
	}
}

func (l *LineReader) refreshLine() {
	// move to origin of the current line
	l.setCursor(0, -l.y)
	// assuming the prompt won't wrap
	fmt.Print(l.Prompt)
	bufStr := l.buf.String()
	n := len(bufStr)
	pl := len(l.Prompt)
	if n > l.cols-pl {
		n = l.cols - pl
	}
	fmt.Print(bufStr[:n])
	bufStr = bufStr[n:]
	l.lines = 0
	wrapCursor := n == l.cols-pl
	n = len(bufStr)
	for n > 0 {
		if n > l.cols {
			n = l.cols
		}
		fmt.Print(bufStr[:n])
		bufStr = bufStr[n:]
		wrapCursor = n == l.cols
		n = len(bufStr)
		l.lines++
	}
	l.eraseToEnd()
	if wrapCursor {
		l.lines++
		// move to next line
		fmt.Print("\n")
	}
	x := (pl + l.pos) % l.cols
	l.y = (pl + l.pos) / l.cols
	l.setCursor(x, l.y-l.lines)
}
