package output

import (
	"fmt"
	"io"
	"os"

	"github.com/fatih/color"
)

var (
	writer io.Writer = os.Stdout
	green            = color.New(color.FgGreen).SprintFunc()
	cyan             = color.New(color.FgCyan).SprintFunc()
	yellow           = color.New(color.FgYellow).SprintFunc()
	red              = color.New(color.FgRed).SprintFunc()
)

// SetWriter redirects all output to the given writer (e.g. os.Stderr for MCP mode).
func SetWriter(w io.Writer) {
	writer = w
}

func Success(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(writer, "%s %s\n", green("✓"), msg)
}

func Info(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(writer, "%s %s\n", cyan("→"), msg)
}

func Warning(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(writer, "%s %s\n", yellow("!"), msg)
}

func Error(format string, a ...interface{}) {
	msg := fmt.Sprintf(format, a...)
	fmt.Fprintf(writer, "%s %s\n", red("✗"), msg)
}

func Create(path string) {
	fmt.Fprintf(writer, "  %s %s\n", green("create"), path)
}

func Mkdir(path string) {
	fmt.Fprintf(writer, "  %s %s\n", green("mkdir"), path)
}

func Skip(path string) {
	fmt.Fprintf(writer, "  %s %s\n", yellow("skip"), path)
}

func Route(method, path string) {
	fmt.Fprintf(writer, "  %s %-7s %s\n", green("route"), method, path)
}

func Inject(file, description string) {
	fmt.Fprintf(writer, "  %s %s (%s)\n", cyan("inject"), file, description)
}
