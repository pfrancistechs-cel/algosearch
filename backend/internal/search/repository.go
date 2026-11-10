package search

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Repository handles database queries for SearchLite.
type Repository struct {
	pool *pgxpool.Pool
}

// NewRepository creates a new search repository.
func NewRepository(pool *pgxpool.Pool) *Repository {
	return &Repository{pool: pool}
}

// Search executes hybrid full-text and trigram fuzzy matching with facets.
func (r *Repository) Search(ctx context.Context, params SearchQueryParams) ([]Hit, map[string]int, int, error) {
	cleanQuery := strings.TrimSpace(params.Query)
	limit := params.Limit
	if limit <= 0 || limit > 100 {
		limit = 20
	}

	if cleanQuery == "" {
		// Empty query: Browse mode ordered by popularity
		return r.searchEmptyQuery(ctx, params, limit)
	}

	return r.searchHybridQuery(ctx, cleanQuery, params, limit)
}

func (r *Repository) searchEmptyQuery(ctx context.Context, params SearchQueryParams, limit int) ([]Hit, map[string]int, int, error) {
	hits := make([]Hit, 0)
	facets := make(map[string]int)

	var whereClauses []string
	var args []interface{}
	argIdx := 1

	if params.Category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, params.Category)
		argIdx++
	}
	if params.MaxPrice > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argIdx))
		args = append(args, params.MaxPrice)
		argIdx++
	}

	whereSQL := ""
	if len(whereClauses) > 0 {
		whereSQL = "WHERE " + strings.Join(whereClauses, " AND ")
	}

	// 1. Total Hits
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM documents %s", whereSQL)
	var totalHits int
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&totalHits); err != nil {
		return nil, nil, 0, fmt.Errorf("count query failed: %w", err)
	}

	// 2. Fetch Hits
	selectSQL := fmt.Sprintf(`
		SELECT id, title, description, category, tags, price, popularity_score
		FROM documents
		%s
		ORDER BY popularity_score DESC, created_at DESC
		LIMIT $%d
	`, whereSQL, argIdx)
	argsWithLimit := append(args, limit)

	rows, err := r.pool.Query(ctx, selectSQL, argsWithLimit...)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("search query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h Hit
		var tags []string
		if err := rows.Scan(&h.ID, &h.Title, &h.Description, &h.Category, &tags, &h.Price, &h.PopularityScore); err != nil {
			return nil, nil, 0, fmt.Errorf("scan error: %w", err)
		}
		h.Tags = tags
		h.Score = float64(h.PopularityScore) / 100.0
		h.Highlight = HighlightResult{
			Title:       h.Title,
			Description: h.Description,
		}
		hits = append(hits, h)
	}

	// 3. Facets grouped by category (without category filter to show all options)
	var facetWhereClauses []string
	var facetArgs []interface{}
	fIdx := 1
	if params.MaxPrice > 0 {
		facetWhereClauses = append(facetWhereClauses, fmt.Sprintf("price <= $%d", fIdx))
		facetArgs = append(facetArgs, params.MaxPrice)
		fIdx++
	}
	facetWhereSQL := ""
	if len(facetWhereClauses) > 0 {
		facetWhereSQL = "WHERE " + strings.Join(facetWhereClauses, " AND ")
	}

	facetSQL := fmt.Sprintf(`
		SELECT category, COUNT(*) 
		FROM documents 
		%s 
		GROUP BY category 
		ORDER BY COUNT(*) DESC
	`, facetWhereSQL)

	facetRows, err := r.pool.Query(ctx, facetSQL, facetArgs...)
	if err == nil {
		defer facetRows.Close()
		for facetRows.Next() {
			var cat string
			var count int
			if err := facetRows.Scan(&cat, &count); err == nil {
				facets[cat] = count
			}
		}
	}

	return hits, facets, totalHits, nil
}

func (r *Repository) searchHybridQuery(ctx context.Context, cleanQuery string, params SearchQueryParams, limit int) ([]Hit, map[string]int, int, error) {
	hits := make([]Hit, 0)
	facets := make(map[string]int)

	// Combine Full-Text ts_rank and Trigram similarity for typo tolerance ("iphne" -> "iPhone")
	baseFilter := `(
		search_vector @@ websearch_to_tsquery('english', $1)
		OR similarity(title, $1) > 0.18
		OR title ILIKE '%' || $1 || '%'
		OR description ILIKE '%' || $1 || '%'
	)`

	var whereClauses []string
	whereClauses = append(whereClauses, baseFilter)
	args := []interface{}{cleanQuery}
	argIdx := 2

	if params.Category != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("category = $%d", argIdx))
		args = append(args, params.Category)
		argIdx++
	}
	if params.MaxPrice > 0 {
		whereClauses = append(whereClauses, fmt.Sprintf("price <= $%d", argIdx))
		args = append(args, params.MaxPrice)
		argIdx++
	}

	whereSQL := "WHERE " + strings.Join(whereClauses, " AND ")

	// 1. Total Hits Count
	countSQL := fmt.Sprintf("SELECT COUNT(*) FROM documents %s", whereSQL)
	var totalHits int
	if err := r.pool.QueryRow(ctx, countSQL, args...).Scan(&totalHits); err != nil {
		return nil, nil, 0, fmt.Errorf("hybrid count failed: %w", err)
	}

	// 2. Fetch Ranked Hits with Highlights
	// Ranking formula: 60% ts_rank + 30% trigram similarity + 10% normalized popularity
	selectSQL := fmt.Sprintf(`
		SELECT 
			id, 
			title, 
			description, 
			category, 
			tags, 
			price, 
			popularity_score,
			(
				(coalesce(ts_rank(search_vector, websearch_to_tsquery('english', $1)), 0.0) * 0.6) +
				(similarity(title, $1) * 0.3) +
				((popularity_score::float / 100.0) * 0.1)
			) AS rank_score,
			coalesce(
				nullif(ts_headline('english', title, websearch_to_tsquery('english', $1), 'StartSel=<mark>, StopSel=</mark>, HighlightAll=TRUE'), ''),
				title
			) AS hl_title,
			coalesce(
				nullif(ts_headline('english', description, websearch_to_tsquery('english', $1), 'StartSel=<mark>, StopSel=</mark>, MaxWords=30, MinWords=15'), ''),
				description
			) AS hl_desc
		FROM documents
		%s
		ORDER BY rank_score DESC, popularity_score DESC
		LIMIT $%d
	`, whereSQL, argIdx)

	argsWithLimit := append(args, limit)
	rows, err := r.pool.Query(ctx, selectSQL, argsWithLimit...)
	if err != nil {
		return nil, nil, 0, fmt.Errorf("hybrid query failed: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var h Hit
		var tags []string
		var hlTitle, hlDesc string
		if err := rows.Scan(
			&h.ID, 
			&h.Title, 
			&h.Description, 
			&h.Category, 
			&tags, 
			&h.Price, 
			&h.PopularityScore, 
			&h.Score, 
			&hlTitle, 
			&hlDesc,
		); err != nil {
			return nil, nil, 0, fmt.Errorf("scan hybrid row error: %w", err)
		}
		h.Tags = tags

		// Fallback for typo highlight if ts_headline produced no marks
		if !strings.Contains(hlTitle, "<mark>") && strings.EqualFold(h.Title, cleanQuery) {
			hlTitle = fmt.Sprintf("<mark>%s</mark>", h.Title)
		}

		h.Highlight = HighlightResult{
			Title:       hlTitle,
			Description: hlDesc,
		}
		hits = append(hits, h)
	}

	// 3. Dynamic Facets for matched search
	var facetWhereClauses []string
	facetWhereClauses = append(facetWhereClauses, baseFilter)
	facetArgs := []interface{}{cleanQuery}
	fIdx := 2
	if params.MaxPrice > 0 {
		facetWhereClauses = append(facetWhereClauses, fmt.Sprintf("price <= $%d", fIdx))
		facetArgs = append(facetArgs, params.MaxPrice)
		fIdx++
	}
	facetWhereSQL := "WHERE " + strings.Join(facetWhereClauses, " AND ")

	facetSQL := fmt.Sprintf(`
		SELECT category, COUNT(*) 
		FROM documents 
		%s 
		GROUP BY category 
		ORDER BY COUNT(*) DESC
	`, facetWhereSQL)

	facetRows, err := r.pool.Query(ctx, facetSQL, facetArgs...)
	if err == nil {
		defer facetRows.Close()
		for facetRows.Next() {
			var cat string
			var count int
			if err := facetRows.Scan(&cat, &count); err == nil {
				facets[cat] = count
			}
		}
	}

	return hits, facets, totalHits, nil
}

// CreateDocument inserts a new document and lets PostgreSQL generate the weighted search_vector.
func (r *Repository) CreateDocument(ctx context.Context, input CreateDocumentInput) (*Document, error) {
	doc := &Document{
		Title:           input.Title,
		Description:     input.Description,
		Category:        input.Category,
		Tags:            input.Tags,
		Price:           input.Price,
		PopularityScore: input.PopularityScore,
	}

	query := `
		INSERT INTO documents (title, description, category, tags, price, popularity_score)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id, created_at
	`

	err := r.pool.QueryRow(
		ctx,
		query,
		doc.Title,
		doc.Description,
		doc.Category,
		doc.Tags,
		doc.Price,
		doc.PopularityScore,
	).Scan(&doc.ID, &doc.CreatedAt)
	if err != nil {
		return nil, fmt.Errorf("failed to insert document: %w", err)
	}

	return doc, nil
}

// DeleteDocument removes a document by ID.
func (r *Repository) DeleteDocument(ctx context.Context, id string) (bool, error) {
	result, err := r.pool.Exec(ctx, "DELETE FROM documents WHERE id = $1", id)
	if err != nil {
		return false, fmt.Errorf("failed to delete document: %w", err)
	}
	return result.RowsAffected() > 0, nil
}
