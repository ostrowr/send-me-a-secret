package utils

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"syscall"

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
		log.Fatal(message, "\nError: ", err)
	}
}

func MustClose(f *os.File) {
	err := f.Close()
	FatallyLogOnError("Error closing file", err)
}
