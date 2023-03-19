package query

import (
	"log"

	"com.ak.gooverlord/indexer"
	"com.ak.gooverlord/models"
	"com.ak.gooverlord/partitioner"
	"github.com/meilisearch/meilisearch-go"
)

type Planner interface {
	FullTextSearch(query models.Query) interface{}
}

// Different planner strategy needs to be defined based on query type
// dumb all scan strategy is used here
// planning and query optimization can be done based on metastore for indexed data
// this needs to be updated to generic response
func FullTextSearch(query models.Query) []*meilisearch.SearchResponse {
	chr := partitioner.GetConsistentHashRing()
	nodes := chr.RealNodesSet
	var results []*meilisearch.SearchResponse
	for node := range nodes {
		result, err := indexer.GetInstance().GetClient(node).Index(query.Index).Search(query.SearchQuery, &query.SearchRequest)
		if err != nil {
			log.Println(err.Error())
			continue
		}
		results = append(results, result)
	}
	return results
}
