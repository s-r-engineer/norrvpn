package main

import (
	"flag"

	libraryLogging "github.com/s-r-engineer/library/logging"
)

func main() {
	var err error

	server := flag.Bool("server", false, "Set server mode")
	debug := flag.Bool("debug", false, "Set debug on")
	flag.Parse()
	if !*debug {
		libraryLogging.InitLogger(true)
	}
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
