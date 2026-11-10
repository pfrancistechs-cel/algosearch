'use client';

import React, { useEffect, useRef } from 'react';
import { Search, X, Loader2 } from 'lucide-react';

interface SearchBarProps {
  query: string;
  onChange: (value: string) => void;
  isLoading?: boolean;
}

export function SearchBar({ query, onChange, isLoading = false }: SearchBarProps) {
  const inputRef = useRef<HTMLInputElement>(null);

  // Global Cmd+K / Ctrl+K shortcut
  useEffect(() => {
    const handleKeyDown = (e: KeyboardEvent) => {
      if ((e.metaKey || e.ctrlKey) && e.key.toLowerCase() === 'k') {
        e.preventDefault();
        inputRef.current?.focus();
        inputRef.current?.select();
      }
      if (e.key === 'Escape' && document.activeElement === inputRef.current) {
        inputRef.current?.blur();
      }
    };

    window.addEventListener('keydown', handleKeyDown);
    return () => window.removeEventListener('keydown', handleKeyDown);
  }, []);

  return (
    <div className="relative w-full group">
      <div className="absolute inset-y-0 left-0 pl-4 flex items-center pointer-events-none text-slate-400 group-focus-within:text-indigo-500 transition-colors">
        {isLoading ? (
          <Loader2 className="w-5 h-5 animate-spin text-indigo-500" />
        ) : (
          <Search className="w-5 h-5" />
        )}
      </div>

      <input
        ref={inputRef}
        type="text"
        value={query}
        onChange={(e) => onChange(e.target.value)}
        placeholder="Search products, models, specs, or typos (e.g. 'macbok', 'iphne', 'anc')..."
        className="w-full pl-12 pr-28 py-3.5 bg-white dark:bg-slate-900 border border-slate-200 dark:border-slate-800 rounded-2xl shadow-sm text-base text-slate-900 dark:text-slate-100 placeholder:text-slate-400 focus:outline-none focus:ring-2 focus:ring-indigo-500/20 focus:border-indigo-500 transition-all"
        autoFocus
        spellCheck={false}
      />

      <div className="absolute inset-y-0 right-0 pr-3 flex items-center gap-1.5">
        {query && (
          <button
            onClick={() => {
              onChange('');
              inputRef.current?.focus();
            }}
            className="p-1.5 text-slate-400 hover:text-slate-600 dark:hover:text-slate-200 hover:bg-slate-100 dark:hover:bg-slate-800 rounded-lg transition-colors"
            title="Clear search"
          >
            <X className="w-4 h-4" />
          </button>
        )}

        <kbd className="hidden sm:inline-flex items-center gap-0.5 px-2 py-1 text-xs font-mono font-medium text-slate-500 dark:text-slate-400 bg-slate-100 dark:bg-slate-800 border border-slate-200 dark:border-slate-700 rounded-md shadow-xs select-none pointer-events-none">
          <span className="text-xs">⌘</span>K
        </kbd>
      </div>
    </div>
  );
}
