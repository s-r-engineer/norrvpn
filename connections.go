package main

import (
	"fmt"
	"io"
	"net"
	"os"

	libraryErrors "github.com/s-r-engineer/library/errors"
)

const socketPath = "/var/run/norrvpn.sock"

func parseConnection(conn net.Conn) (fullData []byte, err error) {
	var data = make([]byte, defaultReadBlockSize)

	var counter int

	for {
		counter, err = conn.Read(data)
		if err != nil && err != io.EOF {
			return
		}

		fullData = append(fullData, data[:counter]...)

		if counter < defaultReadBlockSize {
			break
		}
	}
	return fullData, err
}

func getSocketListener() (net.Listener, error) {
	if _, err := os.Stat(socketPath); err == nil {
		err = os.Remove(socketPath)
		if err != nil {
			return nil, err
		}
	}
	listener, err := net.Listen("unix", socketPath)
	if err != nil {
		return nil, err
	}
	err = os.Chmod(socketPath, 0777)
	if err != nil {
		libraryErrors.Errorer(listener.Close())
		return nil, fmt.Errorf("failed to create Unix socket: %v", err)
	}
	return listener, err
}

func getSocketConnector() (net.Conn, error) {
	return net.Dial("unix", socketPath)
}
