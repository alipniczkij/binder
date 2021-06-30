package storage

type Mapper interface {
	Get(id, command string) ([]string, bool)
	Store(id, value, command string) error
	Delete(id, value, command string) error
	Close() error
}
