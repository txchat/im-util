package wallet

/*
	文件的读写驱动，用于实现将用户从文件格式化读出和写入
*/
import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"strings"
)

type FSDriver struct {
	filePath  string
	formatter Formatter

	rows []string
}

func NewFSDriver(filePath string, formatter Formatter) *FSDriver {
	return &FSDriver{
		filePath:  filePath,
		formatter: formatter,
	}
}

func (d *FSDriver) Load() ([]*Metadata, error) {
	// load
	var err error
	d.rows, err = d.readLines(d.filePath)
	if err != nil {
		return nil, err
	}

	ret := make([]*Metadata, len(d.rows))
	for i, row := range d.rows {
		ret[i], err = d.formatter.FromRow(row)
		if err != nil {
			return nil, err
		}
	}
	return ret, nil
}

func (d *FSDriver) Save(metadata []*Metadata) error {
	f, err := os.OpenFile(d.filePath, os.O_CREATE|os.O_TRUNC|os.O_RDWR, os.ModePerm)
	if err != nil {
		panic(err)
	}
	defer f.Close()
	rb := bufio.NewWriter(f)
	defer rb.Flush()

	d.rows = make([]string, len(metadata))
	for i, md := range metadata {
		item := d.formatter.ToRow(md)
		_, err = rb.WriteString(item + "\n")
		if err != nil {
			return err
		}
		d.rows[i] = item
	}
	return nil
}

func (d *FSDriver) readLines(uri string) ([]string, error) {
	f, err := os.OpenFile(uri, os.O_RDONLY|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	b, err := ioutil.ReadAll(f)
	if err != nil {
		return nil, err
	}
	list := strings.Split(string(b), "\n")
	if len(list) > 0 && list[len(list)-1] == "" {
		return list[:len(list)-1], nil
	}
	return list, nil
}

func LoadMetadata(readPath, readSplit string) ([]*Metadata, error) {
	if readPath == "" || readSplit == "" {
		return nil, errors.New("readPath or readSplit can not empty")
	}
	formatter := NewSplitFormatter(readSplit)
	readDriver := NewFSDriver(readPath, formatter)
	return readDriver.Load()
}
