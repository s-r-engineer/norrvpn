package main

import (
	"testing"

	libraryStrings "github.com/s-r-engineer/library/strings"
)

func TestEncryptor(t *testing.T) {
	var encrypted, decrypted []byte
	salt := libraryStrings.RandString(16)
	data := libraryStrings.RandString(666)
	passphrase := libraryStrings.RandString(32)
	encrypted, _ = encryptAES(passphrase, []byte(data), salt)
	decrypted, _ = decryptAES(passphrase, encrypted, salt)
	if data != string(decrypted) {
		t.Fail()
	}
}
