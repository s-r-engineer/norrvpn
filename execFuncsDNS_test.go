package main

import (
	libraryLogging "github.com/s-r-engineer/library/logging"
	"testing"
)

func TestCheckDNSResolver(t *testing.T) {
	err := checkDNSResolver()
	libraryLogging.Dumper(err, string(resolver))
}
