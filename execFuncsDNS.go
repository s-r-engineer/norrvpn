package main

import (
	"errors"
	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryExec "github.com/s-r-engineer/library/exec"
)

var resolver byte

const defaultNordVPNDNS = "103.86.96.100 103.86.99.100"

func init() {
	libraryErrors.Panicer(checkDNSResolver())
}

func checkDNSResolver() error {
	_, code, err := libraryExec.Run("resolvectl", "status")
	if code == 0 {
		resolver = 'r'
		return nil
	} else if code == 1 {
		resolver = 'n'
		return nil
	}
	return err
}

func setDNS(interfaceName, dnsServers string) error {
	switch resolver {
	case 'r':
		return setResolvedDNS(interfaceName, dnsServers)
	case 'n':
		return setNMDNS(interfaceName, dnsServers)
	default:
		return errors.New("resolver undefined")
	}
}

func setNMDNS(interfaceName, dnsServers string) (err error) {
	_, _, err = libraryExec.Run("nmcli", "connection", "modify", interfaceName, "ipv4.dns", dnsServers)
	if err != nil {
		return err
	}
	_, _, err = libraryExec.Run("nmcli", "connection", "modify", interfaceName, "ipv4.ignore-auto-dns", "yes")
	return err
}
func setResolvedDNS(interfaceName, dnsServers string) error {
	_, _, err := libraryExec.Run("resolvectl", "dns", interfaceName, dnsServers)
	if err != nil {
		return err
	}
	_, _, err = libraryExec.Run("resolvectl", "domain", interfaceName, "~.")
	return err
}
