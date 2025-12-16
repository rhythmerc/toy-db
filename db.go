package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

type DB interface {
	read(key string) (string, error)
	write(key string, value string) error
}

type HashDB struct {
	file    *os.File
	indices map[string]int64
}

func (db *HashDB) Read(key string) (string, error) {
	idx, ok := db.indices[key]
	if !ok {
		return "", &KeyNotFoundError{key}
	}

	bytes := make([]byte, 100)
	db.file.Seek(idx+int64(len(key))+1, 0)

	for buf := make([]byte, 1); ; {
		if _, err := db.file.Read(buf); err != nil {
			if err == io.EOF {
				break
			} else {
				return "", err
			}
		}

		char := buf[0]
		if char == '\n' {
			break
		}

		bytes = append(bytes, char)
	}

	return string(bytes), nil
}

func (db *HashDB) Write(key string, value string) error {
	value = strings.ReplaceAll(value, "\n", " ")
	value = strings.TrimSpace(value)

	offset, err := db.file.Seek(0, 2)
	if err != nil {
		return err
	}

	entry := fmt.Sprintf("%s,%s\n", key, value)
	if _, err := db.file.WriteString(entry); err != nil {
		return err
	}

	db.indices[key] = offset
	return nil
}

func (db *HashDB) PopulateIndices() error {
	prevOffset, err := db.file.Seek(0, 1)
	if err != nil {
		return err
	}

	_, err = db.file.Seek(0, 0)
	if err != nil {
		return err
	}

	reader := bufio.NewReader(db.file)
	offset := int64(0)

	for {
		bytes, err := reader.ReadBytes('\n')
		if err != nil {
			if err == io.EOF {
				break
			}

			return err
		}

		key, _, valid := strings.Cut(string(bytes), ",")
		if !valid {
			offset += int64(len(bytes))
			continue
		}

		db.indices[key] = offset
		offset += int64(len(bytes))
	}

	_, err = db.file.Seek(prevOffset, 0)
	return err
}

func NewHashDB(dbpath string) (*HashDB, error) {
	file, err := os.OpenFile("db.txt", os.O_CREATE|os.O_RDWR, 0644)
	if err != nil {
		return nil, err
	}

	db := HashDB{
		file:    file,
		indices: map[string]int64{},
	}

	if err := db.PopulateIndices(); err != nil {
		return nil, err
	}

	return &db, nil
}
