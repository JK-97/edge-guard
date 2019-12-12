package store

type KV interface {
	Get(key string) ([]byte, error)
	GetDefault(key string, value []byte) ([]byte, error)
	Set(key string, value []byte) error
}
