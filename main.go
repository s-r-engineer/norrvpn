package main

import (
	"flag"

	libraryLogging "github.com/s-r-engineer/library/logging"
)

func main() {
	var err error

	server := flag.Bool("server", false, "Set server mode")
	prod := flag.Bool("prod", false, "Set debug on")
	flag.Parse()
	libraryLogging.InitLogger(*prod)
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
