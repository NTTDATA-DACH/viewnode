//go:build !windows

package cmd

import (
	"os"
	"syscall"
)

func watchSignals() []os.Signal {
	return []os.Signal{os.Interrupt, syscall.SIGTERM}
}
