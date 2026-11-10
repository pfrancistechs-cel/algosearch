package search

import (
	"time"
)

// Document represents an indexed record in SearchLite.
type Document struct {
	ID              string    `json:"id"`
	Title           string    `json:"title"`
	Description     string    `json:"description"`
	Category        string    `json:"category"`
	Tags            []string  `json:"tags"`
	Price           float64   `json:"price"`
	PopularityScore int       `json:"popularity_score"`
	CreatedAt       time.Time `json:"created_at"`
}

// CreateDocumentInput defines the payload for indexing a new record.
type CreateDocumentInput struct {
	Title           string   `json:"title"`
	Description     string   `json:"description"`
	Category        string   `json:"category"`
	Tags            []string `json:"tags"`
	Price           float64  `json:"price"`
	PopularityScore int      `json:"popularity_score"`
}

// HighlightResult provides snippet strings with <mark> tags.
type HighlightResult struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

// Hit represents a matched document with ranking score and highlight snippets.
type Hit struct {
	ID              string          `json:"id"`
	Title           string          `json:"title"`
	Description     string          `json:"description"`
	Category        string          `json:"category"`
	Tags            []string        `json:"tags"`
	Price           float64         `json:"price"`
	PopularityScore int             `json:"popularity_score"`
	Score           float64         `json:"score"`
	Highlight       HighlightResult `json:"_highlight"`
}

// SearchQueryParams contains filters and search criteria.
type SearchQueryParams struct {
	Query    string
	Category string
	MaxPrice float64
	Limit    int
}

// SearchResponse holds Algolia-style search results, facets, and telemetry.
type SearchResponse struct {
	Hits             []Hit          `json:"hits"`
	Facets           map[string]int `json:"facets"`
	ProcessingTimeMs float64        `json:"processing_time_ms"`
	TotalHits        int            `json:"total_hits"`
}
