package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"math/big"
	"net"
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

var (
	dhSalt  = "859175917822340865368638064107081497694254587808895442203370028558935503499994850711769500560687917387996768239918428042499870100570565248853042582404776690722910909663737163697823565083322024596921969408679679684623926697230750621513269868859237316042323122167565768320567347076220049134816200670544019755296817827651530589714950784022357602389003599739840349286412810816129335969325183307433472606170742633830125025545888778050968674491297940842926513238898325060867925523477823975245029044529077641837348938000631321885554368210797390201234225046308653922848700418655869152082329314883861616235413074340679363173284798206828660211800038059756719146718785468126278681878932204273090518684577559059731331540279848972781697305693590914167156502951853761408842456883658002977311365465365591577572426412664317017664995014998647487062188906539444264256306881959923233851315633272005391738204253803288850919363747368143967995176822283961878000735019030487924401109281400728036696331659531624912345718255359133309514044716579869332145390676404168167476282417960501878912045843293942818852573638482519709562186976228846659642577267023708026690051154431474522837705668821587798886212396513437493247990402608493967711814978249121775526845053"
	pSource string
	p, _    = new(big.Int).SetString(pSource, 16)
	g       = big.NewInt(886)
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

func encryptAES(passphrase string, plaintextBytes []byte, salt string) ([]byte, error) {
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

func decryptAES(passphrase string, encryptedBytes []byte, salt string) ([]byte, error) {
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

func getToken(pin string) (string, error) {
	_, tokenFullPath := getTokenPath()
	data, err := os.ReadFile(tokenFullPath)
	if err != nil {
		return "", err
	}
	var token Token
	err = json.Unmarshal(data, &token)
	if err != nil {
		return "", err
	}
	tokenEncryptedBytes, err := base64.StdEncoding.DecodeString(token.Token)
	if err != nil {
		return "", err
	}
	b, err := decryptAES(pin, tokenEncryptedBytes, token.Salt)
	return string(b), err
}

func setToken(pin, token string) error {
	tokenPath, tokenFullPath := getTokenPath()
	err := libraryIO.CreateDirs(tokenPath)
	if err != nil {
		return err
	}
	salt := libraryStrings.RandString(666)
	file, err := os.OpenFile(tokenFullPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	encryptedToken, err := encryptAES(pin, []byte(token), salt)
	if err != nil {
		return err
	}
	tokenObject := Token{Salt: salt, Token: base64.StdEncoding.EncodeToString(encryptedToken)}
	data, err := json.MarshalIndent(tokenObject, "", "  ")
	if err != nil {
		return err
	}
	_, err = file.Write(data)
	return err
}

func getDHSecret(conn net.Conn) (string, error) {
	priv, err := rand.Int(rand.Reader, p)
	if err != nil {
		return "", err
	}
	pub := new(big.Int).Exp(g, priv, p)
	otherSidePub := make([]byte, p.BitLen()/8+1)
	_, err = conn.Write(pub.Bytes())
	if err != nil {
		return "", err
	}
	n, err := conn.Read(otherSidePub)
	if err != nil {
		return "", err
	}
	otherSide := new(big.Int).SetBytes(otherSidePub[:n])
	sharedSecret := new(big.Int).Exp(otherSide, priv, p)
	symmetricKey := sha256.Sum256(sharedSecret.Bytes())
	return fmt.Sprintf("%x", symmetricKey), nil
}

type Token struct {
	Salt, Token string
}
