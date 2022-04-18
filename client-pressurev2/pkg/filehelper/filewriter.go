package filehelper

import (
	"io"
	"os"
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
