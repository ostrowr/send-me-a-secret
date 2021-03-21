package main

import (
	"fmt"

	"github.com/ostrowr/send-me-a-secret/internal/encryptor"
	"github.com/ostrowr/send-me-a-secret/internal/utils"
)

func main() {
	message, err := utils.GetMessageFromStdin()
	utils.FatallyLogOnError("Could not read message", err)
	ciphertext, err := encryptor.Encrypt(message)
	utils.FatallyLogOnError("Failed to encrypt", err)
	fmt.Println(ciphertext)
}
