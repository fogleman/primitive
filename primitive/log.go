package primitive

import "fmt"

var LogLevel int

func Log(level int, format string, a ...interface{}) {
	if LogLevel >= level {
		fmt.Printf(format, a...)
	}
}

func v(format string, a ...interface{}) {
	Log(1, format, a...)
}

func vv(format string, a ...interface{}) {
	Log(2, "  "+format, a...)
}

func vvv(format string, a ...interface{}) {
	Log(3, "    "+format, a...)
}
