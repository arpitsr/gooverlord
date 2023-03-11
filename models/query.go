package models

import "github.com/meilisearch/meilisearch-go"

type Query struct {
	SearchQuery   string `json:"q"`
	Index         string `json:"index"`
	SearchRequest meilisearch.SearchRequest
}
