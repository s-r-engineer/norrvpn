package main

import (
	"encoding/json"
	"fmt"
	"net"
	"os"

	libraryEncryption "github.com/s-r-engineer/library/encryption"
	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryLogging "github.com/s-r-engineer/library/logging"
	libraryNordvpn "github.com/s-r-engineer/library/nordvpn"
)

const (
	defaultNordvpnAddress = "10.5.0.2/32"
	defaultReadBlockSize  = 1024
	interfaceName         = "norrvpn01"
	defaultWGPort         = "51820"
	defaultRouteTable     = "212450"
)

type VerbType int

var (
	up     VerbType = 1
	down   VerbType = 2
	list   VerbType = 3
	rotate VerbType = 4
	//status VerbType = 4

	//currentCountry string
)

type response struct {
	Result bool   `json:"result"`
	Error  string `json:"error,omitempty"`

	Country string `json:"country,omitempty"`

	Countries libraryNordvpn.Countries `json:"countries,omitempty"`
}

func serverMode() error {
	listener, err := getSocketListener()
	if err != nil {
		return err
	}

	defer func() {
		libraryErrors.Errorer(listener.Close())
	}()
	defer func() {
		libraryErrors.Errorer(os.Remove(socketPath))
	}()

	fmt.Println("Server listening on", socketPath)

	for {
		libraryLogging.Info("accepting connections")
		conn, err := listener.Accept()
		if err != nil {
			return err
		}

		symmetricKey, err := libraryEncryption.GetDHSecretFromConnection(conn, p, g)
		if err != nil {
			return err
		}

		go serve(conn, symmetricKey)
	}
}

func buildResponse(err error) response {
	if err != nil {
		return response{Result: false, Error: err.Error()}
	}
	return response{Result: true}
}

func serve(conn net.Conn, secret string) {
	defer func() {
		libraryErrors.Errorer(conn.Close())
	}()
	fullData, err := parseConnection(conn)
	if err != nil {
		libraryLogging.Error(err.Error())
		return
	}
	decryptedData, err := libraryEncryption.DecryptAES(secret, dhSalt, fullData)
	if err != nil {
		libraryLogging.Error(err.Error())
		return
	}
	requestStruct := request{}
	err = json.Unmarshal(decryptedData, &requestStruct)
	if err != nil {
		libraryLogging.Error(err.Error())
		return
	}

	var ip, key string

	var result response

	switch requestStruct.Verb {
	case up, rotate:
		var countryCode string
		if requestStruct.Country != "" {
			code, err1 := libraryNordvpn.GetCountryCode(requestStruct.Country)
			if err1 != nil {
				libraryLogging.Error(err1.Error())
				return
			}
			_, ip, key, countryCode, err = libraryNordvpn.FetchServerData(code)
		} else {
			_, ip, key, countryCode, err = libraryNordvpn.FetchServerData(-1)
		}
		if err != nil {
			libraryLogging.Error(err.Error())
			return
		}
		privateKey, err := libraryNordvpn.FetchOwnPrivateKey(requestStruct.Token)
		if err != nil {
			libraryLogging.Error(err.Error())
			return
		}
		err = execWGup(interfaceName, privateKey, key, ip, defaultNordvpnAddress, defaultWGPort, defaultRouteTable)
		result = buildResponse(err)
		result.Country = countryCode
	case down:
		err = execWGdown(interfaceName, defaultNordvpnAddress, defaultRouteTable)
		result = buildResponse(err)
	case list:
		countries, err := libraryNordvpn.GetCountryList()
		result = buildResponse(err)
		result.Countries = countries
	}
	data, err := json.Marshal(result)
	libraryErrors.Errorer(err)
	encryptedData, err := libraryEncryption.EncryptAES(secret, dhSalt, data)
	libraryErrors.Errorer(err)
	_, err = conn.Write(encryptedData)
	libraryErrors.Errorer(err)
}
