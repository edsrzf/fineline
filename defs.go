// +build ignore

package fineline

/*
#include <termios.h>
#include <sys/ioctl.h>
*/
import "C"

const (
	TCGETS     = C.TCGETS
	TCSETS     = C.TCSETS
	TCSETSW    = C.TCSETSW
	TCSETSF    = C.TCSETSF
	TIOCGWINSZ = C.TIOCGWINSZ

	/* c_iflag bits */
	IGNBRK  = C.IGNBRK
	BRKINT  = C.BRKINT
	IGNPAR  = C.IGNPAR
	PARMRK  = C.PARMRK
	INPCK   = C.INPCK
	ISTRIP  = C.ISTRIP
	INLCR   = C.INLCR
	IGNCR   = C.IGNCR
	ICRNL   = C.ICRNL
	IUCLC   = C.IUCLC
	IXON    = C.IXON
	IXANY   = C.IXANY
	IXOFF   = C.IXOFF
	IMAXBEL = C.IMAXBEL
	IUTF8   = C.IUTF8

	/* c_oflag bits */
	OPOST  = C.OPOST
	OLCUC  = C.OLCUC
	ONLCR  = C.ONLCR
	OCRNL  = C.OCRNL
	ONOCR  = C.ONOCR
	ONLRET = C.ONLRET
	OFILL  = C.OFILL
	OFDEL  = C.OFDEL

	CS8 = C.CS8

	/* c_lflag bits */
	ISIG   = C.ISIG
	ICANON = C.ICANON
	ECHO   = C.ECHO
	ECHOE  = C.ECHOE
	ECHOK  = C.ECHOK
	ECHONL = C.ECHONL
	NOFLSH = C.NOFLSH
	TOSTOP = C.TOSTOP
	IEXTEN = C.IEXTEN

	/* c_cc characters */
	VINTR    = C.VINTR
	VQUIT    = C.VQUIT
	VERASE   = C.VERASE
	VKILL    = C.VKILL
	VEOF     = C.VEOF
	VTIME    = C.VTIME
	VMIN     = C.VMIN
	VSWTC    = C.VSWTC
	VSTART   = C.VSTART
	VSTOP    = C.VSTOP
	VSUSP    = C.VSUSP
	VEOL     = C.VEOL
	VREPRINT = C.VREPRINT
	VDISCARD = C.VDISCARD
	VWERASE  = C.VWERASE
	VLNEXT   = C.VLNEXT
	VEOL2    = C.VEOL2

	/* tcsetattr uses these */
	TCSANOW   = C.TCSANOW
	TCSADRAIN = C.TCSADRAIN
	TCSAFLUSH = C.TCSAFLUSH
)

type termios C.struct_termios
type winsize C.struct_winsize
