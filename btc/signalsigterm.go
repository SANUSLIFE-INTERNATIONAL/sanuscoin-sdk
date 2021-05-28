// Copyright © 2021 The Sanuscoin Team

// +build darwin dragonfly freebsd linux netbsd openbsd solaris

package btc

import (
	"os"
	"syscall"
)

func init() {
	interruptSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
}
