package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type Store interface {
	Put(key, val string) error
	Get(key string) (string, error)
}

type FileStore struct {
	writeFile *os.File
	readFile  *os.File
	index     map[string]int64
}

func NewFileStore(path string) (Store, error) {
	readFile, err := os.OpenFile(path+"/data.store", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return nil, err
	}

	writeFile, err := os.OpenFile(path+"/data.store", os.O_RDONLY, 0644)
	if err != nil {
		return nil, err
	}

	return &FileStore{
		readFile,
		writeFile,
		map[string]int64{},
	}, nil
}

type LenData struct {
	Lkey int32 `json:"lkey"`
	LVal int32 `json:"lval"`
}

func (f *FileStore) Put(key, val string) error {
	offset, err := f.writeFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// storing the offset
	f.index[key] = offset

	err = binary.Write(f.writeFile, binary.LittleEndian, LenData{
		Lkey: int32(len(key)),
		LVal: int32(len(val)),
	})

	if _, err := f.writeFile.Write([]byte(fmt.Sprintf("%s%s", key, val))); err != nil {
		return err
	}
	if err != nil {
		return err
	}

	return nil
}

func (f *FileStore) Get(key string) (string, error) {
	var value string
	var offset int64
	var ok bool
	if offset, ok = f.index[key]; !ok {
		return value, errors.New("key not found")
	}

	_, err := f.readFile.Seek(offset, io.SeekStart)
	if err != nil {
		return value, err
	}

	data := LenData{}
	if err := binary.Read(f.readFile, binary.LittleEndian, &data); err != nil {
		return value, err
	}

	keyVal := make([]byte, data.Lkey+data.LVal)
	_, err = f.readFile.Read(keyVal)
	if err != nil {
		return value, err
	}

	storedKey := string(keyVal[0:data.Lkey])
	value = string(keyVal[data.Lkey : data.Lkey+data.LVal])

	if strings.Compare(storedKey, key) != 0 {
		return value, errors.New("key doest match")
	}

	return value, nil
}
