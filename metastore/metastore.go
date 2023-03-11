package metastore

type Backend interface {
	IsIndexReady(index string) error
}
