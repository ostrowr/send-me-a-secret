package main

import (
	"encoding/base64"
	"fmt"

	"github.com/ostrowr/send-me-a-secret/internal/decryptor"
	"github.com/ostrowr/send-me-a-secret/internal/utils"
)

func main() {
	encodedCiphertext, err := utils.GetMessageFromStdin()
	utils.FatallyLogOnError("Could not read ciphertext", err)
	ciphertext, err := base64.StdEncoding.DecodeString(string(encodedCiphertext))
	utils.FatallyLogOnError("Could not decode base64 ciphertext", err)
	password, err := utils.ReadPassword("Enter password: ")
	utils.FatallyLogOnError("Could not read password", err)
	message, err := decryptor.Decrypt(ciphertext, password)
	utils.FatallyLogOnError("Could not parse decrypt message", err)
	fmt.Println(string(message))
}
