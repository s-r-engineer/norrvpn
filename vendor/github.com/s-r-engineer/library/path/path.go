package libraryPath

import (
	"os"
	"os/user"
)

func GetHomeDir() (string, error) {
	if val, ok := os.LookupEnv("SUDO_HOME"); ok {
		return val, nil
	}
	usr, err := user.Current()
	return usr.HomeDir, err
}
