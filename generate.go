//+build ignore

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"os"
	"strings"
	"text/template"

	"github.com/ostrowr/send-me-a-secret/utils"
	"golang.org/x/term"
)

func main() {
	fmt.Print("Enter Password: ")
	password, err := term.ReadPassword(0)
	utils.FatallyLogOnError("Could not read password", err)
	fmt.Println("\nGenerating key pair...")
	keyPair, err := generateKeyPair(password)
	utils.FatallyLogOnError("Could not generate key pair", err)
	fmt.Println("New key pair generated.")
	fillTemplateFile("./decrypt/privateKey.go.template", keyPair)
	fillTemplateFile("./encrypt/publicKey.go.template", keyPair)
}

func fillTemplateFile(path string, data interface{}) error {
	const templateSuffix = ".template"
	if !strings.HasSuffix(path, templateSuffix) {
		return errors.New("invalid template name")
	}
	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}
	f, err := os.Create(strings.TrimSuffix(path, templateSuffix))
	if err != nil {
		return err
	}
	defer f.Close()
	err = t.Option("missingkey=error").Execute(f, data)
	if err != nil {
		return err
	}
	return nil
}

type KeyPair struct {
	PrivateKeyPem string
	PublicKeyPem  string
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
		PrivateKeyPem: string(privateKeyPem), PublicKeyPem: string(publicKeyPem),
	}, nil
}
