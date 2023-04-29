package primitive

import "fmt"

// LogLevel is the loglevel for the application
var LogLevel int

// Log outputs a log entry to standard output
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
