// Created by cgo -godefs - DO NOT EDIT
// cgo -godefs defs_posix.go

package fineline

const (
    TIOCGWINSZ  = 0x40087468

    IGNBRK  = 0x1
    BRKINT  = 0x2
    IGNPAR  = 0x4
    PARMRK  = 0x8
    INPCK   = 0x10
    ISTRIP  = 0x20
    INLCR   = 0x40
    IGNCR   = 0x80
    ICRNL   = 0x100
    IXON    = 0x200
    IXANY   = 0x800
    IXOFF   = 0x400
    IMAXBEL = 0x2000
    IUTF8   = 0x4000

    OPOST   = 0x1
    ONLCR   = 0x2
    OCRNL   = 0x10
    ONOCR   = 0x20
    ONLRET  = 0x40
    OFILL   = 0x80
    OFDEL   = 0x20000

    CS8 = 0x300

    ISIG    = 0x80
    ICANON  = 0x100
    ECHO    = 0x8
    ECHOE   = 0x2
    ECHOK   = 0x4
    ECHONL  = 0x10
    NOFLSH  = 0x80000000
    TOSTOP  = 0x400000
    IEXTEN  = 0x400

    VINTR       = 0x8
    VQUIT       = 0x9
    VERASE      = 0x3
    VKILL       = 0x5
    VEOF        = 0x0
    VTIME       = 0x11
    VMIN        = 0x10
    VSTART      = 0xc
    VSTOP       = 0xd
    VSUSP       = 0xa
    VEOL        = 0x1
    VREPRINT    = 0x6
    VDISCARD    = 0xf
    VWERASE     = 0x4
    VLNEXT      = 0xe
    VEOL2       = 0x2

    TCSANOW     = 0x0
    TCSADRAIN   = 0x1
    TCSAFLUSH   = 0x2
)

type termios struct {
    Iflag   uint64
    Oflag   uint64
    Cflag   uint64
    Lflag   uint64
    Cc  [20]uint8
    Pad_cgo_0   [4]byte
    Ispeed  uint64
    Ospeed  uint64
}
type winsize struct {
    Row uint16
    Col uint16
    Xpixel  uint16
    Ypixel  uint16
}
