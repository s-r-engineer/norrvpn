package libraryNordvpn

import (
	"encoding/json"
	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryHttp "github.com/s-r-engineer/library/http"
	"io"
	"strings"
)

const DefaultCountriesListURL = "https://api.nordvpn.com/v1/countries"

var countryList map[string]int

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

func GetCountryCode(code string) (int, error) {
	wrapper := libraryErrors.PartWrapError("GetCountryCode")
	err := populateCountryList()
	if err != nil {
		return -2, wrapper(err)
	}
	if id, ok := countryList[code]; ok {
		return id, nil
	}
	return -1, nil
}

func populateCountryList() error {
	if countryList != nil {
		return nil
	}
	countries, err := GetCountryList()
	if err != nil {
		return err
	}
	countryList = make(map[string]int)
	for _, c := range countries {
		countryList[strings.ToLower(c.Code)] = c.ID
	}
	countryList["uk"] = countryList["gb"]
	return nil
}

type Countries []struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
