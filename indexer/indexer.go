package indexer

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"com.ak.gooverlord/models"
	"com.ak.gooverlord/partitioner"
	"com.ak.gooverlord/utils"
	"github.com/meilisearch/meilisearch-go"
)

type IndexerBackend interface {
	GetMeilisearchClient(node string) *meilisearch.Client
	IndexEntry(le models.LogEntry) error
}

var once sync.Once
var searchBackend *meilisearchBackend

type meilisearchBackend struct {
	connsMap map[string]*meilisearch.Client
}

func GetMeilisearchClient(node string) *meilisearch.Client {
	once.Do(func() {
		searchBackend = &meilisearchBackend{connsMap: make(map[string]*meilisearch.Client, 0)}
		nodes := partitioner.GetConsistentHashRing().RealNodesSet
		for node := range nodes {
			client := meilisearch.NewClient(meilisearch.ClientConfig{
				Host:   node,
				APIKey: os.Getenv(utils.MASTER_KEY),
			})
			searchBackend.connsMap[node] = client
		}
	})
	return searchBackend.connsMap[node]
}

func IndexEntries(entries []models.LogEntry) {
	// Map to assign which log entries should go to which node based on
	// consistent hashing ring
	nodeToLogEntriesMap := make(map[string][]models.LogEntry, 0)

	for _, entry := range entries {
		chr := partitioner.GetConsistentHashRing()
		node, err := chr.GetNode(entry.Appname + entry.Hostname)
		if err != nil {
			log.Println("Could not get node", err)
			return
		}
		nodeToLogEntriesMap[node] = append(nodeToLogEntriesMap[node], entry)
	}

	fmt.Println(nodeToLogEntriesMap)

	// for each node index the entries by fetching client for that meilisearch node
	for node, logEntries := range nodeToLogEntriesMap {
		client := GetMeilisearchClient(node)
		indexToEntryMap := make(map[string][]models.LogEntry, 0)
		for _, val := range logEntries {
			t, err := time.Parse(time.RFC3339, val.Timestamp)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Format the hour and minute as a string
			hourMinute := t.Format("15-04")
			indexToEntryMap[hourMinute] = append(indexToEntryMap[hourMinute], val)
		}

		for key, val := range indexToEntryMap {
			index := client.Index(key)
			_, err := index.AddDocuments(val)
			if err != nil {
				fmt.Println(err)
				continue
			}
		}
	}
}
