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
	dir := filepath.Dir(path)
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, defaultDirectoryMod); err != nil {
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

func CopyFile(from string, to string) (int64, error) {

	if err := CreateBaseDir(to); err != nil {
		return -1, err
	}

	reader, err := os.Open(from)
	if err != nil {
		return -1, errors.Wrapf(err, "unable to open source file %s", from)
	}
	defer reader.Close()

	writer, err := os.Create(to)
	if err != nil {
		return -1, errors.Wrapf(err, "unable to open destination file %s", to)
	}
	defer writer.Close()

	return writer.ReadFrom(reader)
}

func DeleteFile(path string) error {
	err := os.Remove(path)
	if err != nil {
		return errors.Wrapf(err, "unable to delete file %s", path)
	}
	return nil
}

func DeleteDir(path string) error {
	err := os.RemoveAll(path)
	if err != nil {
		return errors.Wrapf(err, "unable to delete dir %s", path)
	}
	return nil
}
