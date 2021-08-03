package storage

type Mapper interface {
	Get(key, command string) ([]string, bool)
	GetKeys(string) ([]string, error)
	Store(key, value, command string) error
	Delete(key, value, command string) error
	Close() error
}
