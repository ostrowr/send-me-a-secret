package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"

	"golang.org/x/term"
)

func main() {
	fmt.Print("Enter Password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	generatePrivateKey(password)
}

func generatePrivateKey(password []byte) error {
	rng := rand.Reader
	privateKey, err := rsa.GenerateKey(rng, 4096)
	if err != nil {
		return err
	}

	block := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(privateKey),
	}

	encryptedBlock, err := x509.EncryptPEMBlock(
		rand.Reader,
		block.Type,
		block.Bytes,
		password,
		x509.PEMCipherAES256,
	)
	if err != nil {
		return err
	}

	privateKeyPem := pem.EncodeToMemory(encryptedBlock)
	if err != nil {
		return err
	}
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		},
	)
	fmt.Println(string(privateKeyPem))
	fmt.Println(string(publicKeyPem))

	return nil
}

// func encrypt(message string) {
// 	secretMessage := []byte("send reinforcements, we're going to advance")
// 	label := []byte("orders")

// 	// crypto/rand.Reader is a good source of entropy for randomizing the
// 	// encryption function.
// 	rng := rand.Reader

// 	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, &test2048Key.PublicKey, secretMessage, label)
// 	if err != nil {
// 		fmt.Fprintf(os.Stderr, "Error from encryption: %s\n", err)
// 		return
// 	}

// 	// Since encryption is a randomized function, ciphertext will be
// 	// different each time.
// 	fmt.Printf("Ciphertext: %x\n", ciphertext)
// }
