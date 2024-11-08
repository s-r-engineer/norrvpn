// Command wgctrl is a testing utility for interacting with WireGuard via package
// wgctrl.
package main

import (
	"flag"
	"fmt"
	"os"

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
		pin1 := readSecretInput("Enter PIN")
		privateKey := fetchOwnPrivateKey(pin1)
		execWGup(interfaceName, privateKey, key, host, defaultNordvpnAddress)
	case "down":
		execWGdown(interfaceName, defaultNordvpnAddress)
	case "init":
		pin1 := readSecretInput("Enter PIN")
		pin2 := readSecretInput("Enter PIN again")
		if pin1 != pin2 {
			panicer(fmt.Errorf("pins not match"))
		}
		token := readSecretInput("Enter TOKEN")
		setToken(pin1, token)
	case "showToken":
		pin1 := readSecretInput("Enter PIN")
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
