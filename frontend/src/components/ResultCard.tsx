'use client';

import React, { useState } from 'react';
import { Hit } from '@/types';
import { HighlightedText } from './HighlightedText';
import { Tag, Sparkles, Trash2, Loader2 } from 'lucide-react';

interface ResultCardProps {
  hit: Hit;
  onDelete: (id: string) => Promise<void>;
}

export function ResultCard({ hit, onDelete }: ResultCardProps) {
  const [isDeleting, setIsDeleting] = useState(false);

  const handleDelete = async (e: React.MouseEvent) => {
    e.stopPropagation();
    if (!confirm(`Are you sure you want to remove "${hit.title}" from the search index?`)) {
      return;
    }
    setIsDeleting(true);
    try {
      await onDelete(hit.id);
    } finally {
      setIsDeleting(false);
    }
  };

  return (
    <article className="group relative bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl p-5 hover:border-indigo-400 dark:hover:border-indigo-500/50 hover:shadow-md transition-all flex flex-col justify-between">
      <div>
        {/* Top meta row */}
        <div className="flex items-center justify-between gap-2 mb-2.5">
          <div className="flex items-center gap-2">
            <span className="inline-flex items-center text-[11px] font-semibold uppercase tracking-wider text-indigo-600 dark:text-indigo-400 bg-indigo-50 dark:bg-indigo-950/60 px-2.5 py-0.5 rounded-full">
              {hit.category}
            </span>
            <span className="inline-flex items-center gap-1 text-[11px] font-mono text-slate-400 bg-slate-100 dark:bg-slate-800 px-2 py-0.5 rounded-full">
              <Sparkles className="w-3 h-3 text-amber-500" />
              <span>{hit.popularity_score}% pop</span>
            </span>
          </div>

          <div className="flex items-center gap-2">
            <span className="text-base font-bold text-slate-900 dark:text-white font-mono">
              ${hit.price.toFixed(2)}
            </span>
            <button
              onClick={handleDelete}
              disabled={isDeleting}
              className="opacity-0 group-hover:opacity-100 p-1.5 text-slate-400 hover:text-rose-500 hover:bg-rose-50 dark:hover:bg-rose-950/30 rounded-lg transition-all"
              title="Delete from index"
            >
              {isDeleting ? (
                <Loader2 className="w-4 h-4 animate-spin text-rose-500" />
              ) : (
                <Trash2 className="w-4 h-4" />
              )}
            </button>
          </div>
        </div>

        {/* Title */}
        <h3 className="text-base font-semibold text-slate-900 dark:text-slate-100 group-hover:text-indigo-600 dark:group-hover:text-indigo-400 transition-colors line-clamp-1 mb-2">
          <HighlightedText html={hit._highlight?.title || hit.title} />
        </h3>

        {/* Description snippet */}
        <p className="text-sm text-slate-600 dark:text-slate-300 leading-relaxed line-clamp-2 mb-4">
          <HighlightedText html={hit._highlight?.description || hit.description} />
        </p>
      </div>

      {/* Tags */}
      {hit.tags && hit.tags.length > 0 && (
        <div className="flex flex-wrap items-center gap-1.5 pt-3 border-t border-slate-100 dark:border-slate-800/60">
          <Tag className="w-3 h-3 text-slate-400" />
          {hit.tags.map((tag) => (
            <span
              key={tag}
              className="text-[11px] font-mono text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 px-1.5 py-0.5 rounded"
            >
              #{tag}
            </span>
          ))}
        </div>
      )}
    </article>
  );
}
