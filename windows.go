package fineline

import (
	"fmt"
	"os"
	"strings"
	"syscall"
	"unsafe"
)

type term struct {
	termCommon
	// platform-specific
	// console handle
	h uintptr
	rows int
}

var (
	modkernel32 = loadDll("kernel32.dll")

	procFillConsoleOutputCharacter = getSysProcAddr(modkernel32, "FillConsoleOutputCharacterW")
	procGetConsoleMode = getSysProcAddr(modkernel32, "GetConsoleMode")
	procGetConsoleScreenBufferInfo = getSysProcAddr(modkernel32, "GetConsoleScreenBufferInfo")
	procSetConsoleCursorPosition = getSysProcAddr(modkernel32, "SetConsoleCursorPosition")
	procSetConsoleMode = getSysProcAddr(modkernel32, "SetConsoleMode")
)

func loadDll(fname string) uint32 {
	h, e := syscall.LoadLibrary(fname)
	if e != 0 {
		panic("LoadLibrary failed")
	}
	return h
}

func getSysProcAddr(m uint32, pname string) uintptr {
	p, e := syscall.GetProcAddress(m, pname)
	if e != 0 {
		panic("GetProcAddress failed on " + pname)
	}
	return uintptr(p)
}

const (
	_ENABLE_PROCESSED_INPUT = 0x0001
	_ENABLE_LINE_INPUT = 0x0002
	_ENABLE_ECHO_INPUT = 0x0004
	_ENABLE_WINDOW_INPUT = 0x0008
	_ENABLE_MOUSE_INPUT = 0x0010
	_ENABLE_INSERT_MODE = 0x0020
	_ENABLE_QUICK_EDIT_MODE = 0x0040
	_ENABLE_EXTENDED_FLAGS = 0x0080
)

type coord struct {
	x, y int16
}

type smallRect struct {
	left, top, right, bottom int16
}

type consoleScreenBufferInfo struct {
	dwSize coord
	dwCursorPosition coord
	wAttributes uint16
	srWindow smallRect
	dwMaximumWindowSize coord
}

var origTerm uint32

// enable raw mode and gather metrics, like number of columns
func (t *term) init() {
	// STD_OUTPUT_HANDLE
	h, errno := syscall.GetStdHandle(-11)
	t.h = uintptr(h)
	if int32(t.h) == -1 {
		err := os.Errno(errno)
		panic(err)
	}
	ok, _, e := syscall.Syscall(procGetConsoleMode, 2,
		t.h, uintptr(unsafe.Pointer(&origTerm)), 0)
	if ok == 0 {
		err := os.NewSyscallError("GetConsoleMode", int(e))
		panic(err)
	}

	raw := origTerm
	raw &^= _ENABLE_LINE_INPUT | _ENABLE_ECHO_INPUT | _ENABLE_PROCESSED_INPUT | _ENABLE_WINDOW_INPUT
	ok, _, e = syscall.Syscall(procSetConsoleMode, 2, t.h, uintptr(raw), 0)
	if ok == 0 {
		err := os.NewSyscallError("SetConsoleMode", int(e))
		panic(err)
	}

	win := t.getConsoleInfo()
	t.cols = int(win.dwSize.x)
	t.rows = int(win.dwSize.y)

	t.buf = new(buffer)
}

func (t *term) disableRawMode() {
	ok, _, e := syscall.Syscall(procSetConsoleMode, 2, t.h, uintptr(origTerm), 0)
	if ok == 0 {
		err := os.NewSyscallError("SetConsoleMode", int(e))
		panic(err)
	}
}

// x is absolute, y is relative
func (t *term) setCursor(x, y int) {
	pos := t.getCursorPos()
	pos.x = int16(x)
	pos.y += int16(y)
	t.setCursorPos(pos)
}

func (t *term) eraseToEnd() {
	pos := t.getCursorPos()
	length := t.cols*(t.rows - int(pos.y) + 1) - int(pos.x)
	t.fillConsoleOutputCharacter(' ' << 8, length, pos)
}

func (t *term) printCandidates() {
	str := "\n\x1b[0G" + strings.Join(t.candidates, "\n\x1b[0G") + "\n"
	fmt.Print(str)
	t.refreshLine()
}

func (t *term) clearScreen() {
	pos := coord{}
	t.fillConsoleOutputCharacter(' ' << 8, t.cols*t.rows, pos)
	t.setCursorPos(pos)
}

func (t *term) getCursorPos() coord {
	win := t.getConsoleInfo()
	return win.dwCursorPosition
}

func (t *term) setCursorPos(pos coord) {
	ok, _, e := syscall.Syscall(procSetConsoleCursorPosition, 2, t.h,
		*(*uintptr)(unsafe.Pointer(&pos)), 0)
	if ok == 0 {
		err := os.NewSyscallError("SetConsoleCursorPosition", int(e))
		panic(err)
	}
}

func (t *term) getConsoleInfo() *consoleScreenBufferInfo {
	var win consoleScreenBufferInfo
	ok, _, e := syscall.Syscall(procGetConsoleScreenBufferInfo, 2, t.h,
		uintptr(unsafe.Pointer(&win)), 0)
	if ok == 0 {
		err := os.NewSyscallError("GetConsoleScreenBufferInfo", int(e))
		panic(err)
	}
	return &win
}

func (t *term) fillConsoleOutputCharacter(char uint16, length int, pos coord) {
	coord := uintptr(unsafe.Pointer(&pos))
	out := uintptr(unsafe.Pointer(new(int)))
	ok, _, e := syscall.Syscall6(procFillConsoleOutputCharacter, 5,
		t.h, uintptr(char), uintptr(length), coord, out, 0)
	if ok == 0 {
		err := os.NewSyscallError("FillConsoleOutputCharacter", int(e))
		panic(err)
	}
}
