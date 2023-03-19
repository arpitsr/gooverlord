package metastore

type Backend interface {
	Set(key string, value bool) (any, bool)
	Get(key string) (any, bool)
}
