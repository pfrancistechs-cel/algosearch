'use client';

import React from 'react';
import { Filter, RotateCcw, Tag, DollarSign } from 'lucide-react';

interface FacetSidebarProps {
  facets: Record<string, number>;
  selectedCategory: string;
  onSelectCategory: (category: string) => void;
  maxPrice: number;
  onMaxPriceChange: (price: number) => void;
  onReset: () => void;
}

export function FacetSidebar({
  facets,
  selectedCategory,
  onSelectCategory,
  maxPrice,
  onMaxPriceChange,
  onReset,
}: FacetSidebarProps) {
  const categories = Object.entries(facets).sort((a, b) => b[1] - a[1]);
  const hasActiveFilters = selectedCategory !== 'all' || maxPrice < 4000;

  return (
    <aside className="w-full lg:w-64 space-y-6">
      {/* Header */}
      <div className="flex items-center justify-between pb-3 border-b border-slate-200 dark:border-slate-800">
        <div className="flex items-center gap-2 font-semibold text-slate-800 dark:text-slate-100 text-sm">
          <Filter className="w-4 h-4 text-indigo-500" />
          <span>Faceted Filters</span>
        </div>
        {hasActiveFilters && (
          <button
            onClick={onReset}
            className="flex items-center gap-1 text-xs text-indigo-600 dark:text-indigo-400 hover:text-indigo-700 font-medium transition-colors"
          >
            <RotateCcw className="w-3 h-3" />
            <span>Reset</span>
          </button>
        )}
      </div>

      {/* Categories Facet */}
      <div className="space-y-3">
        <div className="flex items-center gap-1.5 text-xs font-semibold uppercase tracking-wider text-slate-400">
          <Tag className="w-3.5 h-3.5" />
          <span>Categories</span>
        </div>

        <div className="space-y-1">
          <button
            onClick={() => onSelectCategory('all')}
            className={`w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm font-medium transition-all ${
              selectedCategory === 'all'
                ? 'bg-indigo-50 dark:bg-indigo-950/50 text-indigo-600 dark:text-indigo-400 font-semibold'
                : 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800/60'
            }`}
          >
            <span>All Categories</span>
          </button>

          {categories.map(([category, count]) => {
            const isSelected = selectedCategory === category;
            return (
              <button
                key={category}
                onClick={() => onSelectCategory(isSelected ? 'all' : category)}
                className={`w-full flex items-center justify-between px-3 py-2 rounded-xl text-sm transition-all ${
                  isSelected
                    ? 'bg-indigo-50 dark:bg-indigo-950/50 text-indigo-600 dark:text-indigo-400 font-semibold'
                    : 'text-slate-600 dark:text-slate-300 hover:bg-slate-100 dark:hover:bg-slate-800/60'
                }`}
              >
                <span className="truncate">{category}</span>
                <span
                  className={`text-xs px-2 py-0.5 rounded-full font-mono ${
                    isSelected
                      ? 'bg-indigo-200/60 dark:bg-indigo-900/60 text-indigo-700 dark:text-indigo-300'
                      : 'bg-slate-100 dark:bg-slate-800 text-slate-500 dark:text-slate-400'
                  }`}
                >
                  {count}
                </span>
              </button>
            );
          })}
        </div>
      </div>

      {/* Max Price Slider */}
      <div className="space-y-3 pt-4 border-t border-slate-200 dark:border-slate-800">
        <div className="flex items-center justify-between text-xs font-semibold uppercase tracking-wider text-slate-400">
          <div className="flex items-center gap-1.5">
            <DollarSign className="w-3.5 h-3.5" />
            <span>Max Price</span>
          </div>
          <span className="text-slate-800 dark:text-slate-200 font-mono text-sm">
            {maxPrice >= 4000 ? 'Any' : `$${maxPrice}`}
          </span>
        </div>

        <input
          type="range"
          min="50"
          max="4000"
          step="50"
          value={maxPrice}
          onChange={(e) => onMaxPriceChange(Number(e.target.value))}
          className="w-full accent-indigo-600 cursor-pointer"
        />

        <div className="flex justify-between text-[11px] text-slate-400 font-mono">
          <span>$50</span>
          <span>$2,000</span>
          <span>$4,000+</span>
        </div>
      </div>
    </aside>
  );
}
