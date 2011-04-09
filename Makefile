# Copyright 2011 The Go Authors.  All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

include $(GOROOT)/src/Make.inc

TARG=github.com/edsrzf/fineline

GOFILES=\
	buffer.go\
	completer.go\
	fineline.go\
	ops.go\
	termios.go\

GOFILES_linux=\
	ansi.go\
	ioctl.go\

GOFILES_windows=\
	windows.go\

GOFILES+=$(GOFILES_$(GOOS))

include $(GOROOT)/src/Make.pkg

termios.go: defs.c
	godefs -g fineline defs.c > termios.go
