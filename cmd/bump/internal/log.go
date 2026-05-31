package internal

import (
	"fmt"
	"os"
	"strings"
)

func Debug(format string, args ...any) {
	if DebugFlag != nil && QuietFlag != nil && *DebugFlag && !*QuietFlag {
		fmt.Printf(format, args...)
	}
}

func Info(format string, args ...any) {
	if QuietFlag != nil && !*QuietFlag {
		fmt.Printf(format, args...)
	}
}

func Error(format string, args ...any) {
	fmt.Fprintf(os.Stderr, format, args...)
}

func SliceString(s []string) string {
	first := true
	b := strings.Builder{}
	for _, s := range s {
		if first {
			first = false
		} else {
			b.WriteString(", ")
		}
		b.WriteString(s)
	}
	return b.String()
}
