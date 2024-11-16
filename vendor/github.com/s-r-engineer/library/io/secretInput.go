package libraryIO

import (
	"fmt"
	"strings"
	"syscall"

	"golang.org/x/term"
)

func ReadSecretInput(message string) (string, error) {
	fmt.Println(message)
	pinBytes, err := term.ReadPassword(int(syscall.Stdin))
	return strings.TrimSpace(string(pinBytes)), err
}
