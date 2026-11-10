package search

import (
	"context"
	"errors"
	"math"
	"strings"
	"time"
)

// Service coordinates search operations and calculates execution telemetry.
type Service struct {
	repo *Repository
}

// NewService creates a new search service.
func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// Search processes the query, records elapsed time in milliseconds, and returns the response.
func (s *Service) Search(ctx context.Context, params SearchQueryParams) (*SearchResponse, error) {
	start := time.Now()

	hits, facets, totalHits, err := s.repo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	elapsedMs := float64(time.Since(start).Microseconds()) / 1000.0
	// Round to two decimal places
	elapsedMs = math.Round(elapsedMs*100) / 100

	if hits == nil {
		hits = []Hit{}
	}
	if facets == nil {
		facets = make(map[string]int)
	}

	return &SearchResponse{
		Hits:             hits,
		Facets:           facets,
		ProcessingTimeMs: elapsedMs,
		TotalHits:        totalHits,
	}, nil
}

// IndexDocument validates and creates a new document.
func (s *Service) IndexDocument(ctx context.Context, input CreateDocumentInput) (*Document, error) {
	if strings.TrimSpace(input.Title) == "" {
		return nil, errors.New("title is required")
	}
	if strings.TrimSpace(input.Description) == "" {
		return nil, errors.New("description is required")
	}
	if strings.TrimSpace(input.Category) == "" {
		return nil, errors.New("category is required")
	}
	if input.Tags == nil {
		input.Tags = []string{}
	}
	if input.Price < 0 {
		input.Price = 0
	}
	if input.PopularityScore < 0 {
		input.PopularityScore = 50
	}

	return s.repo.CreateDocument(ctx, input)
}

// DeleteDocument deletes an indexed document.
func (s *Service) DeleteDocument(ctx context.Context, id string) (bool, error) {
	id = strings.TrimSpace(id)
	if id == "" {
		return false, errors.New("document id is required")
	}
	return s.repo.DeleteDocument(ctx, id)
}
