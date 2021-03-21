package decryptor

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/pem"
)

func Decrypt(ciphertext []byte, password []byte) ([]byte, error) {
	privateKey, err := pemToPrivateKey([]byte(privateKeyPem), password)
	if err != nil {
		return nil, err
	}
	return decrypt(privateKey, ciphertext)
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
