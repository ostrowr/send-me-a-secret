package utils

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"syscall"

	"golang.org/x/term"
)

func GetMessageFromStdin() []byte {
	// todo print a nice message if in terminal,
	// like "paste encrypted message here"
	message := make([]byte, 0)
	reader := bufio.NewReader(os.Stdin)
	for {
		b, err := reader.ReadByte()
		if err != nil {
			if err == io.EOF {
				return message
			}
			panic(err) // todo
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
		defer tty.Close()
		fd = int(tty.Fd())
	}

	pass, err := term.ReadPassword(fd)
	fmt.Fprintln(os.Stderr)
	return pass, err
}
