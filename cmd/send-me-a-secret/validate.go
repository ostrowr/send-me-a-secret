package main

import (
	"errors"

	"github.com/ostrowr/send-me-a-secret/internal/githubapi"
	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
	"github.com/ostrowr/send-me-a-secret/internal/utils"
)

func validate(githubUsername, privateKeyPassword string) error {
	client := githubapi.GetGithubClient("")
	utils.PrintDefaultf("Validating send-me-a-secret\n")

	utils.PrintDefaultf("Fetching public key from GitHub...\n")
	publicKey, err := githubapi.GetPublicKeyFromGithubUnauthenticated(client, githubUsername, rsahelpers.IsValidSendMeASecretKey)
	utils.FatallyLogOnError("Couldn't fetch public key from GitHub", err)
	utils.PrintCyanf("Public key successfully fetched\n\n")
	utils.PrintDefaultf("Encrypting a test message using that public key\n")
	testMessage := "Hello, my name is Inigo Montoya."
	ciphertext, err := rsahelpers.Encrypt(publicKey, []byte(testMessage))
	utils.FatallyLogOnError("Couldn't encrypt message", err)
	utils.PrintCyanf("Message successfully encrypted\n\n")
	utils.PrintDefaultf("Decrypting the test message using your private key\n")
	if privateKeyPassword == "" {
		var passwordBytes []byte
		passwordBytes, err = utils.ReadPassword("Enter passphrase for private key: ")
		utils.FatallyLogOnError("Unable to read passphrase", err)
		privateKeyPassword = string(passwordBytes)
	}

	privateKey, err := rsahelpers.ReadPrivateKeyFromFile([]byte(privateKeyPassword))
	utils.FatallyLogOnError("Couldn't read private key", err)
	decrypted, err := rsahelpers.Decrypt(privateKey, ciphertext)
	utils.FatallyLogOnError("Couldn't decrypt message", err)
	if string(decrypted) != testMessage {
		utils.FatallyLogOnError("", errors.New("messages don't match"))
	}
	utils.PrintCyanf("Message successfully decrypted\n\n")
	utils.PrintGreenf("Validation succeeded! You're good to go.\n")
	return nil
}
