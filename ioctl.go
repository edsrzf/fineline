package fineline

import (
	"syscall"
	"unsafe"
)

func winIoctl(fd int, cmd int, win *winsize) {
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(win)))
	return
}

func ttyIoctl(fd int, cmd int, term *termios) {
	_, _, _ = syscall.Syscall(syscall.SYS_IOCTL, uintptr(fd), uintptr(cmd), uintptr(unsafe.Pointer(term)))
	return
}
