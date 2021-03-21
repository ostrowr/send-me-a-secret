package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"

	"github.com/ostrowr/send-me-a-secret/utils"
	"github.com/webview/webview"
)

func main() {
	publicKey, err := pemToPublicKey([]byte(publicKeyPem))
	if err != nil {
		panic(err) // todo
	}
	message := utils.GetMessageFromStdin()
	ciphertext, err := encrypt(publicKey, message)
	if err != nil {
		panic(err) // todo
	}
	fmt.Println(ciphertext)
}

func openWebview() {
	debug := true
	w := webview.New(debug)
	defer w.Destroy()
	w.SetTitle("Minimal webview example")
	w.SetSize(800, 600, webview.HintNone)
	w.Navigate("https://en.m.wikipedia.org/wiki/Main_Page")
	w.Run()
}

func pemToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func encrypt(publicKey *rsa.PublicKey, message []byte) (string, error) {
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, message, []byte("todo"))
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
