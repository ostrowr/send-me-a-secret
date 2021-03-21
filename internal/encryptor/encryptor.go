package encryptor

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
)

func Encrypt(message []byte) (string, error) {
	publicKey, err := pemToPublicKey([]byte(publicKeyPem))
	if err != nil {
		return "", err
	}
	ciphertext, err := encrypt(publicKey, message)
	if err != nil {
		return "", err
	}
	return ciphertext, nil
}

func pemToPublicKey(pub []byte) (*rsa.PublicKey, error) {
	block, _ := pem.Decode(pub)
	return x509.ParsePKCS1PublicKey(block.Bytes)
}

func encrypt(publicKey *rsa.PublicKey, message []byte) (string, error) {
	rng := rand.Reader

	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, message, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(ciphertext), nil
}
