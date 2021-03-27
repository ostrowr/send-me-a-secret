package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"

	"github.com/fatih/color"

	"golang.org/x/term"
)

func GetMessageFromStdin() ([]byte, error) {
	message := make([]byte, 0)
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return message, nil
			}
			return nil, err
		}
		message = append(message, b)
	}
}

func ReadPassword(prompt string) ([]byte, error) {
	fmt.Fprint(os.Stderr, prompt)
	var fd int
	if term.IsTerminal(syscall.Stdin) {
		fd = syscall.Stdin
	} else {
		tty, err := os.Open("/dev/tty")
		if err != nil {
			return nil, err
		}
		defer MustClose(tty)
		fd = int(tty.Fd())
	}

	pass, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return pass, err
}

func FatallyLogOnError(message string, err error) {
	if err != nil {
		color.Set(color.FgRed)
		defer color.Unset()
		log.Fatal(message, "\nError: ", err)
	}
}

func MustClose(f *os.File) {
	err := f.Close()
	FatallyLogOnError("Error closing file", err)
}

func PrintRed(format string, a ...interface{}) {
	color.New(color.FgRed).Fprintf(os.Stderr, format, a...)
}

func PrintCyan(format string, a ...interface{}) {
	color.New(color.FgCyan).Fprintf(os.Stderr, format, a...)
}

func PrintGreen(format string, a ...interface{}) {
	color.New(color.FgGreen).Fprintf(os.Stderr, format, a...)
}

func PrintYellow(format string, a ...interface{}) {
	color.New(color.FgYellow).Fprintf(os.Stderr, format, a...)
}

func PrintDefault(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, format, a...)
}
