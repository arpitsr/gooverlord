package indexer

import (
	"com.ak.gooverlord/models"
)

// we can extend multiple backend implementation based on req or POC results
type IndexerBackend interface {
	GetClient(node string) *any
	IndexEntry(le models.LogEntry) error
}
