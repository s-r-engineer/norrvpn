package main

import (
	"encoding/json"
	"io"
	"net/http"
	"os"
	"strings"

	"github.com/olekukonko/tablewriter"
)

func getCountryList() (Countries, error) {
	resp, err := http.Get("https://api.nordvpn.com/v1/countries")
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

func getCountryCode(code string) int {
	countries, _ := getCountryList()
	for _, country := range countries {
		if strings.EqualFold(country.Code, code) {
			return country.ID
		}
	}
	return -1
}

func formatTable(c Countries) {
	table := tablewriter.NewWriter(os.Stdout)
	for _, country := range c {
		table.Append([]string{country.Name, country.Code})
	}
	headers := []string{"Country", "Code"}
	table.SetHeader(headers)
	table.Render()
}

type Countries []struct {
	ID   int    `json:"id"`
	Code string `json:"code"`
	Name string `json:"name"`
}
