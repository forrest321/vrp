package debug

import (
	"fmt"
)

var Debug = false

func Out(format string, args ...interface{}) {
	if Debug {
		fmt.Printf(format+"\n", args...)
	}
}

func SetDebug(debug bool) {
	Debug = debug
}
