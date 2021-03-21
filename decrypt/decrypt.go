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
)

func main() {

	encodedCiphertext := string(utils.GetMessageFromStdin())
	ciphertext, err := base64.StdEncoding.DecodeString(encodedCiphertext)
	if err != nil {
		panic(err) // todo
	}

	password, err := utils.ReadPassword("Enter password: ")
	if err != nil {
		panic(err) // todo
	}
	privateKey, err := pemToPrivateKey([]byte(privateKeyPem), password)
	if err != nil {
		panic(err) // todo
	}
	message, err := decrypt(privateKey, ciphertext)
	if err != nil {
		panic(err) // todo
	}
	fmt.Println(string(message))
}

func pemToPrivateKey(priv []byte, password []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(priv)
	isEncrypted := x509.IsEncryptedPEMBlock(block)
	pemBytes := block.Bytes
	if isEncrypted {
		var err error
		pemBytes, err = x509.DecryptPEMBlock(block, password)
		if err != nil {
			return nil, err
		}
	}
	return x509.ParsePKCS1PrivateKey(pemBytes)
}

func decrypt(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	rng := rand.Reader
	return rsa.DecryptOAEP(sha256.New(), rng, privateKey, ciphertext, []byte("todo"))
}
