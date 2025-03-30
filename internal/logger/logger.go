package logger

import (
	"fmt"
	"strings"
)

const (
	ColorGreen  = "\033[0;32m"
	ColorRed    = "\033[0;31m"
	ColorBlue   = "\033[0;34m"
	ColorYellow = "\033[0;33m"
	ColorReset  = "\033[0m"
)

func PrintHeader(command string) {
	separator := fmt.Sprintf("%s%s%s", ColorBlue, strings.Repeat("=", 50), ColorReset)
	fmt.Printf("%s\n", separator)
	fmt.Printf("%s🕯️ %s%s\n", ColorBlue, command, ColorReset)
	fmt.Printf("%s\n", separator)
	fmt.Println()
}

func Info(format string, args ...any) {
	fmt.Printf("%s[INFO] %s%s\n", ColorBlue, fmt.Sprintf(format, args...), ColorReset)
}

func Success(format string, args ...any) {
	fmt.Printf("%s[SUCCESS] %s%s\n", ColorGreen, fmt.Sprintf(format, args...), ColorReset)
}

func Error(format string, args ...any) {
	fmt.Printf("%s[ERROR] %s%s\n", ColorRed, fmt.Sprintf(format, args...), ColorReset)
}
