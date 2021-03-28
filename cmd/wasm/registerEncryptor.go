package main

import (
	"crypto/rsa"
	"fmt"
	"syscall/js"

	"github.com/ostrowr/send-me-a-secret/internal/githubapi"
	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
)

func encryptWrapper() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "No message provided."
		}
		pubKey := interface{}(args[0]).(*rsa.PublicKey)
		message := args[1].String()
		encrypted, err := rsahelpers.Encrypt(pubKey, []byte(message))
		if err != nil {
			return err.Error()
		}
		return encrypted
	})
	return f
}

func getPublicKeyWrapper() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "No message provided."
		}
		username := args[0].String()
		client := githubapi.GetGithubClient("")
		fmt.Println(client)
		publicKey, err := githubapi.GetPublicKeyFromGithubUnauthenticated(client, username, rsahelpers.IsValidSendMeASecretKey)
		if err != nil {
			return err.Error()
		}
		fmt.Println(publicKey)
		return username

		// return publicKey.E
	})
	return f
}

func main() {
	js.Global().Set("encrypt", encryptWrapper())
	js.Global().Set("getPublicKey", getPublicKeyWrapper())
	<-make(chan bool) // never exit
}
