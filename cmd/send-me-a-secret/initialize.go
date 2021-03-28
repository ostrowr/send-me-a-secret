package main

import (
	"fmt"

	"github.com/ostrowr/send-me-a-secret/internal/githubapi"
	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
	"github.com/ostrowr/send-me-a-secret/internal/utils"
)

func initialize(skipVerify bool, skipGithub bool, password string, githubToken string) error {
	utils.PrintDefault("Initializing send-me-a-secret\n")
	if password == "" {
		passwordBytes, err := utils.ReadPassword("Enter passphrase (empty for no passphrase): ")
		// todo confirm password
		utils.FatallyLogOnError("Unable to read passphrase", err)
		password = string(passwordBytes)
	}

	utils.PrintDefault("Generating new private key (this may take a moment)...\n")
	privateKey, err := rsahelpers.GenerateKey()
	utils.FatallyLogOnError("Couldn't generate private key", err)
	utils.PrintCyan("Private key generated. It's saved in %s\n\n", rsahelpers.PathToKeyFile())
	// todo alert before overwriting
	err = rsahelpers.WritePrivateKeyToFile([]byte(password), privateKey)
	utils.FatallyLogOnError("Couldn't write private key", err)
	publicKeyBytes, err := rsahelpers.GetSSHPublicKey(privateKey)
	utils.FatallyLogOnError("Couldn't generate public key", err)
	publicKey := string(publicKeyBytes)
	if skipGithub {
		utils.PrintYellow("--skip-github specified; not uploading key to GitHub. Copy the following key into https://github.com/settings/ssh/new with the title %s\n\n", githubapi.PublicKeyName)
		fmt.Println(publicKey)
		return nil
	}

	if githubToken == "" {
		utils.PrintYellow("send-me-a-secret adds a new public key to your GitHub account. If you'd prefer to manually add the key to your account, re-run `initialize` with the --skip-github flag.\n")
		githubTokenBytes, err := utils.ReadPassword("Enter a GitHub token with write:public_key access. (You can generate one at https://github.com/settings/tokens/new): ")
		utils.FatallyLogOnError("Unable to read GitHub token", err)
		githubToken = string(githubTokenBytes)
	}

	client := githubapi.GetGithubClient(githubToken)

	err = githubapi.UploadPublicKeyToGithub(client, publicKey)
	utils.FatallyLogOnError("Couldn't upload public key to GitHub", err)
	utils.PrintCyan("Public key uploaded.\n\n")
	if skipVerify {
		return nil
	}

	utils.PrintDefault("Getting authenticated GitHub username...\n")
	username, err := githubapi.GetCurrentUsername(client)
	utils.FatallyLogOnError("Couldn't get GitHub username\n", err)
	utils.PrintCyan("Got username: %s.\n\n", username)

	validate(username, password)
	return nil
}
