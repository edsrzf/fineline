package fineline

import (
	"syscall"
)

func tcgetattr(fd int, t *termios) {
	ttyIoctl(0, syscall.TCGETS, t)
}

func tcsetattr(fd, op int, t *termios) {
	var cmd int
	switch op {
	case TCSANOW:
		cmd = syscall.TCSETS
	case TCSADRAIN:
		cmd = TCSETSW
	case TCSAFLUSH:
		cmd = TCSETSF
	}
	ttyIoctl(0, cmd, t)
}
