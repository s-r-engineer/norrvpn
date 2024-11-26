package libraryIO

import (
	"os"
	"strings"
)

const create = os.O_CREATE
const appendData = os.O_APPEND
const read = os.O_RDONLY
const write = os.O_WRONLY

func ReadFileToString(filePath string) (string, error) {
	data, err := ReadFileToBytes(filePath)
	return string(data), err
}

func ReadFileToBytes(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

func CreateAndOpenFileForWrite(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, create|write, 0600)
}

func CreateAndOpenFileForAppend(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, create|appendData, 0600)
}

func OpenFile(filePath string) (*os.File, error) {
	return os.OpenFile(filePath, read, 0600)
}

func CheckFileExist(filePath string) (bool, error) {
	_, err := os.Stat(filePath)
	if err == nil {
		return true, nil
	}
	return false, err
}

func CreateDirs(dirPath string) error {
	return os.MkdirAll(dirPath, 0700)
}

func CreateTempFoder(preffixAndSuffix ...string) (string, func(), error) {
	tempDir, err := os.MkdirTemp("", strings.Join(preffixAndSuffix,"*"))
	if err != nil {
		return tempDir, func() {}, err
	}
	return tempDir, func() { os.RemoveAll(tempDir) }, err
}

func CreateTempFile(tempFolder string) {
	os.CreateTemp(tempFolder,"")
}