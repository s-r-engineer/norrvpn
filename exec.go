package main

import (
	"fmt"
	"regexp"
	"strings"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryExec "github.com/s-r-engineer/library/exec"
	libraryStrings "github.com/s-r-engineer/library/strings"
)

const defaultWGPort = "51820"

func getEndpointIP(lines []string) string {
	re := regexp.MustCompile(`^219.*from all to ([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}) lookup main$`)
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); matches != nil {
			return libraryStrings.Trim(matches[1])
		}
	}
	libraryErrors.Panicer("no IP found:\n" + strings.Join(lines, "\n"))
	return ""
}

func execWGdown(interfaceName, interfaceIP string) error {
	_, _, err := libraryExec.Run("ip", "route", "delete", "default", "dev", interfaceName, "table", "212450")
	if err != nil {
		return libraryErrors.WrapError("delete default route", err)
	}
	out, _, err := libraryExec.Run("ip", "rule", "show")
	if err != nil {
		return libraryErrors.WrapError("ip rule show", err)
	}
	_, _, err = libraryExec.Run("ip", "rule", "delete", "to", getEndpointIP(strings.Split(out, "\n")), "table", "main", "priority", "219")
	if err != nil {
		return libraryErrors.WrapError("delete rule for server", err)
	}
	_, _, err = libraryExec.Run("ip", "rule", "delete", "lookup", "212450", "priority", "220")
	if err != nil {
		return libraryErrors.WrapError("delete lookup rule", err)
	}
	_, _, err = libraryExec.Run("ip", "link", "set", "down", "dev", interfaceName)
	if err != nil {
		return libraryErrors.WrapError("link down", err)
	}
	_, _, err = libraryExec.Run("ip", "address", "del", interfaceIP, "dev", interfaceName)
	if err != nil {
		return libraryErrors.WrapError("delete address from the interface", err)
	}
	_, _, err = libraryExec.Run("ip", "link", "delete", "dev", interfaceName)
	if err != nil {
		return libraryErrors.WrapError("delete interface", err)
	}
	return nil
}

func execWGup(interfaceName, privateKey, publicKey, endpointIP, interfaceIP string) (err error) {
	_, code, _ := libraryExec.Run("ip", "link", "show", interfaceName)
	if code == 1 {
		_, _, err = libraryExec.Run("ip", "link", "add", "dev", interfaceName, "type", "wireguard")
		if err != nil {
			return libraryErrors.WrapError("add interface", err)
		}
	}
	_, _, err = libraryExec.RunWithStdin(privateKey, "wg", "set", interfaceName, "private-key", "/dev/stdin")
	if err != nil {
		return libraryErrors.WrapError("set private key", err)
	}
	_, _, err = libraryExec.Run("wg", "set", interfaceName, "peer", publicKey, "endpoint", endpointIP+":"+defaultWGPort, "allowed-ips", "0.0.0.0/0")
	if err != nil {
		return libraryErrors.WrapError("set peer", err)
	}
	_, _, err = libraryExec.Run("ip", "address", "add", interfaceIP, "dev", interfaceName)
	if err != nil {
		return libraryErrors.WrapError("set ip address", err)
	}
	_, _, err = libraryExec.Run("ip", "link", "set", "up", "dev", interfaceName)
	if err != nil {
		return libraryErrors.WrapError("up interface", err)
	}
	_, _, err = libraryExec.Run("ip", "route", "add", "default", "dev", interfaceName, "table", "212450")
	if err != nil {
		return libraryErrors.WrapError("add default route", err)
	}
	_, _, err = libraryExec.Run("ip", "rule", "add", "to", endpointIP, "table", "main", "priority", "219")
	if err != nil {
		return libraryErrors.WrapError("add server route", err)
	}
	_, _, err = libraryExec.Run("ip", "rule", "add", "lookup", "212450", "priority", "220")
	if err != nil {
		return libraryErrors.WrapError("set lookup rule", err)
	}
	return nil
}

func setAddress(interfaceName, interfaceIP string) error {
	_, _, err := libraryExec.Run("ip", "address", "add", interfaceIP, "dev", interfaceName)
	return libraryErrors.WrapError(fmt.Sprintf("setAddress -> %s, %s", interfaceName, interfaceIP), err)
}
