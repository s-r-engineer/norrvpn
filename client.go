package main

import (
	"encoding/json"
	"flag"
	"fmt"
	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryIO "github.com/s-r-engineer/library/io"
)

type request struct {
	Verb    VerbType `json:"verb"`
	Country string   `json:"country,omitempty"`
	Token   string   `json:"token"`
}

func clientMode() error {
	function := flag.Arg(0)
	requestData := request{}
	switch function {
	case "up":
		if flag.NArg() == 2 {
			requestData.Country = flag.Arg(1)
		}
		pin1, err := libraryIO.ReadSecretInput("Enter PIN")
		if err != nil {
			return err
		}
		requestData.Token, err = getToken(pin1)
		if err != nil {
			return err
		}
		requestData.Verb = up
	case "down":
		requestData.Verb = down
	case "countriesList":
	case "list":
		requestData.Verb = list
	case "init":
		pin1, err := libraryIO.ReadSecretInput("Enter PIN")
		if err != nil {
			return err
		}
		pin2, err := libraryIO.ReadSecretInput("Enter PIN again")
		if err != nil {
			return err
		}
		if pin1 != pin2 {
			libraryErrors.Panicer(fmt.Errorf("pins not match"))
		}
		token, err := libraryIO.ReadSecretInput("Enter TOKEN")
		if err != nil {
			return err
		}
		return setToken(pin1, token)
	}
	conn, err := getSocketConnector()
	if err != nil {
		return err
	}
	defer func() {
		libraryErrors.Errorer(conn.Close())
	}()
	symmetricKey, err := getDHSecret(conn)
	if err != nil {
		return err
	}
	bytes, err := json.Marshal(requestData)
	if err != nil {
		return err
	}
	encryptedData, err := encryptAES(symmetricKey, bytes, dhSalt)
	if err != nil {
		return err
	}
	_, err = conn.Write(encryptedData)
	if err != nil {
		return err
	}
	fullDataEncrypted, err := parseConnection(conn)
	if err != nil {
		return err
	}
	fullData, err := decryptAES(symmetricKey, fullDataEncrypted, dhSalt)
	if err != nil {
		return err
	}
	responseData := response{}
	err = json.Unmarshal(fullData, &responseData)
	if err != nil {
		return err
	}
	if !responseData.Result {
		return libraryErrors.WrapError("failed with message", fmt.Errorf(responseData.Error))
	}
	if requestData.Verb == list {
		formatTable(responseData.Countries)
	}
	return nil
}
