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
	encodedCiphertext, err := utils.GetMessageFromStdin()
	utils.FatallyLogOnError("Could not read ciphertext", err)
	ciphertext, err := base64.StdEncoding.DecodeString(string(encodedCiphertext))
	utils.FatallyLogOnError("Could not decode base64 ciphertext", err)
	password, err := utils.ReadPassword("Enter password: ")
	utils.FatallyLogOnError("Could not read password", err)
	privateKey, err := pemToPrivateKey([]byte(privateKeyPem), password)
	utils.FatallyLogOnError("Could not parse private key", err)
	message, err := decrypt(privateKey, ciphertext)
	utils.FatallyLogOnError("Failed to decrypt", err)
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
