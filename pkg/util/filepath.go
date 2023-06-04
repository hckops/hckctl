package util

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

const (
	defaultDirectoryMod os.FileMode = 0755
)

func CreateBaseDir(path string) error {
	return createBaseDirMod(path, defaultDirectoryMod)
}

func createBaseDirMod(path string, mod os.FileMode) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, mod); err != nil {
			return errors.Wrapf(err, "unable to create dir %s", dir)
		}
	}
	return nil
}

func ReadFile(path string) (string, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return "", errors.Wrapf(err, "unable to read file %s", path)
	}
	return string(data), nil
}