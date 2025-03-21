package libraryNordvpn

import (
	"encoding/json"
	libraryHttp "github.com/s-r-engineer/library/http"
	"io"
	"strings"
)

const DefaultCountriesListURL = "https://api.nordvpn.com/v1/countries"

func GetCountryList() (Countries, error) {
	resp, err := libraryHttp.GetUrl(DefaultCountriesListURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	c := Countries{}
	err = json.Unmarshal(data, &c)
	return c, err
}

func GetCountryCode(code string) int {
	countries, _ := GetCountryList()
	for _, country := range countries {
		if strings.EqualFold(country.Code, code) {
			return country.ID
		}
	}
	return -1
}

type Countries []struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
