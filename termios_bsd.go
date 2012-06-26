// +build darwin freebsd netbsd openbsd

package fineline

import (
	"syscall"
)

func tcgetattr(fd int, t *termios) {
	ttyIoctl(0, syscall.TIOCGETA, t)
}

func tcsetattr(fd, op int, t *termios) {
	var cmd uintptr
	switch op {
	case TCSANOW:
		cmd = syscall.TIOCSETA
	case TCSADRAIN:
		cmd = syscall.TIOCSETAW
	case TCSAFLUSH:
		cmd = syscall.TIOCSETAF
	}
	ttyIoctl(0, cmd, t)
}
