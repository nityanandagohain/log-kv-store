package store

import (
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type FileStore struct {
	writeFile *os.File
	readFile  *os.File
	index     map[string]int64
}

func NewFileStore(path string) (Store, error) {
	writeFile, err := os.OpenFile(path+"/data.store", os.O_APPEND|os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	readFile, err := os.OpenFile(path+"/data.store", os.O_RDONLY, os.ModePerm)
	if err != nil {
		return nil, err
	}

	index := map[string]int64{}
	err = RecreteIndex(readFile, &index)
	if err != nil {
		return nil, err
	}

	return &FileStore{
		writeFile,
		readFile,
		index,
	}, nil
}

type LenData struct {
	Lkey int32 `json:"lkey"`
	LVal int32 `json:"lval"`
}

// recreate the index from the file
func RecreteIndex(readFile *os.File, index *map[string]int64) error {
	offset, err := readFile.Seek(0, io.SeekStart)
	if err != nil {
		if err == io.EOF {
			return nil
		}
		return err
	}

	for {
		data := LenData{}
		if err := binary.Read(readFile, binary.LittleEndian, &data); err != nil {
			if err == io.EOF {
				break
			}
			return err
		}

		keyVal := make([]byte, data.Lkey+data.LVal)
		_, err = readFile.Read(keyVal)
		if err != nil {
			return err
		}

		key := string(keyVal[0:data.Lkey])
		_ = string(keyVal[data.Lkey : data.Lkey+data.LVal])
		(*index)[key] = offset
		curr, err := readFile.Seek(0, io.SeekCurrent)
		if err != nil {
			return err
		}
		offset += curr
	}

	return nil
}

func (f *FileStore) Put(key, val string) error {
	offset, err := f.writeFile.Seek(0, io.SeekEnd)
	if err != nil {
		return err
	}

	// storing the offset
	f.index[key] = offset

	if err = binary.Write(f.writeFile, binary.LittleEndian, LenData{
		Lkey: int32(len(key)),
		LVal: int32(len(val)),
	}); err != nil {
		return err
	}

	if _, err := f.writeFile.Write([]byte(fmt.Sprintf("%s%s", key, val))); err != nil {
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
