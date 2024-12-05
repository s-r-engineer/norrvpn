package main

import (
	"os"

	libraryIO "github.com/s-r-engineer/library/io"

	"github.com/olekukonko/tablewriter"
	libraryNordvpn "github.com/s-r-engineer/library/nordvpn"
)

func formatTable(c libraryNordvpn.Countries) {
	table := tablewriter.NewWriter(os.Stdout)
	for _, country := range c {
		table.Append([]string{country.Name, country.Code})
	}
	headers := []string{"Country", "Code"}
	table.SetHeader(headers)
	table.Render()
}

func getCountry() (string, error) {
	_, _, countryPath := getConfigPath()
	data, err := os.ReadFile(countryPath)
	if err != nil {
		return "", err
	}
	return string(data), err
}

func setCountry(country string) error {
	configPath, _, countryPath := getConfigPath()
	err := libraryIO.CreateDirs(configPath)
	if err != nil {
		return err
	}
	file, err := os.OpenFile(countryPath, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return err
	}
	_, err = file.Write([]byte(country))
	return err
}
