package tools

import (
	"fmt"
	"os"
)

func LogErrorAndExit(e error, debug bool) {
	if debug {
		fmt.Printf("%+v", e)
	} else {
		fmt.Println(e.Error())
	}
	os.Exit(1)
}
