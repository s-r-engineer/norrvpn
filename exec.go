package main

import (
	"errors"
	"fmt"
	"regexp"
	"strings"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryExec "github.com/s-r-engineer/library/exec"
	libraryStrings "github.com/s-r-engineer/library/strings"
)

type checkError struct {
}

func (checkError) Error() string {
	return "check failed"
}

//type checkResult int
//
//var (
//	fail checkResult = 2
//	pass checkResult = 2
//	do   checkResult = 1
//	redo checkResult = 0
//)
//
//type checkStatus struct {
//	Status checkResult
//	Error  error
//}
//
//func newCheckStatus(result checkResult, err error) (c checkStatus) {
//	c.Status = result
//	c.Error = err
//	return
//}

var checkErrorInstance checkError

func execWGdown(interfaceName, interfaceIP, defaultRouteTable string) (err error) {
	if err = checkDefaultRoute(interfaceName, defaultRouteTable); err == nil || errors.Is(err, checkErrorInstance) {
		err = deleteDefaultRoute(interfaceName, defaultRouteTable)
		if err != nil {
			return err
		}
	}
	err = deleteServerRule()
	if err != nil {
		return err
	}
	err = deleteLookupRule(defaultRouteTable)
	if err != nil {
		return err
	}
	err = linkDown(interfaceName)
	if err != nil {
		return err
	}
	err = deleteIpAddress(interfaceIP, interfaceName)
	if err != nil {
		return err
	}
	return deleteInterface(interfaceName)
}

func execWGup(interfaceName, privateKey, publicKey, endpointIP, interfaceIP, defaultWGPort, defaultRouteTable string) (err error) {
	if !checkInterface(interfaceName) {
		err = addInterface(interfaceName)
		if err != nil {
			return err
		}
	}

	err = checkPrivateKey(interfaceName, privateKey)
	if errors.Is(err, checkErrorInstance) {
		err = setPrivateKey(interfaceName, privateKey)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkIfPeerOk(interfaceName, publicKey, endpointIP, defaultWGPort)
	if errors.Is(err, checkErrorInstance) {
		err = setPeer(interfaceName, publicKey, endpointIP, defaultWGPort)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkIfAddressOk(interfaceName, interfaceIP)
	if errors.Is(err, checkErrorInstance) {
		err = setAddress(interfaceName, interfaceIP)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkIfLinkDown(interfaceName)
	if errors.Is(err, checkErrorInstance) {
		err = linkDown(interfaceName)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}
	err = linkUp(interfaceName)
	if err != nil {
		return err
	}

	err = checkDefaultRoute(interfaceName, defaultRouteTable)
	if errors.Is(err, checkErrorInstance) {
		err = addDefaultRoute(interfaceName, defaultRouteTable)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	err = checkServerRule(endpointIP)
	if errors.Is(err, checkErrorInstance) {
		err = deleteServerRule()
		if err != nil {
			return err
		}
		err = addServerRule(endpointIP)
		if err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	if err != nil {
		return err
	}
	err = checkLookupRule(defaultRouteTable)
	if errors.Is(err, checkErrorInstance) {
		return addLookupRule(defaultRouteTable)
	}
	return err
}

func getCurrentEndpointIpFromRulesList(lines string) (string, error) {
	re := regexp.MustCompile(`.*219.*from all to ([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}) lookup main\n`)
	//for _, line := range lines {
	if matches := re.FindStringSubmatch(lines); matches != nil {
		return libraryStrings.Trim(matches[1]), nil
	}
	//}
	return "", fmt.Errorf("no IP found:\n %s", lines)
}

func getRules() (string, error) {
	out, _, err := libraryExec.Run(
		"ip",
		"rule",
		"list")
	if err != nil {
		return "", libraryErrors.WrapError("ip rule list", err)
	}
	return out, nil
}

func addDefaultRoute(interfaceName, defaultRouteTable string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"route",
		"add",
		"default",
		"dev",
		interfaceName,
		"table", defaultRouteTable)
	if err != nil {
		return libraryErrors.WrapError("add default route", err)
	}
	return nil
}

func addInterface(interfaceName string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"link",
		"add",
		"dev",
		interfaceName,
		"type",
		"wireguard")
	if err != nil {
		return libraryErrors.WrapError("add interface", err)
	}
	return nil
}
func addLookupRule(defaultRouteTable string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"rule",
		"add",
		"lookup",
		defaultRouteTable,
		"priority",
		"220")
	if err != nil {
		return libraryErrors.WrapError("add lookup rule", err)
	}
	return nil
}
func addServerRule(endpointIP string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"rule",
		"add",
		"to",
		endpointIP,
		"table",
		"main",
		"priority",
		"219")
	if err != nil {
		return libraryErrors.WrapError("add server rule", err)
	}
	return nil
}
func setAddress(interfaceName, interfaceIP string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"address",
		"add",
		interfaceIP,
		"dev",
		interfaceName)
	if err != nil {
		return libraryErrors.WrapError("set ip address", err)
	}
	return nil
}
func setPeer(interfaceName, publicKey, endpointIP, defaultWGPort string) error {
	_, _, err := libraryExec.Run(
		"wg",
		"set",
		interfaceName,
		"peer",
		publicKey,
		"endpoint",
		endpointIP+":"+defaultWGPort,
		"allowed-ips",
		"0.0.0.0/0")
	if err != nil {
		return libraryErrors.WrapError("set peer", err)
	}
	return nil
}
func setPrivateKey(interfaceName, privateKey string) error {
	_, _, err := libraryExec.RunWithStdin(privateKey,
		"wg",
		"set",
		interfaceName,
		"private-key",
		"/dev/stdin")
	if err != nil {
		return libraryErrors.WrapError("set private key", err)
	}
	return nil
}

func deleteInterface(interfaceName string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"link",
		"delete",
		"dev",
		interfaceName)
	if err != nil {
		return libraryErrors.WrapError("delete interface", err)
	}
	return nil
}

func deleteDefaultRoute(interfaceName, defaultRouteTable string) error {
	_, _, err := libraryExec.Run(
		"ip",
		"route",
		"delete",
		"default",
		"dev",
		interfaceName,
		"table",
		defaultRouteTable)
	if err != nil {
		return libraryErrors.WrapError("delete default route", err)
	}
	return nil
}

func deleteServerRule() error {
	rules, err := getRules()
	if err != nil {
		return libraryErrors.WrapError("delete rule for server", err)
	}
	ip, err := getCurrentEndpointIpFromRulesList(rules)
	if err != nil {
		//nolint:nilerr
		return nil // because this only can happen if no ip found -> no rule exist
	}
	_, _, err = libraryExec.Run("ip",
		"rule",
		"delete",
		"to",
		ip,
		"table",
		"main",
		"priority",
		"219")
	if err != nil {
		return libraryErrors.WrapError("delete rule for server", err)
	}
	return nil
}
func deleteLookupRule(defaultRouteTable string) error {
	_, _, err := libraryExec.Run("ip",
		"rule",
		"delete",
		"lookup",
		defaultRouteTable,
		"priority",
		"220")
	if err != nil {
		return libraryErrors.WrapError("delete lookup rule", err)
	}
	return nil
}
func deleteIpAddress(interfaceIP, interfaceName string) error {
	_, _, err := libraryExec.Run("ip",
		"address",
		"del",
		interfaceIP,
		"dev",
		interfaceName)
	if err != nil {
		return libraryErrors.WrapError("delete address from the interface", err)
	}
	return nil
}

func linkUp(interfaceName string) error {
	_, _, err := libraryExec.Run("ip",
		"link",
		"set",
		"up",
		"dev",
		interfaceName)
	if err != nil {
		return libraryErrors.WrapError("up interface", err)
	}
	return nil
}
func linkDown(interfaceName string) error {
	_, _, err := libraryExec.Run("ip",
		"link",
		"set",
		"down",
		"dev",
		interfaceName)
	if err != nil {
		return libraryErrors.WrapError("down interface", err)
	}
	return nil
}

func checkIfLinkDown(interfaceName string) error {
	out, _, err := libraryExec.Run("ip",
		"link",
		"show",
		interfaceName,
		"up")
	if err != nil {
		return libraryErrors.WrapError("check interface", err)
	}
	if len(out) == 0 {
		return nil
	}
	return checkError{}
}

func checkInterface(interfaceName string) bool {
	_, code, _ := libraryExec.Run(
		"ip",
		"link",
		"show",
		interfaceName)
	return code != 1
}

func checkDefaultRoute(interfaceName, defaultRouteTable string) error {
	out, _, err := libraryExec.Run("ip",
		"route",
		"show",
		"table",
		defaultRouteTable)
	if err != nil {
		if strings.Contains(out, "table does not exist") {
			return checkError{}
		}
		return libraryErrors.WrapError("check default route", err)
	}
	if strings.Contains(out, "default dev "+interfaceName) {
		return nil
	}
	return checkError{}
}

func checkLookupRule(defaultRouteTable string) error {
	re := regexp.MustCompile(`.*220.*from all lookup ` + defaultRouteTable + `\n`)
	rules, err := getRules()
	if err != nil {
		return err
	}
	if matches := re.FindStringSubmatch(rules); matches != nil {
		return nil
	}
	return checkError{}
}
func checkServerRule(endpointIP string) error {
	rules, err := getRules()
	if err != nil {
		return libraryErrors.WrapError("check server rule", nil)
	}
	ip, err := getCurrentEndpointIpFromRulesList(rules)
	if err != nil || ip != endpointIP {
		return checkError{}
	}
	return nil
}

func checkIfAddressOk(interfaceName, interfaceIP string) error {
	out, _, err := libraryExec.Run("ip",
		"address",
		"show",
		interfaceName)
	if err != nil {
		return libraryErrors.WrapError("check ip address", err)
	}
	if strings.Contains(out, interfaceIP) {
		return nil
	}
	return checkError{}
}
func checkIfPeerOk(interfaceName, publicKey, endpointIP, defaultWGPort string) error {
	_, _, err := libraryExec.Run("wg",
		"set",
		interfaceName,
		"peer",
		publicKey,
		"endpoint",
		endpointIP+":"+defaultWGPort,
		"allowed-ips",
		"0.0.0.0/0")
	if err != nil {
		return libraryErrors.WrapError("get peer", err)
	}
	return nil
}
func checkPrivateKey(interfaceName, privateKey string) error {
	currentPrivateKey, _, err := libraryExec.Run("wg",
		"show",
		interfaceName,
		"private-key")
	if err != nil {
		return libraryErrors.WrapError("check private key", err)
	}
	if currentPrivateKey == privateKey {
		return nil
	}
	return checkError{}
}
