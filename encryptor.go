package main

import (
	"encoding/base64"
	"encoding/json"
	"math/big"
	"os"

	libraryEncryption "github.com/s-r-engineer/library/encryption"
	libraryIO "github.com/s-r-engineer/library/io"
	libraryPath "github.com/s-r-engineer/library/path"
	libraryStrings "github.com/s-r-engineer/library/strings"
)

const (
	keyLength   = 32
	nonceLength = 12
	iterations  = 100000
	dhSalt      = "859175917822340865368638064107081497694254587808895442203370028558935503499994850711769500560687917387996768239918428042499870100570565248853042582404776690722910909663737163697823565083322024596921969408679679684623926697230750621513269868859237316042323122167565768320567347076220049134816200670544019755296817827651530589714950784022357602389003599739840349286412810816129335969325183307433472606170742633830125025545888778050968674491297940842926513238898325060867925523477823975245029044529077641837348938000631321885554368210797390201234225046308653922848700418655869152082329314883861616235413074340679363173284798206828660211800038059756719146718785468126278681878932204273090518684577559059731331540279848972781697305693590914167156502951853761408842456883658002977311365465365591577572426412664317017664995014998647487062188906539444264256306881959923233851315633272005391738204253803288850919363747368143967995176822283961878000735019030487924401109281400728036696331659531624912345718255359133309514044716579869332145390676404168167476282417960501878912045843293942818852573638482519709562186976228846659642577267023708026690051154431474522837705668821587798886212396513437493247990402608493967711814978249121775526845053"
	cycles      = 5000
	letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
)

var (
	pSource string
	p, _    = new(big.Int).SetString(pSource, 16)
	g       = big.NewInt(886)
)

func getConfigPath() (string, string, string) {
	homedir, _ := libraryPath.GetHomeDir()
	tokenPath := homedir + "/.config/norrvpn"
	tokenFullPath := tokenPath + "/token.json"
	countryFullPAth := tokenPath + "/country"
	return tokenPath, tokenFullPath, countryFullPAth
}

func parseToken() (string, error) {
	pin1, err := libraryIO.ReadSecretInput("Enter PIN")
	if err != nil {
		return "", err
	}
	return getToken(pin1)
}

func getToken(pin string) (string, error) {
	_, tokenFullPath, _ := getConfigPath()
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
	b, err := libraryEncryption.DecryptAES(pin, token.Salt, tokenEncryptedBytes)
	return string(b), err
}

func setToken(pin, token string) error {
	tokenPath, tokenFullPath, _ := getConfigPath()
	err := libraryIO.CreateDirs(tokenPath)
	if err != nil {
		return err
	}
	salt := libraryStrings.RandString(666)
	file, err := os.OpenFile(tokenFullPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	encryptedToken, err := libraryEncryption.EncryptAES(pin, salt, []byte(token))
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

type Token struct {
	Salt, Token string
}
