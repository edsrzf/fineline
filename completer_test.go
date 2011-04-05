package fineline

import (
	"testing"
)

var simpleList = []string {
	"cat",
	"catch",
	"cats",
	"caught",
	"cough",
	"dog",
}

var simpleTests = []struct {
	input string
	expected []string
}{
	{"", []string{"cat", "catch", "cats", "caught", "cough", "dog"}},
	{"c", []string{"cat", "catch", "cats", "caught", "cough"}},
	{"d", []string{"dog"}},
	{"e", nil},
	{"ca", []string{"cat", "catch", "cats", "caught"}},
	{"co", []string{"cough"}},
	{"do", []string{"dog"}},
	{"cat", []string{"cat", "catch", "cats"}},
	{"dog", []string{"dog"}},
	{"cat ", []string{"cat", "catch", "cats", "caught", "cough", "dog"}},
	{"dog ca", []string{"cat", "catch", "cats", "caught"}},
}

func listsEqual(x, y []string) bool {
	if len(x) != len(y) {
		return false
	}
	for i := range x {
		if x[i] != y[i] {
			return false
		}
	}
	return true
}

func TestSimpleCompleter(t *testing.T) {
	c := NewSimpleCompleter(simpleList)
	for _, test := range simpleTests {
		list := c.Complete(test.input, len(test.input))
		if !listsEqual(test.expected, list) {
			t.Errorf("expected list\n%v\ngot list\n%v\n", test.expected, list)
		}
	}
}
