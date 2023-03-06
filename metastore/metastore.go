package metastore

type Backend interface {
	Get(key string) (string, error)
	Set(key string, value interface{}) error
}
