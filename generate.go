package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"

	"golang.org/x/term"
)

func main() {
	fmt.Print("Enter Password: ")
	password, err := term.ReadPassword(0)
	if err != nil {
		panic(err)
	}
	keyPair, err := generateKeyPair(password)
	if err != nil {
		log.Fatal("Error generating key pair")
	}
	fmt.Println(keyPair)
}

type KeyPair struct {
	privateKeyPem string
	publicKeyPem  string
}

func generateKeyPair(password []byte) (*KeyPair, error) {
	rng := rand.Reader
	privateKey, err := rsa.GenerateKey(rng, 4096)
	if err != nil {
		return nil, err
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
		return nil, err
	}

	privateKeyPem := pem.EncodeToMemory(encryptedBlock)
	if err != nil {
		return nil, err
	}
	publicKeyPem := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(&privateKey.PublicKey),
		},
	)
	return &KeyPair{
		privateKeyPem: string(privateKeyPem), publicKeyPem: string(publicKeyPem),
	}, nil
}
