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

var debugVar = "0"

var tokenPath = getHomeDir() + "/.config/norrvpn"
var tokenFullPath = tokenPath + "/token.json"

func debug() bool {
	return debugVar == "1"
}

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
	wrapIfError("", err)
	return trim(string(pinBytes))
}

func trim(s string) string {
	return strings.TrimSpace(s)
}

func wrapIfError(msg string, err error) error {
	if err != nil {
		return fmt.Errorf("%s -> %w", msg, err)
	}
	return err
}

// func wrapIfError(msg string, err error) error {
// 	if err != nil {
// 		if debug() {
// 			var (
// 				pc        uintptr
// 				file      string
// 				line      int
// 				traceback []string
// 			)
// 			for i := 1; i > 0; i++ {
// 				pc, file, line, _ = runtime.Caller(i)
// 				if pc == 0 {
// 					break
// 				}
// 				function := runtime.FuncForPC(pc).Name()
// 				traceback = append(traceback, fmt.Sprintf("%s:%d %s", file, line, function))
// 			}
// 			slices.Reverse(traceback)
// 			traceback = append(traceback, msg, err.Error())
// 			return fmt.Errorf(strings.Join(traceback, "\n↓\n"))
// 		}
// 		return fmt.Errorf("%s -> %w", msg, err)
// 	}
// 	return err
// }
