package indexer

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"sync"
	"time"

	"com.ak.gooverlord/config"
	"com.ak.gooverlord/metastore"
	"com.ak.gooverlord/models"
	"com.ak.gooverlord/partitioner"
	"github.com/meilisearch/meilisearch-go"
)

var instance *meilisearchBackend

type meilisearchBackend struct {
	lock     sync.RWMutex
	connsMap sync.Map
}

func GetInstance() *meilisearchBackend {
	if instance == nil {
		UpdateInstance()
	}
	return instance
}

func init() {
	instance = &meilisearchBackend{}
	UpdateInstance()
}

func UpdateInstance() {
	instance.lock.Lock()
	defer instance.lock.Unlock()
	instance = &meilisearchBackend{}
	nodes := partitioner.GetConsistentHashRing().RealNodesSet
	for node := range nodes {
		client := meilisearch.NewClient(meilisearch.ClientConfig{
			Host:   node,
			APIKey: os.Getenv(config.MASTER_KEY),
		})
		instance.connsMap.Store(node, client)
	}
}

func (mb *meilisearchBackend) GetClient(node string) *meilisearch.Client {
	client, _ := instance.connsMap.Load(node)
	return client.(*meilisearch.Client)
}

func (mb *meilisearchBackend) IndexEntries(entries []models.LogEntry) {
	// Map to assign which log entries should go to which node based on
	// consistent hashing ring
	nodeToLogEntriesMap := make(map[string][]models.LogEntry, 0)

	for _, entry := range entries {
		chr := partitioner.GetConsistentHashRing()
		// Here we should do defined sharding
		shard := rand.Intn(config.DEFAULT_NO_OF_SHARDS)
		node, err := chr.GetNode(entry.Appname + fmt.Sprintf("%d", shard))
		if err != nil {
			log.Printf("Could not get node %s", err)
			return
		}
		nodeToLogEntriesMap[node] = append(nodeToLogEntriesMap[node], entry)
	}

	// for each node index the entries by fetching client for that meilisearch node
	for node, logEntries := range nodeToLogEntriesMap {
		client := mb.GetClient(node)
		indexToEntryMap := make(map[string][]models.LogEntry, 0)
		for _, val := range logEntries {
			t, err := time.Parse(time.RFC3339, val.Timestamp)
			if err != nil {
				fmt.Println(err)
				continue
			}

			// Format the hour and minute as a string
			hourMinute := t.Format(config.DATE_INDEX_FORMAT)
			mb.IsIndexReady(hourMinute, node) // Setting up index and filter attr
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

func (imb *meilisearchBackend) IsIndexReady(index string, node string) bool {
	_, b := metastore.GetInMemoryBacked().Get(index)
	if !b {
		client := imb.GetClient(node)
		client.CreateIndex(&meilisearch.IndexConfig{
			Uid:        index,
			PrimaryKey: "id",
		})
		//TODO: move this to config
		client.Index(index).UpdateFilterableAttributes(&[]string{
			"appname",
			"hostname",
		})
		metastore.GetInMemoryBacked().Set(index, true)
	}
	return b
}
