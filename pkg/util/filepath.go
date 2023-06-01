package util

import (
	"github.com/pkg/errors"
	"os"
	"path/filepath"
)

const (
	defaultDirectoryMod os.FileMode = 0755
	defaultFileMod      os.FileMode = 0600
)

func CreateBaseDir(path string) error {
	return CreateBaseDirMod(path, defaultDirectoryMod)
}

func CreateBaseDirMod(path string, mod os.FileMode) error {
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, mod); err != nil {
			return errors.Wrapf(err, "unable to create dir %s", dir)
		}
	}
	return nil
}
