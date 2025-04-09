package libraryEncryption

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"io"

	libraryErrors "github.com/s-r-engineer/library/errors"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength   = 32
	nonceLength = 12
	iterations  = 100000
)

func deriveKey(passphrase, salt string) []byte {
	return pbkdf2.Key([]byte(passphrase), []byte(salt), iterations, keyLength, sha512.New)
}

func EncryptAES(passphrase, salt string, plaintextBytes []byte) ([]byte, error) {
	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	nonce := make([]byte, nonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}
	ciphertext := aesGCM.Seal(nonce, nonce, plaintextBytes, nil)
	return ciphertext, nil
}

func DecryptAES(passphrase, salt string, encryptedBytes []byte) ([]byte, error) {
	key := deriveKey(passphrase, salt)
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}
	if len(encryptedBytes) < nonceLength {
		return nil, libraryErrors.WrapError("ciphertext too short", nil)
	}
	nonce, ciphertext := encryptedBytes[:nonceLength], encryptedBytes[nonceLength:]

	return aesGCM.Open(nil, nonce, ciphertext, nil)
}
