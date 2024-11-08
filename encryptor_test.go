package main

import (
	"testing"
)

func TestEncryptor(t *testing.T) {
	var encrypted, decrypted string
	salt := randString(16)
	data := randString(666)
	passphrase := randString(32)
	encrypted, _ = encryptAES(passphrase, data, salt)
	decrypted, _ = decryptAES(passphrase, encrypted, salt)
	if data != decrypted {
		t.Fail()
	}
}
