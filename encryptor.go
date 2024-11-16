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

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryIO "github.com/s-r-engineer/library/io"
	libraryPath "github.com/s-r-engineer/library/path"
	libraryStrings "github.com/s-r-engineer/library/strings"
	"golang.org/x/crypto/pbkdf2"
)

const (
	keyLength   = 32
	nonceLength = 12
	iterations  = 100000
	cycles      = 5000
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

func getTokenPath() (string, string) {
	homedir, _ := libraryPath.GetHomeDir()
	tokenPath := homedir + "/.config/norrvpn"
	tokenFullPath := tokenPath + "/token.json"
	return tokenPath, tokenFullPath
}

func deriveKey(passphrase, salt string) []byte {
	return pbkdf2.Key([]byte(passphrase), []byte(salt), iterations, keyLength, sha512.New)
}

func encryptAES(passphrase, plaintext, salt string) (string, error) {
	key := deriveKey(passphrase, salt)
	plaintextBytes := []byte(plaintext)

	block, err := aes.NewCipher(key)
	if err != nil {
		libraryErrors.Panicer(fmt.Errorf("could not create cipher: %w", err))
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		libraryErrors.Panicer(fmt.Errorf("could not create GCM mode: %w", err))
	}

	nonce := make([]byte, nonceLength)
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		libraryErrors.Panicer(fmt.Errorf("could not generate nonce: %w", err))
	}

	ciphertext := aesGCM.Seal(nonce, nonce, plaintextBytes, nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

func decryptAES(passphrase, encryptedBase64, salt string) (string, error) {
	key := deriveKey(passphrase, salt)

	ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		libraryErrors.Panicer(fmt.Errorf("could not decode base64: %w", err))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		libraryErrors.Panicer(fmt.Errorf("could not create cipher: %w", err))
	}

	aesGCM, err := cipher.NewGCM(block)
	if err != nil {
		libraryErrors.Panicer(fmt.Errorf("could not create GCM mode: %w", err))
	}

	if len(ciphertext) < nonceLength {
		libraryErrors.Panicer("ciphertext too short")
	}
	nonce, ciphertext := ciphertext[:nonceLength], ciphertext[nonceLength:]

	plaintext, err := aesGCM.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		libraryErrors.Panicer(fmt.Errorf("decryption failed: %w", err))
	}

	return string(plaintext), nil
}

func getToken(pin string) string {
	_, tokenFullPath := getTokenPath()
	data, err := os.ReadFile(tokenFullPath)
	libraryErrors.Panicer(err)
	var token Token
	err = json.Unmarshal(data, &token)
	libraryErrors.Panicer(err)
	str, err := decryptAES(pin, token.Token, token.Salt)
	libraryErrors.Panicer(err)
	return str
}

func setToken(pin, token string) {
	tokenPath, tokenFullPath := getTokenPath()
	libraryErrors.Panicer(libraryIO.CreataeDirs(tokenPath))
	salt := libraryStrings.RandString(666)
	file, err := os.OpenFile(tokenFullPath, os.O_CREATE|os.O_WRONLY, 0600)
	libraryErrors.Panicer(err)
	encryptedToken, err := encryptAES(pin, token, salt)
	libraryErrors.Panicer(err)
	tokenObject := Token{Salt: salt, Token: encryptedToken}
	data, err := json.MarshalIndent(tokenObject, "", "  ")
	libraryErrors.Panicer(err)
	_, err = file.Write(data)
	libraryErrors.Panicer(err)
}

type Token struct {
	Salt, Token string
}
