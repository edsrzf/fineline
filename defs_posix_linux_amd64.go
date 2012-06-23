// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_posix.go

package fineline

const (
	TIOCGWINSZ	= 0x5413

	IGNBRK	= 0x1
	BRKINT	= 0x2
	IGNPAR	= 0x4
	PARMRK	= 0x8
	INPCK	= 0x10
	ISTRIP	= 0x20
	INLCR	= 0x40
	IGNCR	= 0x80
	ICRNL	= 0x100
	IUCLC	= 0x200
	IXON	= 0x400
	IXANY	= 0x800
	IXOFF	= 0x1000
	IMAXBEL	= 0x2000
	IUTF8	= 0x4000

	OPOST	= 0x1
	OLCUC	= 0x2
	ONLCR	= 0x4
	OCRNL	= 0x8
	ONOCR	= 0x10
	ONLRET	= 0x20
	OFILL	= 0x40
	OFDEL	= 0x80

	CS8	= 0x30

	ISIG	= 0x1
	ICANON	= 0x2
	ECHO	= 0x8
	ECHOE	= 0x10
	ECHOK	= 0x20
	ECHONL	= 0x40
	NOFLSH	= 0x80
	TOSTOP	= 0x100
	IEXTEN	= 0x8000

	VINTR		= 0x0
	VQUIT		= 0x1
	VERASE		= 0x2
	VKILL		= 0x3
	VEOF		= 0x4
	VTIME		= 0x5
	VMIN		= 0x6
	VSWTC		= 0x7
	VSTART		= 0x8
	VSTOP		= 0x9
	VSUSP		= 0xa
	VEOL		= 0xb
	VREPRINT	= 0xc
	VDISCARD	= 0xd
	VWERASE		= 0xe
	VLNEXT		= 0xf
	VEOL2		= 0x10

	TCSANOW		= 0x0
	TCSADRAIN	= 0x1
	TCSAFLUSH	= 0x2
)

type termios struct {
	Iflag	uint32
	Oflag	uint32
	Cflag	uint32
	Lflag	uint32
	Line	uint8
	Cc	[32]uint8
	Pad_cgo_0	[3]byte
	Ispeed	uint32
	Ospeed	uint32
}
type winsize struct {
	Row	uint16
	Col	uint16
	Xpixel	uint16
	Ypixel	uint16
}
