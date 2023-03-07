package query

import (
	"fmt"
	"log"

	"com.ak.gooverlord/indexer"
	"com.ak.gooverlord/models"
	"com.ak.gooverlord/partitioner"
	"github.com/meilisearch/meilisearch-go"
)

func Search(query models.Query) []*meilisearch.SearchResponse {
	chr := partitioner.GetConsistentHashRing()
	nodes := chr.RealNodesSet
	fmt.Printf("Query Model: %v\n", query)
	var results []*meilisearch.SearchResponse
	for node := range nodes {
		result, err := indexer.GetMeilisearchClient(node).Index(query.Index).Search(query.SearchQuery, &meilisearch.SearchRequest{})
		if err != nil {
			log.Println(err.Error())
			continue
		}
		results = append(results, result)
	}
	return results
}
