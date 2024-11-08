package main

import (
	"os/exec"
	"regexp"
	"strings"
)

const defaultWGPort = "51820"

func getEndpointIP(lines []string) string {
	re := regexp.MustCompile(`^219.*from all to ([0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}) lookup main$`)
	for _, line := range lines {
		if matches := re.FindStringSubmatch(line); matches != nil {
			return trim(matches[1])
		}
	}
	panicer("no IP found:\n" + strings.Join(lines, "\n"))
	return ""
}
func run(command ...string) (string, error, int) {
	cmd := exec.Command(command[0], command[1:]...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		dumper(command)
		dumper(string(output))
		panicer(err)
	}
	return string(output), err, cmd.ProcessState.ExitCode()
}

func execWGdown(interfaceName, interfaceIP string) {
	run("ip", "route", "delete", "default", "dev", interfaceName, "table", "212450")
	out, _, _ := run("ip", "rule", "show")
	run("ip", "rule", "delete", "to", getEndpointIP(strings.Split(out, "\n")), "table", "main", "priority", "219")
	run("ip", "rule", "delete", "lookup", "212450", "priority", "220")
	run("ip", "link", "set", "down", "dev", interfaceName)
	run("ip", "address", "del", interfaceIP, "dev", interfaceName)
	run("ip", "link", "delete", "dev", interfaceName)
}

func execWGup(interfaceName, privateKey, publicKey, endpointIP, interfaceIP string) {
	var cmd *exec.Cmd
	cmd = exec.Command("ip", "link", "show", interfaceName)
	cmd.Run()
	if cmd.ProcessState.ExitCode() == 1 {
		cmd = exec.Command("ip", "link", "add", "dev", interfaceName, "type", "wireguard")
		panicer(cmd.Run())
	}
	cmd = exec.Command("wg", "set", interfaceName, "private-key", "/dev/stdin")
	cmd.Stdin = strings.NewReader(privateKey)
	b, err := cmd.CombinedOutput()
	if err != nil {
		dumper(string(b))
		panicer(err)
	}
	run("wg", "set", interfaceName, "peer", publicKey, "endpoint", endpointIP+":"+defaultWGPort, "allowed-ips", "0.0.0.0/0")
	run("ip", "address", "add", interfaceIP, "dev", interfaceName)
	run("ip", "link", "set", "up", "dev", interfaceName)
	run("ip", "route", "add", "default", "dev", interfaceName, "table", "212450")
	run("ip", "rule", "add", "to", endpointIP, "table", "main", "priority", "219")
	run("ip", "rule", "add", "lookup", "212450", "priority", "220")
}
