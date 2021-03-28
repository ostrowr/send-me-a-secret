package rsahelpers

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"io"
	"os"
	"os/user"
	"path/filepath"

	"github.com/ostrowr/send-me-a-secret/internal/utils"
	"golang.org/x/crypto/ssh"
)

// Github strips out comments and doesn't allow options on public keys, and I haven't figured out an elegant way
// to mark a public key as belonging to send-me-a-secret without doing something weird like
// creating a gist or updating a bio to point to the right key. In the meantime, we'll use this
// nontraditional key length and assume the user has no other keys of length 4567.
// Encryption will fail if there is not exactly one key of length WeirdKeyLength in the github user's account.
const WeirdKeyLength = 4568

var KeyFilename = ".send-me-a-secret"

func PathToKeyFile() string {
	usr, err := user.Current()
	if err != nil {
		panic(err)
	}
	return filepath.Join(usr.HomeDir, KeyFilename)
}

// GenerateKey generates a new RSA private key with key length WEIRD_KEY_LENGTH
func GenerateKey() (*rsa.PrivateKey, error) {
	rng := rand.Reader
	return rsa.GenerateKey(rng, WeirdKeyLength)
}

// WritePrivateKeyToFile writes an rsa private key to ~/.send-me-a-secret
// This path is not configurable; don't want a user to be able to forget
// where they saved their key.
func WritePrivateKeyToFile(password []byte, privateKey *rsa.PrivateKey) error {
	keyfile, err := os.Create(PathToKeyFile())
	if err != nil {
		return err
	}
	err = os.Chmod(PathToKeyFile(), 0600)
	if err != nil {
		return err
	}
	defer utils.MustClose(keyfile)
	return writePrivateKey(password, privateKey, keyfile)
}

func writePrivateKey(password []byte, privateKey *rsa.PrivateKey, dest io.Writer) error {
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

	err = pem.Encode(dest, encryptedBlock)
	if err != nil {
		return err
	}

	return nil
}

// ReadPrivateKeyFromFile reads an rsa private key from ~/.send-me-a-secret
// This path is not configurable; don't want a user to be able to forget
// where they saved their key.
func ReadPrivateKeyFromFile(password []byte) (*rsa.PrivateKey, error) {
	keyPem, err := os.ReadFile(PathToKeyFile())
	if err != nil {
		return nil, err
	}
	return readPrivateKey(password, keyPem)
}

func readPrivateKey(password, keyPem []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(keyPem)
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

// GetSSHPublicKey generates a public key suitable for openssh (and thus GitHub) from a private key.
func GetSSHPublicKey(privateKey *rsa.PrivateKey) ([]byte, error) {
	publicKey, err := ssh.NewPublicKey(&privateKey.PublicKey)
	if err != nil {
		return nil, err
	}
	marshaled := ssh.MarshalAuthorizedKey(publicKey)
	return marshaled, nil
}

var ErrInvalidPublicKey = errors.New("invalid public key")

func SSHPubKeyToRSAPubKey(sshPubKey []byte) (*rsa.PublicKey, error) {
	pubKey, _, _, _, err := ssh.ParseAuthorizedKey(sshPubKey)
	if err != nil {
		return nil, err
	}
	keyAsCryptoKey, ok := pubKey.(ssh.CryptoPublicKey)
	if !ok {
		return nil, ErrInvalidPublicKey
	}
	rsaPubKey, ok := keyAsCryptoKey.CryptoPublicKey().(*rsa.PublicKey)
	if !ok {
		return nil, ErrInvalidPublicKey
	}
	return rsaPubKey, nil
}

// Encrypt encrypts a message under the given public key, suitable for decrypting via `Decrypt`
func Encrypt(publicKey *rsa.PublicKey, message []byte) (string, error) {
	rng := rand.Reader
	ciphertext, err := rsa.EncryptOAEP(sha256.New(), rng, publicKey, message, nil)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts a message using the given private key which was encrypted by `Encrypt`
func Decrypt(privateKey *rsa.PrivateKey, base64EncodedCiphertext string) ([]byte, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(base64EncodedCiphertext)
	if err != nil {
		return nil, err
	}
	rng := rand.Reader
	return rsa.DecryptOAEP(sha256.New(), rng, privateKey, ciphertext, nil)
}

// IsValidSendMeASecretKey checks if the key fetched from GitHub is the key uploaded by send-me-a-secret.
// Right now, this just checks that length of the key is WEIRD_KEY_LENGTH, hoping that the user doesn't have
// any other keys of that length, but hopefully in the future we'll be able to do something a bit cleverer.
func IsValidSendMeASecretKey(publicKey *rsa.PublicKey) bool {
	return publicKey.Size() == WeirdKeyLength/8
}
