package storage

type Mapper interface {
	Get(key, command string) ([]string, bool)
	Store(key, value, command string) error
	Delete(key, value, command string) error
	Close() error
}
