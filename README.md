# AlgoSearch 🔍 ⚡

> **A blazing-fast, typo-tolerant search and indexing engine inspired by Algolia, built with Go & Next.js.**

AlgoSearch provides sub-millisecond search-as-you-type querying, inverted index posting lists, prefix matching with finite state automata/tries, and Levenshtein distance typo tolerance.

---

## 🏗️ Architecture

```text
                     ┌───────────────────────────┐
                     │ AlgoSearch Instant UI     │
                     │ (Next.js Search-As-You-Type)│
                     └─────────────┬─────────────┘
                                   │ HTTP/REST & WebSocket
                                   ▼
┌────────────────────────────────────────────────────────────────────────┐
│                        AlgoSearch Core (Go)                            │
│                                                                        │
│  ┌───────────────────────┐   ┌──────────────────────────────────────┐  │
│  │   Search API Router   │   │         Query Parser & Ranking       │  │
│  │  (Instant Search API) │──▶│   (BM25, Levenshtein, Prefix Trie)   │  │
│  └───────────────────────┘   └──────────────────┬───────────────────┘  │
│                                                 │                      │
│                                                 ▼                      │
│  ┌───────────────────────┐   ┌──────────────────────────────────────┐  │
│  │    Document Store     │   │            Inverted Index            │  │
│  │ (JSON Document Store) │◀──│  (Postings Lists & Token Shards)     │  │
│  └───────────────────────┘   └──────────────────────────────────────┘  │
└────────────────────────────────────────────────────────────────────────┘
```

---

## 📂 Project Structure

```text
algosearch/
├── backend/                  # Go search engine core
│   ├── cmd/
│   │   ├── server/           # Search API HTTP daemon
│   │   └── cli/              # Batch indexer & benchmark CLI
│   ├── pkg/
│   │   ├── api/              # Search & indexing endpoints
│   │   ├── index/            # Inverted index, tokenizer, trie
│   │   ├── query/            # Fuzzy matching, prefix search, ranking
│   │   └── storage/          # Document persistence & write-ahead log
│   └── go.mod
│
├── frontend/                 # Next.js instant search dashboard
│   ├── public/
│   └── src/
│       ├── app/
│       ├── components/
│       ├── hooks/
│       ├── lib/
│       └── types/
│
├── .gitignore
└── README.md
```

---

## 📜 License
MIT © pfrancistechs-cel
