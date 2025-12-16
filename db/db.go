package db

type DB interface {
	Read(key string) (string, error)
	Write(key string, value string) error
}
