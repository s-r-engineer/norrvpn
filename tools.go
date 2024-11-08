package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/user"
	"strings"
	"syscall"

	"github.com/davecgh/go-spew/spew"
	"golang.org/x/term"
)

var tokenPath = getHomeDir() + "/.config/norrvpn"
var tokenFullPath = tokenPath + "/token.json"

func panicer(err any) {
	if err != nil {
		panic(err)
	}
}

func dumper(v any) {
	spew.Dump(v)
}

func getHomeDir() string {
	if val, ok := os.LookupEnv("SUDO_HOME"); ok {
		return val
	}
	usr, err := user.Current()
	if err != nil {
		panicer(err)
	}
	return usr.HomeDir
}

func randString(n int) string {
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func readSecretInput(message string) string {
	fmt.Println(message)
	pinBytes, err := term.ReadPassword(int(syscall.Stdin))
	panicer(err)
	return trim(string(pinBytes))
}

func trim(s string) string {
	return strings.TrimSpace(s)
}
