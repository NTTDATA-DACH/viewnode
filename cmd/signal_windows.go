//go:build windows

package cmd

import "os"

func watchSignals() []os.Signal {
	return []os.Signal{os.Interrupt}
}
