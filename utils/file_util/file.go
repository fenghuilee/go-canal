package file_util

import (
	"os"
)

func IsExist(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func IsDir(path string) bool {
	if stat, err := os.Stat(path); err == nil {
		return stat.IsDir()
	}
	return false
}

func CreateFileIfNecessary(path string) bool {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if file, err := os.Create(path); err == nil {
			file.Close()
		}
	}
	exist := IsExist(path)
	return exist
}

func MkdirIfNecessary(path string) error {
	_, err := os.Stat(path)
	if err != nil && os.IsNotExist(err) {
		if err := os.MkdirAll(path, os.ModePerm); err != nil {
			// os.Chmod(path, 0777)
			return err
		}
	}
	return nil
}
