# SearchLite (AlgoSearch) 🔍 ⚡

> **A blazing-fast, typo-tolerant Algolia-style instant search engine built with Go, PostgreSQL (`tsvector` + `pg_trgm`), and Next.js 15 (App Router).**

Zero external search clusters (Elasticsearch/Meilisearch). Zero Redis. Production-ready sub-millisecond search-as-you-type using PostgreSQL's native weighted full-text search and trigram similarity indexing.

---

## 🏗️ Architecture

```text
                     ┌─────────────────────────────────────────┐
                     │          SearchLite Instant UI          │
                     │  (Next.js 15, Debounced 150ms, ⌘K Input) │
                     └────────────────────┬────────────────────┘
                                          │ HTTP GET /api/search
                                          ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                             Go Search API Backend                                │
│                                                                                  │
│   ┌─────────────────────┐   ┌───────────────────────┐   ┌────────────────────┐   │
│   │   HTTP Router       │──▶│   Search Service      │──▶│   Telemetry Timer  │   │
│   │   (Go 1.22+ Mux)    │   │   (Hybrid Scoring)    │   │   (Sub-ms Tracking)│   │
│   └─────────────────────┘   └───────────┬───────────┘   └────────────────────┘   │
│                                         │                                        │
│                                         ▼                                        │
│                         pgx/v5 High-Performance Pool                             │
└─────────────────────────────────────────┼────────────────────────────────────────┘
                                          │
                                          ▼
┌──────────────────────────────────────────────────────────────────────────────────┐
│                         PostgreSQL (Neon / Local)                                │
│                                                                                  │
│   ┌────────────────────────────────┐       ┌─────────────────────────────────┐   │
│   │   search_vector (GIN Index)    │       │     title (pg_trgm GIN Index)   │   │
│   │   Weighted TSVECTOR [A/B/C]    │       │     Typo Tolerance & Similarity │   │
│   └────────────────────────────────┘       └─────────────────────────────────┘   │
│   ┌────────────────────────────────┐       ┌─────────────────────────────────┐   │
│   │   ts_headline Highlighting     │       │     Dynamic Facet Counts        │   │
│   │   <mark> Snippets Generation   │       │     Grouped Category Breakdown  │   │
│   └────────────────────────────────┘       └─────────────────────────────────┘   │
└──────────────────────────────────────────────────────────────────────────────────┘
```

---

## 🚀 Key Features

- **⚡ Sub-Millisecond Instant Search:** Live search-as-you-type with 150ms debounce and exact `processing_time_ms` telemetry tracking.
- **🎯 Typo Tolerance & Fuzzy Matching:** Powered by PostgreSQL's `pg_trgm` extension. Querying misspelled terms like `"iphne"`, `"macbok"`, or `"cancelling"` matches exact products accurately.
- **🏆 Hybrid Relevance Ranking:** Weighted formula combining full-text `ts_rank` on weighted attributes (`title` [A], `category` [B], `description` [C]), trigram similarity (`similarity(title, q)`), and normalized popularity scores.
- **🏷️ Real-Time Dynamic Facets:** Dynamic category count breakdown computed per query alongside active max-price range filtering.
- **✨ In-Line Match Highlighting:** Matched terms wrapped with `<mark>` tags via `ts_headline` and highlighted cleanly in the UI.
- **➕ Instant In-Browser Record Indexing:** "+ Index Record" modal indexing new documents via `POST /api/documents` into PostgreSQL in real time.
- **⌨️ Keyboard Shortcuts:** Focus search instantly from anywhere using `⌘K` or `Ctrl+K`.

---

## 📂 Project Structure

```text
algosearch/
├── backend/
│   ├── cmd/
│   │   └── api/
│   │       └── main.go           # Entrypoint, CORS, graceful shutdown
│   ├── internal/
│   │   ├── db/
│   │   │   └── db.go             # pgxpool connection, auto-migrations & 30 seed records
│   │   └── search/
│   │       ├── model.go          # Domain models, Hit, SearchResponse
│   │       ├── repository.go     # Hybrid SQL queries, ts_headline, facets
│   │       ├── service.go        # Telemetry timer & business logic
│   │       └── handler.go        # HTTP endpoints (/api/search, /api/documents)
│   ├── .env.example
│   └── go.mod
│
├── frontend/
│   ├── src/
│   │   ├── app/
│   │   │   ├── layout.tsx
│   │   │   └── page.tsx          # Main SearchLite instant search experience
│   │   ├── components/
│   │   │   ├── SearchBar.tsx     # ⌘K Command bar with clear button
│   │   │   ├── TelemetryPill.tsx # "Found X results in Y ms" badge
│   │   │   ├── FacetSidebar.tsx  # Dynamic category counts & price slider
│   │   │   ├── ResultCard.tsx    # Cards with price, popularity & delete
│   │   │   ├── HighlightedText.tsx # Styled <mark> tag renderer
│   │   │   └── AddDocumentModal.tsx # Live index record modal
│   │   ├── lib/
│   │   │   └── api.ts            # Client API fetcher
│   │   └── types/
│   │       └── index.ts          # TypeScript interfaces
│   ├── .env.example
│   └── package.json
│
├── docker-compose.yml            # Local PostgreSQL 16 container with pg_trgm
├── .gitignore
└── README.md
```

---

## 🛠️ Quickstart & Setup

### 1. Database Setup (Choose Neon or Local Docker)

#### Option A: Cloud Database (Neon PostgreSQL)
Create a PostgreSQL database on [Neon.tech](https://neon.tech) and copy your connection string (`sslmode=require`).

#### Option B: Local Docker PostgreSQL
```bash
docker-compose up -d
```
*(Runs PostgreSQL 16 on `localhost:5433`)*

---

### 2. Backend Setup (`/backend`)

1. Copy `.env.example` to `.env`:
   ```bash
   cd backend
   cp .env.example .env
   ```
2. Configure `DATABASE_URL` in `.env`:
   ```ini
   PORT=8081
   # For Neon:
   DATABASE_URL=postgres://neondb_owner:YOUR_PASSWORD@ep-sample-12345.us-east-2.aws.neon.tech/neondb?sslmode=require
   # For Local Docker:
   # DATABASE_URL=postgres://algosearch:algosearch_password@localhost:5433/algosearch?sslmode=disable
   ```
3. Run the Go API server:
   ```bash
   go run ./cmd/api
   ```
   > **Note:** On first startup, the database automatically executes migrations (`pg_trgm`, `search_vector`, GIN indexes) and seeds **30 realistic tech products** if the collection is empty.

---

### 3. Frontend Setup (`/frontend`)

1. Copy `.env.example` to `.env.local`:
   ```bash
   cd ../frontend
   cp .env.example .env.local
   ```
2. Start the Next.js development server:
   ```bash
   npm run dev
   ```
3. Open **[http://localhost:3000](http://localhost:3000)** in your browser.

---

## 📡 REST API Reference

### `GET /api/search`
Executes hybrid full-text and trigram fuzzy matching with facets and telemetry.

**Query Parameters:**
- `q`: Search string (e.g. `iphne`, `macbok`, `anc`)
- `category`: Category filter (e.g. `Laptops`, `Audio`)
- `max_price`: Maximum price threshold
- `limit`: Number of results (default `20`)

**Response:**
```json
{
  "hits": [
    {
      "id": "e93ab720-91a5-48fa-bb64-c71501c51888",
      "title": "Apple MacBook Pro 16-inch M3 Max",
      "description": "Liquid Retina XDR display with 128GB unified memory...",
      "category": "Laptops",
      "tags": ["apple", "laptop", "m3", "pro"],
      "price": 3499.00,
      "popularity_score": 98,
      "score": 0.94,
      "_highlight": {
        "title": "Apple <mark>MacBook</mark> Pro 16-inch M3 Max",
        "description": "Liquid Retina XDR display with 128GB unified memory..."
      }
    }
  ],
  "facets": {
    "Laptops": 5,
    "Smartphones": 3,
    "Audio": 4
  },
  "processing_time_ms": 1.15,
  "total_hits": 12
}
```

### `POST /api/documents`
Indexes a new document in real time into PostgreSQL.

**Payload:**
```json
{
  "title": "Sony WH-1000XM6 Wireless Headphones",
  "description": "Next-generation noise canceling headphones with AI sound calibration.",
  "category": "Audio",
  "tags": ["sony", "anc", "hifi"],
  "price": 449.00,
  "popularity_score": 92
}
```

### `DELETE /api/documents/{id}`
Deletes a document from the search index.

---

## 📜 License
MIT © pfrancistechs-cel
