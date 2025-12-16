package db

import "fmt"

type KeyNotFoundError struct {
	key string
}

func (err *KeyNotFoundError) Error() string {
	return fmt.Sprintf("key %s not found", err.key)
}

var ERR_KNF *KeyNotFoundError
