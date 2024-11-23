package main

import (
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	libraryErrors "github.com/s-r-engineer/library/errors"
	libraryExec "github.com/s-r-engineer/library/exec"
	libraryStrings "github.com/s-r-engineer/library/strings"
)

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

func execWGdown(interfaceName, interfaceIP string) {
	libraryExec.Run("ip", "route", "delete", "default", "dev", interfaceName, "table", "212450")
	out, _, _ := libraryExec.Run("ip", "rule", "show")
	libraryExec.Run("ip", "rule", "delete", "to", getEndpointIP(strings.Split(out, "\n")), "table", "main", "priority", "219")
	libraryExec.Run("ip", "rule", "delete", "lookup", "212450", "priority", "220")
	libraryExec.Run("ip", "link", "set", "down", "dev", interfaceName)
	libraryExec.Run("ip", "address", "del", interfaceIP, "dev", interfaceName)
	libraryExec.Run("ip", "link", "delete", "dev", interfaceName)
}

func execWGup(interfaceName, privateKey, publicKey, endpointIP, interfaceIP string) error {
	var cmd *exec.Cmd
	cmd = exec.Command("ip", "link", "show", interfaceName)
	cmd.Run()
	if cmd.ProcessState.ExitCode() == 1 {
		cmd = exec.Command("ip", "link", "add", "dev", interfaceName, "type", "wireguard")
		libraryErrors.Panicer(cmd.Run())
	}
	cmd = exec.Command("wg", "set", interfaceName, "private-key", "/dev/stdin")
	cmd.Stdin = strings.NewReader(privateKey)
	_, err := cmd.CombinedOutput()
	if err != nil {
		libraryErrors.Panicer(err)
	}
	libraryExec.Run("wg", "set", interfaceName, "peer", publicKey, "endpoint", endpointIP+":"+defaultWGPort, "allowed-ips", "0.0.0.0/0")
	libraryExec.Run("ip", "address", "add", interfaceIP, "dev", interfaceName)
	libraryExec.Run("ip", "link", "set", "up", "dev", interfaceName)
	libraryExec.Run("ip", "route", "add", "default", "dev", interfaceName, "table", "212450")
	libraryExec.Run("ip", "rule", "add", "to", endpointIP, "table", "main", "priority", "219")
	libraryExec.Run("ip", "rule", "add", "lookup", "212450", "priority", "220")
	return nil
}

func setAddress(interfaceName, interfaceIP string) error {
	_, _, err := libraryExec.Run("ip", "address", "add", interfaceIP, "dev", interfaceName)
	return libraryErrors.WrapError(fmt.Sprintf("setAddress -> %s, %s", interfaceName, interfaceIP), err)
}
