package util

import (
	"io"
	"io/ioutil"
	"os"
	"strings"
)

func WriteFile(filepath string) (io.Writer, func(), error) {
	fd, err := os.OpenFile(filepath, os.O_WRONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	return fd, func() {
		fd.Close()
	}, nil
}

func ReadFile(filepath string) (io.Reader, func(), error) {
	fd, err := os.OpenFile(filepath, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, nil, err
	}
	return fd, func() {
		fd.Close()
	}, nil
}

func ReadAllLines(filepath string) ([]string, error) {
	fd, closer, err := ReadFile(filepath)
	if err != nil {
		return nil, err
	}
	defer closer()
	b, err := ioutil.ReadAll(fd)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(b), "\n")
	if len(list) > 0 && list[len(list)-1] == "" {
		return list[:len(list)-1], nil
	}
	return list, nil
}
