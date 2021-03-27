package rsahelpers

import (
	"bytes"
	"crypto/x509"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEncryptDecrypt(t *testing.T) {
	privateKey, err := GenerateKey()
	assert.NoError(t, err)

	t.Run("Encrypt and Decrypt", func(t *testing.T) {
		message := "Hello, my name is Inigo Montoya"
		encrypted, err := Encrypt(&privateKey.PublicKey, []byte(message))
		assert.NoError(t, err)
		decrypted, err := Decrypt(privateKey, encrypted)
		assert.NoError(t, err)
		assert.Equal(t, message, string(decrypted))
	})

	t.Run("Marshaling public keys", func(t *testing.T) {
		pubKeyBytes, err := GetSSHPublicKey(privateKey)
		assert.NoError(t, err)
		pubKey, err := SSHPubKeyToRSAPubKey(pubKeyBytes)
		assert.NoError(t, err)
		assert.Equal(t, pubKey, &privateKey.PublicKey)
	})

	t.Run("Reading/writing private key with correct password", func(t *testing.T) {
		password := []byte("password")
		var buf bytes.Buffer
		err := writePrivateKey(password, privateKey, &buf)
		assert.NoError(t, err)
		writtenBytes := buf.Bytes()
		parsedPrivateKey, err := readPrivateKey(password, writtenBytes)
		assert.NoError(t, err)
		assert.Equal(t, privateKey, parsedPrivateKey)
	})

	t.Run("Reading/writing private key with incorrect password", func(t *testing.T) {
		var buf bytes.Buffer
		err := writePrivateKey([]byte("password"), privateKey, &buf)
		assert.NoError(t, err)
		writtenBytes := buf.Bytes()
		_, err = readPrivateKey([]byte("Incorrect password"), writtenBytes)
		assert.Error(t, err)
		assert.ErrorIs(t, err, x509.IncorrectPasswordError)
	})
}
