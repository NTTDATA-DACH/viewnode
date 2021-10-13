package tools

import (
	"fmt"
	"os"
)

func LogDebug(msg string, debug bool) {
	if debug {
		fmt.Printf("debug: %s\n", msg)
	}
}

func LogErrorAndExit(e error, debug bool) {
	if debug {
		fmt.Printf("%+v", e)
	} else {
		fmt.Println(e.Error())
	}
	os.Exit(1)
}
