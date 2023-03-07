package models

type Query struct {
	SearchQuery string `json:"q"`
	Index       string `json:"index"`
}
