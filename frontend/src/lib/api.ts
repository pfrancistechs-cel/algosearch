import { CreateDocumentInput, Document, Hit, SearchResponse } from '@/types';

const API_BASE = process.env.NEXT_PUBLIC_API_URL || 'http://localhost:8081/api';

export interface SearchParams {
  q?: string;
  category?: string;
  max_price?: number;
  limit?: number;
}

export async function searchDocuments(params: SearchParams): Promise<SearchResponse> {
  const query = new URLSearchParams();
  if (params.q !== undefined && params.q.trim() !== '') {
    query.set('q', params.q.trim());
  }
  if (params.category && params.category !== 'all') {
    query.set('category', params.category);
  }
  if (params.max_price && params.max_price > 0) {
    query.set('max_price', params.max_price.toString());
  }
  if (params.limit) {
    query.set('limit', params.limit.toString());
  }

  const url = `${API_BASE}/search?${query.toString()}`;
  const res = await fetch(url, {
    method: 'GET',
    headers: {
      'Content-Type': 'application/json',
    },
    cache: 'no-store',
  });

  if (!res.ok) {
    const errText = await res.text();
    throw new Error(`Search failed (${res.status}): ${errText}`);
  }

  return res.json();
}

export async function indexDocument(input: CreateDocumentInput): Promise<Hit> {
  const res = await fetch(`${API_BASE}/documents`, {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify(input),
  });

  if (!res.ok) {
    const errText = await res.text();
    throw new Error(`Failed to index document (${res.status}): ${errText}`);
  }

  return res.json();
}

export async function deleteDocument(id: string): Promise<boolean> {
  const res = await fetch(`${API_BASE}/documents/${id}`, {
    method: 'DELETE',
    headers: {
      'Content-Type': 'application/json',
    },
  });

  if (!res.ok) {
    const errText = await res.text();
    throw new Error(`Failed to delete document (${res.status}): ${errText}`);
  }

  return true;
}
