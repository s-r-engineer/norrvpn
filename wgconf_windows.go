package main

import (
	"fmt"
	"os"
)

const defaultWGPort = "51820"

const confTemplate = `[Interface]
PrivateKey = %s
ListenPort = %d
Address = 10.5.0.2/32
DNS = 103.86.96.100, 103.86.99.100

[Peer]
PublicKey = %s
AllowedIPs = 0.0.0.0/0, ::/0
Endpoint = %s:%d
`

func createConfigFile(countryCode, peerEndpoint, peerPublicKey, ownPrivateKey string) (tempFilePath string, f func(), err error) {
	tempFile, err := os.CreateTemp("", countryCode+".conf")
	if err != nil {
		return
	}
	defer tempFile.Close()
	_, err = tempFile.Write([]byte(fmt.Sprintf(confTemplate, ownPrivateKey, defaultWGPort, peerPublicKey, peerEndpoint, defaultWGPort)))
	if err != nil {
		return
	}
	return tempFile.Name(), func() { os.Remove(tempFile.Name()) }, nil
}

// func execWGdown(interfaceName, interfaceIP string) {
// }

// func execWGup(interfaceName, privateKey, publicKey, endpointIP, interfaceIP string) error {

// }

// func setAddress(interfaceName, interfaceIP string) error {
// 	_, _, err := libraryExec.Run("ip", "address", "add", interfaceIP, "dev", interfaceName)
// 	return libraryErrors.WrapError(fmt.Sprintf("setAddress -> %s, %s", interfaceName, interfaceIP), err)
// }
