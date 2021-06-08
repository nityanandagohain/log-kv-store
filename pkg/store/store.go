package store

type Store interface {
	Put(key, val string) error
	Get(key string) (string, error)
}
