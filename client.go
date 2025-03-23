package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"strings"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryIO "github.com/s-r-engineer/library/io"
	libraryLogging "github.com/s-r-engineer/library/logging"
)

type request struct {
	Verb    VerbType `json:"verb"`
	Country string   `json:"country,omitempty"`
	Token   string   `json:"token"`
}

func clientMode() error {
	requestData := request{}

	function := flag.Arg(0)
	switch function {
	case "rotate", "r", "reconnect", "renew":
		country, err := getCountry()

		if err != nil {
			return err
		}
		requestData.Verb = rotate
		requestData.Country = country
		token, err := parseToken()

		if err != nil {
			return err
		}
		requestData.Token = token
	case "up", "u", "c", "connect":
		if flag.NArg() == 2 {
			countryCode := strings.ToLower(flag.Arg(1))
			if countryCode == "uk" {
				requestData.Country = "gb"
			} else {
				requestData.Country = flag.Arg(1)
			}
		}
		token, err := parseToken()
		if err != nil {
			return err
		}
		requestData.Token = token
		requestData.Verb = up
	case "down", "d", "disconnect":
		requestData.Verb = down
	case "listCountries", "countriesList", "list", "l":
		requestData.Verb = list
	case "init", "i":
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
		return libraryErrors.WrapError("failed with message", errors.New(responseData.Error))
	}
	err = setCountry(responseData.Country)
	if err != nil {
		libraryLogging.Warn(fmt.Sprintf("could not set the country -> %v", err))
	}
	if requestData.Verb == list {
		formatTable(responseData.Countries)
	}
	return nil
}
