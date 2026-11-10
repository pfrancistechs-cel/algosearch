'use client';

import React, { useState, useEffect, useCallback, useTransition } from 'react';
import { SearchResponse, Hit } from '@/types';
import { searchDocuments, deleteDocument } from '@/lib/api';
import { SearchBar } from '@/components/SearchBar';
import { TelemetryPill } from '@/components/TelemetryPill';
import { FacetSidebar } from '@/components/FacetSidebar';
import { ResultCard } from '@/components/ResultCard';
import { AddDocumentModal } from '@/components/AddDocumentModal';
import { Search, Plus, Compass, AlertCircle, Sparkles } from 'lucide-react';

export default function SearchPage() {
  const [query, setQuery] = useState('');
  const [debouncedQuery, setDebouncedQuery] = useState('');
  const [selectedCategory, setSelectedCategory] = useState('all');
  const [maxPrice, setMaxPrice] = useState(4000);
  const [data, setData] = useState<SearchResponse>({
    hits: [],
    facets: {},
    processing_time_ms: 0,
    total_hits: 0,
  });
  const [isLoading, setIsLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [, startTransition] = useTransition();

  // 150ms debounce on input query
  useEffect(() => {
    const timer = setTimeout(() => {
      setDebouncedQuery(query);
    }, 150);
    return () => clearTimeout(timer);
  }, [query]);

  // Execute search
  const performSearch = useCallback(async () => {
    setIsLoading(true);
    setError(null);
    try {
      const res = await searchDocuments({
        q: debouncedQuery,
        category: selectedCategory === 'all' ? undefined : selectedCategory,
        max_price: maxPrice < 4000 ? maxPrice : undefined,
        limit: 30,
      });
      startTransition(() => {
        setData(res);
      });
    } catch (err: any) {
      setError(err?.message || 'Error executing search');
    } finally {
      setIsLoading(false);
    }
  }, [debouncedQuery, selectedCategory, maxPrice]);

  useEffect(() => {
    performSearch();
  }, [performSearch]);

  const handleDelete = async (id: string) => {
    try {
      await deleteDocument(id);
      setData((prev) => ({
        ...prev,
        hits: prev.hits.filter((h) => h.id !== id),
        total_hits: Math.max(0, prev.total_hits - 1),
      }));
    } catch (err: any) {
      alert(err?.message || 'Failed to delete record');
    }
  };

  const handleResetFilters = () => {
    setSelectedCategory('all');
    setMaxPrice(4000);
  };

  return (
    <main className="min-h-screen bg-slate-50 dark:bg-slate-950 text-slate-900 dark:text-slate-100 flex flex-col">
      {/* Top Navigation */}
      <header className="sticky top-0 z-30 bg-white/80 dark:bg-slate-900/80 backdrop-blur-md border-b border-slate-200 dark:border-slate-800">
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8 h-16 flex items-center justify-between gap-4">
          <div className="flex items-center gap-3">
            <div className="w-9 h-9 rounded-xl bg-gradient-to-tr from-indigo-600 to-violet-500 flex items-center justify-center text-white shadow-xs">
              <Search className="w-5 h-5" />
            </div>
            <div>
              <div className="flex items-center gap-2">
                <span className="font-bold tracking-tight text-slate-900 dark:text-white text-lg">
                  SearchLite
                </span>
                <span className="text-[10px] uppercase font-mono font-bold tracking-widest px-2 py-0.5 rounded-md bg-indigo-50 dark:bg-indigo-950/60 text-indigo-600 dark:text-indigo-400 border border-indigo-200 dark:border-indigo-800">
                  Algolia Fast
                </span>
              </div>
            </div>
          </div>

          <div className="flex items-center gap-3">
            <button
              onClick={() => setIsModalOpen(true)}
              className="flex items-center gap-1.5 px-3.5 py-2 bg-indigo-600 hover:bg-indigo-700 text-white text-xs sm:text-sm font-medium rounded-xl shadow-xs transition-colors"
            >
              <Plus className="w-4 h-4" />
              <span>Index Record</span>
            </button>
          </div>
        </div>
      </header>

      {/* Hero Search Section */}
      <section className="bg-white dark:bg-slate-900/50 border-b border-slate-200 dark:border-slate-800/80 py-8 px-4 sm:px-6 lg:px-8 shadow-xs">
        <div className="max-w-4xl mx-auto space-y-3">
          <SearchBar query={query} onChange={setQuery} isLoading={isLoading} />
          <TelemetryPill
            totalHits={data.total_hits}
            processingTimeMs={data.processing_time_ms}
            query={debouncedQuery}
          />
        </div>
      </section>

      {/* Main Content Area: Sidebar + Results Grid */}
      <div className="flex-1 max-w-7xl w-full mx-auto px-4 sm:px-6 lg:px-8 py-8">
        <div className="flex flex-col lg:flex-row gap-8 items-start">
          {/* Facet Sidebar */}
          <FacetSidebar
            facets={data.facets}
            selectedCategory={selectedCategory}
            onSelectCategory={setSelectedCategory}
            maxPrice={maxPrice}
            onMaxPriceChange={setMaxPrice}
            onReset={handleResetFilters}
          />

          {/* Results Area */}
          <section className="flex-1 w-full space-y-6">
            {error && (
              <div className="flex items-center gap-3 p-4 bg-rose-50 dark:bg-rose-950/30 border border-rose-200 dark:border-rose-900 rounded-2xl text-rose-700 dark:text-rose-300 text-sm">
                <AlertCircle className="w-5 h-5 shrink-0" />
                <p>Could not load search results: {error}</p>
              </div>
            )}

            {/* Results Grid */}
            {data.hits.length > 0 ? (
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {data.hits.map((hit) => (
                  <ResultCard key={hit.id} hit={hit} onDelete={handleDelete} />
                ))}
              </div>
            ) : !isLoading ? (
              <div className="text-center py-16 px-4 bg-white dark:bg-slate-900 border border-dashed border-slate-300 dark:border-slate-800 rounded-3xl space-y-4">
                <div className="w-12 h-12 mx-auto rounded-2xl bg-indigo-50 dark:bg-indigo-950/50 flex items-center justify-center text-indigo-500">
                  <Compass className="w-6 h-6" />
                </div>
                <div>
                  <h3 className="text-base font-semibold text-slate-900 dark:text-slate-100">
                    No results found
                  </h3>
                  <p className="text-sm text-slate-500 max-w-md mx-auto mt-1">
                    No matches found for &quot;{query}&quot; in{' '}
                    {selectedCategory === 'all' ? 'any category' : selectedCategory}. Try
                    searching for &quot;macbook&quot;, &quot;iphone&quot;, or testing typos like
                    &quot;iphne&quot;.
                  </p>
                </div>
                <button
                  onClick={() => {
                    setQuery('');
                    handleResetFilters();
                  }}
                  className="inline-flex items-center gap-2 px-4 py-2 bg-slate-100 dark:bg-slate-800 hover:bg-slate-200 dark:hover:bg-slate-700 text-slate-700 dark:text-slate-200 text-xs font-semibold rounded-xl transition-colors"
                >
                  <Sparkles className="w-3.5 h-3.5 text-indigo-500" />
                  <span>Reset All Filters</span>
                </button>
              </div>
            ) : (
              // Skeleton loading state
              <div className="grid grid-cols-1 md:grid-cols-2 xl:grid-cols-3 gap-4">
                {[1, 2, 3, 4, 5, 6].map((i) => (
                  <div
                    key={i}
                    className="h-48 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-5 animate-pulse space-y-3"
                  >
                    <div className="flex justify-between">
                      <div className="h-4 w-20 bg-slate-200 dark:bg-slate-800 rounded" />
                      <div className="h-4 w-12 bg-slate-200 dark:bg-slate-800 rounded" />
                    </div>
                    <div className="h-5 w-3/4 bg-slate-200 dark:bg-slate-800 rounded" />
                    <div className="h-4 w-full bg-slate-200 dark:bg-slate-800 rounded" />
                    <div className="h-4 w-2/3 bg-slate-200 dark:bg-slate-800 rounded" />
                  </div>
                ))}
              </div>
            )}
          </section>
        </div>
      </div>

      {/* Index Document Modal */}
      <AddDocumentModal
        isOpen={isModalOpen}
        onClose={() => setIsModalOpen(false)}
        onIndexed={performSearch}
      />
    </main>
  );
}
