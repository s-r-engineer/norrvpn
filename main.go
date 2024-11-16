// Command wgctrl is a testing utility for interacting with WireGuard via package
// wgctrl.
package main

import (
	"flag"
	"fmt"
	"os"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryIO "github.com/s-r-engineer/library/io"

	"github.com/olekukonko/tablewriter"
)

const defaultNordvpnAddress = "10.5.0.2/32"

const interfaceName = "norrvpn01"

func main() {
	flag.Parse()
	function := flag.Arg(0)
	var host, key string
	switch function {
	case "up":
		if flag.NArg() == 2 {
			host, key = FetchServerData(getCountryCode(flag.Arg(1)))
		} else {
			host, key = FetchServerData(-1)
		}
		pin1, _ := libraryIO.ReadSecretInput("Enter PIN")
		privateKey := fetchOwnPrivateKey(pin1)
		execWGup(interfaceName, privateKey, key, host, defaultNordvpnAddress)
	case "down":
		execWGdown(interfaceName, defaultNordvpnAddress)
	case "init":
		pin1, _ := libraryIO.ReadSecretInput("Enter PIN")
		pin2, _ := libraryIO.ReadSecretInput("Enter PIN again")
		if pin1 != pin2 {
			libraryErrors.Panicer(fmt.Errorf("pins not match"))
		}
		token, _ := libraryIO.ReadSecretInput("Enter TOKEN")
		setToken(pin1, token)
	case "showToken":
		pin1, _ := libraryIO.ReadSecretInput("Enter PIN")
		privateKey := getToken(pin1)
		fmt.Println(privateKey)
	case "listCountries":
		table := tablewriter.NewWriter(os.Stdout)
		for _, country := range getCountryList() {
			table.Append([]string{country.Name, country.Code})
		}
		headers := []string{"Country", "Code"}
		table.SetHeader(headers)
		table.Render()
	}
}
