package main

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"

	libraryErrors "github.com/s-r-engineer/library/errors"
)

func getCountryList() countries {
	resp, err := http.Get("https://api.nordvpn.com/v1/countries")
	libraryErrors.Panicer(err)
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	libraryErrors.Panicer(err)
	c := countries{}
	libraryErrors.Panicer(json.Unmarshal(data, &c))
	return c
}

func getCountryCode(code string) int {
	for _, country := range getCountryList() {
		if strings.EqualFold(country.Code, code) {
			return country.ID
		}
	}
	return -1
}

type countries []struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
