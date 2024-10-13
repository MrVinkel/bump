package bump

import (
	"fmt"
	"os"
	"strings"
)

func Debug(format string, args ...interface{}) {
	if *DebugFlag && !*QuietFlag {
		fmt.Printf(format, args...)
	}
}

func Info(format string, args ...interface{}) {
	if !*QuietFlag {
		fmt.Printf(format, args...)
	}
}

func Error(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func PrintSlice(slice []string) string {
	first := true
	b := strings.Builder{}
	for _, s := range slice {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}
		b.WriteString(s)
	}
	return b.String()
}
