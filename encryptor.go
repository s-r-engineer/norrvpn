package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength   = 32
	nonceLength = 12
	iterations  = 100000
	cycles      = 5000
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func deriveKey(passphrase, salt string) []byte {
	return pbkdf2.Key([]byte(passphrase), []byte(salt), iterations, keyLength, sha512.New)
}

func encryptAES(passphrase, plaintext, salt string) (string, error) {
	key := deriveKey(passphrase, salt)
	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		panicer(fmt.Errorf("could not create cipher: %w", err))
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panicer(fmt.Errorf("could not create GCM mode: %w", err))
	}

	nonce := make([]byte, nonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		panicer(fmt.Errorf("could not generate nonce: %w", err))
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintextBytes, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(passphrase, encryptedBase64, salt string) (string, error) {
	key := deriveKey(passphrase, salt)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		panicer(fmt.Errorf("could not decode base64: %w", err))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		panicer(fmt.Errorf("could not create cipher: %w", err))
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		panicer(fmt.Errorf("could not create GCM mode: %w", err))
	}

	if len(ciphertext) < nonceLength {
		panicer("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceLength], ciphertext[nonceLength:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		panicer(fmt.Errorf("decryption failed: %w", err))
	}

	return string(plaintext), nil
}

func getToken(pin string) string {
	data, err := os.ReadFile(tokenFullPath)
	panicer(err)
	var token Token
	err = json.Unmarshal(data, &token)
	panicer(err)
	str, err := decryptAES(pin, token.Token, token.Salt)
	panicer(err)
	return str
}

func setToken(pin, token string) {
	panicer(os.MkdirAll(tokenPath, 0700))
	salt := randString(666)
	file, err := os.OpenFile(tokenFullPath, os.O_CREATE|os.O_WRONLY, 0600)
	panicer(err)
	encryptedToken, err := encryptAES(pin, token, salt)
	panicer(err)
	tokenObject := Token{Salt: salt, Token: encryptedToken}
	data, err := json.MarshalIndent(tokenObject, "", "  ")
	panicer(err)
	_, err = file.Write(data)
	panicer(err)
}

type Token struct {
	Salt, Token string
}
