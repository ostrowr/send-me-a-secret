package main

import (
	"fmt"

	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
	"github.com/ostrowr/send-me-a-secret/internal/utils"
)

func decrypt(privateKeyPassword, base64EncodedCiphertext string) error {
	if base64EncodedCiphertext == "" {
		utils.PrintDefaultf("Waiting for input from stdin... press ctrl-d when done\n")
		messageBytes, err := utils.GetMessageFromStdin()
		utils.FatallyLogOnError("Couldn't read from stdin", err)
		base64EncodedCiphertext = string(messageBytes)
	}

	if privateKeyPassword == "" {
		passwordBytes, err := utils.ReadPassword("Enter passphrase for private key: ")
		utils.FatallyLogOnError("Unable to read passphrase", err)
		privateKeyPassword = string(passwordBytes)
	}

	privateKey, err := rsahelpers.ReadPrivateKeyFromFile([]byte(privateKeyPassword))
	utils.FatallyLogOnError("Couldn't read private key", err)
	decrypted, err := rsahelpers.Decrypt(privateKey, base64EncodedCiphertext)
	utils.FatallyLogOnError("Couldn't decrypt message", err)
	fmt.Println(string(decrypted))
	return nil
}
