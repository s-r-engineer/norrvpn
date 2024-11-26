// Command wgctrl is a testing utility for interacting with WireGuard via package
// wgctrl.
package main

import (
	"flag"
	libraryLogging "github.com/s-r-engineer/library/logging"
)

func main() {
	server := flag.Bool("server", false, "Set server mode")
	flag.Parse()
	var err error
	if *server {
		err = serverMode()
	} else {
		err = clientMode()
	}
	if err != nil {
		libraryLogging.Error(err.Error())
	} else {
		libraryLogging.Info("done")
	}
}
