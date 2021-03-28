package main

import (
	"syscall/js"

	"github.com/ostrowr/send-me-a-secret/internal/rsahelpers"
)

func encryptWrapper() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 2 {
			return "Must provide message and public key."
		}
		message := args[0].String()
		publicKeyBytes := []byte(args[1].String())
		publicKey, err := rsahelpers.SSHPubKeyToRSAPubKey(publicKeyBytes)
		if err != nil {
			return err.Error()
		}
		encrypted, err := rsahelpers.Encrypt(publicKey, []byte(message))
		if err != nil {
			return err.Error()
		}
		return encrypted
	})
	return f
}

func getValidPublicKeyWrapper() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "Must provide list of public keys."
		}

		var validKey string
		for i := 0; i < args[0].Length(); i++ {
			key := args[0].Index(i).String()
			parsed, err := rsahelpers.SSHPubKeyToRSAPubKey([]byte(key))
			if err != nil {
				continue // it's ok if they have some public keys that we can't parse
			}
			if rsahelpers.IsValidSendMeASecretKey(parsed) {
				if validKey != "" {
					return "Too many valid keys. Can't decide."
				}
				validKey = key
			}
		}
		if validKey == "" {
			return "No valid keys"
		}
		return validKey
	})
	return f
}

func main() {
	js.Global().Set("encrypt", encryptWrapper())
	js.Global().Set("getValidPublicKey", getValidPublicKeyWrapper())
	<-make(chan bool) // never exit
}
