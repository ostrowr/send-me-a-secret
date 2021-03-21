package main

import (
	"syscall/js"

	"github.com/ostrowr/send-me-a-secret/internal/encryptor"
)

func encryptWrapper() js.Func {
	f := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		if len(args) != 1 {
			return "No message provided."
		}
		message := args[0].String()
		encrypted, err := encryptor.Encrypt([]byte(message))
		if err != nil {
			return err.Error()
		}
		return encrypted
	})
	return f
}

func main() {
	js.Global().Set("encrypt", encryptWrapper())
	<-make(chan bool) // never exit
}
