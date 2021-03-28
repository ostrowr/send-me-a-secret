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

func getPublicKeyWrapper() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "No message provided."
		}

		// 1. Get a list of all public keys (js-land)
		// 2. Pass that list into this function, which figures out which (if any) are valid
		// 3. Return the (string) public key back from this function
		// 4. Run encrypt using the public key returned from this function, parsed yet again.

		return `ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAACPACro6jmd+F8ZkcS2lmtRsIiuEUwjXDyxr0ZF1U9fBypwF2LzZTqQh0WykHcn0ETycRUonSL8feTLxaPCjv0puUdS7vY16LZsFDDyr4siRfcJFqkD8psf55Fm5bZYvnLkUgUil4dwhk3pKS59OzZrl5su0j01OUf0Ly1WzdiONNtr2U+fNrP5PB3mHCgyhD3rVMEAtVAahbah7PF1aSoCqs4Xazhg5yaGTUiqBREYhqXbQsLNZL4vxLURFF95Ngqn9srWN6neDmyjuWVrH9yKrEnsGYyWcvWdATLj/PslSVsdZ9RKbchxlEoNQVECPV7KVhkaCBGOnQ+tg++7svKhM/NacgoJE/qOnAWMK2mjBEOeuPEgA9qQibD1Ps0XwkHDDyMa9HHCNWTme+s67RUy84mc7wImDJYFn/hTKQgRdY6gTJmOhP/hBsBhbU1U/DYOUXtUzUoE+2eIpewSYshwdM9F7xGwVH+F1ArcUej67Fxl5Fty7talvIyJSYzr10ZIgBxxWijDzZ+wlIn7jwxHIeaZ3ilVPqloBJN/G9zDJZArwm59ALJ7JARDU+FQxXwHxO+TR2AhIGdMg3IF0JVpmQB0rXu/7bip2JraM936DM02494wJj15fxCWBacIugp895gqI8MdacUKQHlCmuDzqXuy67wpI2mokm339rQpG0uyO/U3Ese7i47IdZfHP5JTGNGv7IKAxGtY2Cvmsz/sCDa9ZRED+EzKAMTDTfqF/8HdhsR3WEVsBIwOOX9`
	})
	return f
}

func main() {
	js.Global().Set("encrypt", encryptWrapper())
	js.Global().Set("getPublicKey", getPublicKeyWrapper())
	<-make(chan bool) // never exit
}
