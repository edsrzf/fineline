// +build ignore

package fineline

/*
#include <sys/ioctl.h>
*/
import "C"

const (
	TCSETSF    = C.TCSETSF
	TCSETSW    = C.TCSETSW
)
