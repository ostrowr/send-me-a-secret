package main

import (
	"fmt"

	"github.com/ostrowr/send-me-a-secret/internal/githubapi"
	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
	"github.com/ostrowr/send-me-a-secret/internal/utils"
)

func encrypt(githubUsername string, message string) error {
	if message == "" {
		utils.PrintDefault("Waiting for input from stdin... press ctrl-d when done\n")
		messageBytes, err := utils.GetMessageFromStdin()
		utils.FatallyLogOnError("Couldn't read from stdin", err)
		message = string(messageBytes)
	}

	utils.PrintDefault("Getting public key from GitHub...\n")
	client := githubapi.GetGithubClient("")
	publicKey, err := githubapi.GetPublicKeyFromGithubUnauthenticated(client, githubUsername, rsahelpers.IsValidSendMeASecretKey)
	utils.FatallyLogOnError("Couldn't fetch public key from GitHub", err)
	ciphertext, err := rsahelpers.Encrypt(publicKey, []byte(message))
	utils.FatallyLogOnError("Couldn't encrypt message", err)
	utils.PrintGreen("Encryption succeeded! Ciphertext below:\n\n")
	fmt.Println(ciphertext)
	return nil
}
