package main

import (
	"testing"
	libraryStrings "github.com/s-r-engineer/library/strings"
)

func TestEncryptor(t *testing.T) {
	var encrypted, decrypted string
	salt := libraryStrings.RandString(16)
	data := libraryStrings.RandString(666)
	passphrase := libraryStrings.RandString(32)
	encrypted, _ = encryptAES(passphrase, data, salt)
	decrypted, _ = decryptAES(passphrase, encrypted, salt)
	if data != decrypted {
		t.Fail()
	}
}
