package fineline

import (
	"os"
	"path/filepath"
	"runtime"
	"sort"
	"strings"
)

// A Completer provides candidates for tab-completion.
type Completer interface {
	// Complete takes a string and a cursor position and returns
	// a list of candidate strings for tab-completion.
	Complete(str string, cur int) []string
}

// A SimpleCompleter provides completion candidates from a list
// of strings. Words are separated by Delim, which is a single
// space by default.
type SimpleCompleter struct {
	list []string
	Delim string
}

// NewSimpleCompleter creates a new SimpleCompleter.
// The list is sorted and used to provide completion candidates.
func NewSimpleCompleter(list []string) *SimpleCompleter {
	c := &SimpleCompleter{}
	c.SetList(list)
	c.Delim = " "
	return c
}

// SetList sorts a list of strings and supplies that list to c.
func (c *SimpleCompleter) SetList(list []string) {
	sort.Strings(list)
	c.list = list
}

// AddString inserts a string into c's list, using already allocated space
// if possible.
func (c *SimpleCompleter) AddString(str string) {
	n := len(c.list)
	pos := sort.Search(n, func(i int) bool { return c.list[i] >= str })
	c.list = append(c.list, str)
	if pos < n {
		copy(c.list[pos+1:], c.list[pos:])
		c.list[pos] = str
	}
}

// RemoveString removes a string from c's list. If the string is not in the
// list, RemoveString does nothing.
func (c *SimpleCompleter) RemoveString(str string) {
	n := len(c.list)
	pos := sort.Search(n, func(i int) bool { return c.list[i] >= str })
	if pos >= n {
		return
	}
	if pos < n - 1 {
		copy(c.list[pos:], c.list[pos+1:])
	}
	c.list = c.list[:n-1]
}

func (c *SimpleCompleter) Complete(str string, cur int) []string {
	// find the prefix
	tokStart := strings.LastIndex(str[:cur], c.Delim)
	tokEnd := strings.Index(str[cur:], c.Delim)
	if tokEnd < 0 {
		tokEnd = len(str) - cur
	}
	prefix := str[tokStart+1:cur+tokEnd]

	n := len(c.list)
	searchFunc := func(i int) bool { return c.list[i] >= prefix }
	first := sort.Search(n, searchFunc)
	if first == n || !strings.HasPrefix(c.list[first], prefix) {
		return nil
	}
	last := first + 1
	for last < n && strings.HasPrefix(c.list[last], prefix) {
		last++
	}
	return c.list[first:last]
}

// returns the user's home directory
func getHome() string {
	if runtime.GOOS == "windows" {
		return os.Getenv("USERPROFILE")
	}
	return os.Getenv("HOME")
}

// A FilenameCompleter completes a path string.
type FilenameCompleter struct {
	Delim string
}

func (c *FilenameCompleter) Complete(str string, cur int) []string {
	// find the prefix
	pos := strings.LastIndex(str[:cur], c.Delim)
	prefix := str[pos+1:cur]

	// four cases to consider:
	// 1. No characters
	// 2. First character is '/'
	// 3. First character is '~' and second is '/'
	// 4. First character is '~'
	var dirPath string
	n := len(prefix)
	if filepath.IsAbs(prefix) {
		// use the root directory
		prefix = prefix[1:]
		dirPath = "/"
	} else if n > 0 && prefix[0] == '~' {
		if n > 1 && prefix[1] == '/' {
			prefix = prefix[2:]
			dirPath = getHome()
		} else {
			// what to do?
			// parse /etc/passwd to get users (sigh)
			// for Windows, LsaLookupNames?
		}
	} else {
		// use current directory
		dirPath, _ = os.Getwd()
	}
	dir, err := os.Open(dirPath)
	if err != nil {
		panic(err.String())
	}
	defer dir.Close()
	var candidates []string
	names, err := dir.Readdir(-1)
	if err != nil {
		panic(err.String())
	}
	for _, f := range names {
		if strings.HasPrefix(f.Name, prefix) {
			if f.IsDirectory() {
				candidates = append(candidates, f.Name + "/")
			} else {
				candidates = append(candidates, f.Name)
			}
		}
	}
	return candidates
}

func completeString(str string, cur int, c Completer) {
	candidates := c.Complete(str, cur)
	switch len(candidates) {
	case 0:
		// do nothing
	case 1:
		// we found it
	default:
		// see if there's a common prefix longer than our current prefix
	}
}
