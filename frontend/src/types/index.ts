export interface HighlightResult {
  title: string;
  description: string;
}

export interface Hit {
  id: string;
  title: string;
  description: string;
  category: string;
  tags: string[];
  price: number;
  popularity_score: number;
  score: number;
  _highlight: HighlightResult;
}

export interface SearchResponse {
  hits: Hit[];
  facets: Record<string, number>;
  processing_time_ms: number;
  total_hits: number;
}

export interface Document {
  id: string;
  title: string;
  description: string;
  category: string;
  tags: string[];
  price: number;
  popularity_score: number;
  created_at: string;
}

export interface CreateDocumentInput {
  title: string;
  description: string;
  category: string;
  tags: string[];
  price: number;
  popularity_score: number;
}
