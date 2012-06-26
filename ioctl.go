package fineline

import (
	"syscall"
	"unsafe"
)

func winIoctl(fd int, cmd uintptr, win *winsize) {
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), cmd, uintptr(unsafe.Pointer(win)))
	return
}

func ttyIoctl(fd int, cmd uintptr, term *termios) {
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), cmd, uintptr(unsafe.Pointer(term)))
	return
}
