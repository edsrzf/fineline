package fineline

import (
	"bufio"
	"fmt"
	"os"
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

var keyMap = [...]int {
	1: opHome, // ctrl-a
	2: opLeft, // ctrl-b
	3: opCancel, // ctrl-c
	4: opEof, // ctrl-d
	5: opEnd, // ctrl-e
	6: opRight, // ctrl-f
	7: noop, // ctrl-g
	8: opBackspace, // ctrl-h
	'\t': opComplete,
	'\r': opSubmit,
	11: opDeleteToEnd, // ctrl-k
	12: opClear, // ctrl-l
	'\n': opSubmit,
	14: noop, // ctrl-n
	15: opSubmit, // ctrl-o
	16: opUp, // ctrl-p
	17: noop, // ctrl-q; should be quoted insert
	18: noop, // ctrl-r; should be reverse history search
	19: noop, // ctrl-s; should be forward history search?
	20: opTranspose, // ctrl-t
	21: opDeleteToBeginning, // ctrl-u
	22: noop, // ctrl-v; should be quoted insert
	23: noop, // ctrl-w; should be kill word
	24: noop, // ctrl-x; should be something?
	25: noop, // ctrl-y; should be yank
	26: noop, // ctrl-z
	27: opEscape,
	127: opBackspace,
}

var cancelled = os.NewError("line cancelled")

func (t *term) exec(r *bufio.Reader, op, c int) (bool, os.Error) {
	switch op {
	case opPutc:
		t.putc(c)
	case opHome:
		t.home()
	case opLeft:
		t.left()
	case opCancel:
		return false, cancelled
	case opEof:
		if t.pos < t.buf.len() - 1 {
			t.delete()
			break
		}
		return false, os.EOF
	case opEnd:
		t.end()
	case opRight:
		t.right()
	case opBackspace:
		t.backspace()
	case opComplete:
		if t.c != nil {
			t.complete()
		} else {
			t.putc(c)
		}
	case opDeleteToEnd:
		t.deleteToEnd()
	case opClear:
		t.clearScreen()
	case opSubmit:
		return false, nil
	case opTranspose:
		t.transpose()
	case opDeleteToBeginning:
		t.deleteToBeginning()
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
					t.delete()
				}
				break
			}
			switch seq[1] {
			case 65:
				// down arrow
				currentEntry++
				if currentEntry > len(history) {
					currentEntry = 0
				}
			case 66:
				// up arrow
				currentEntry--
				if currentEntry < 0 {
					currentEntry = len(history) - 1
				}
			case 67:
				// right arrow
				t.right()
			case 68:
				// left arrow
				t.left()
			case 70:
				// end
				t.end()
			case 72:
				// home
				t.home()
			}
		} else if seq[0] == 79 {
			switch seq[1] {
			case 70:
				// end
				t.end()
			case 72:
				// home
				t.home()
			}
		}
	}
	return true, nil
}

func (t *term) putc(c int) {
	t.buf.WriteRune(c, t.pos)
	t.pos++
	t.refreshLine()
}

func (t *term) puts(s string) {
	t.buf.WriteString(s, t.pos)
	t.pos += len(s)
	t.refreshLine()
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

func (t *term) complete() {
	if t.display {
		t.printCandidates()
		return
	}
	str := t.buf.String()
	candidates := t.c.Complete(str, t.pos)
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
		t.display = true
		t.candidates = candidates
	}
	inter := findIntersect(str, complete)
	if inter == complete {
		return
	}
	t.puts(complete[len(inter):])
	t.refreshLine()
}

func (t *term) backspace() {
	if t.pos > 0 {
		t.pos--
		t.buf.remove(t.pos)
		t.refreshLine()
	}
}

// delete the character in front of the cursor, like the delete key
func (t *term) delete() {
	t.buf.remove(t.pos)
	t.refreshLine()
}

func (t *term) deleteToBeginning() {
	t.buf.pretruncate(t.pos)
	t.pos = 0
	t.refreshLine()
}

func (t *term) deleteToEnd() {
	t.buf.truncate(t.pos)
	t.refreshLine()
}

func (t *term) home() {
	t.pos = 0
	t.refreshLine()
}

func (t *term) end() {
	t.pos = t.buf.len()
	t.refreshLine()
}

func (t *term) transpose() {
	t.buf.transpose(t.pos)
	t.refreshLine()
}

// move the cursor left
func (t *term) left() {
	if t.pos > 0 {
		t.pos--
		t.refreshLine()
	}
}

// move the cursor right
func (t *term) right() {
	if t.pos < t.buf.len() {
		t.pos++
		t.refreshLine()
	}
}

func (t *term) refreshLine() {
	// move to origin of the current line
	t.setCursor(0, -t.y)
	// assuming the prompt won't wrap
	fmt.Print(t.prompt)
	bufStr := t.buf.String()
	n := len(bufStr)
	pl := len(t.prompt)
	if n > t.cols - pl {
		n = t.cols - pl
	}
	fmt.Print(bufStr[:n])
	bufStr = bufStr[n:]
	t.lines = 0
	wrapCursor := n == t.cols - pl
	n = len(bufStr)
	for n > 0 {
		if n > t.cols {
			n = t.cols
		}
		fmt.Print(bufStr[:n])
		bufStr = bufStr[n:]
		wrapCursor = n == t.cols
		n = len(bufStr)
		t.lines++
	}
	t.eraseToEnd()
	if wrapCursor {
		t.lines++
		// move to next line
		fmt.Print("\n")
	}
	x := (pl + t.pos) % t.cols
	t.y = (pl + t.pos) / t.cols
	t.setCursor(x, t.y - t.lines)
}
