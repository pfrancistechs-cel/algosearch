import React from 'react';
import { Zap, CheckCircle2 } from 'lucide-react';

interface TelemetryPillProps {
  totalHits: number;
  processingTimeMs: number;
  query?: string;
}

export function TelemetryPill({ totalHits, processingTimeMs, query }: TelemetryPillProps) {
  return (
    <div className="flex flex-wrap items-center justify-between gap-2 py-2 px-1 text-xs text-slate-500 dark:text-slate-400">
      <div className="flex items-center gap-1.5 font-medium">
        <Zap className="w-3.5 h-3.5 text-amber-500 fill-amber-500/20" />
        <span>
          Found <strong className="text-slate-800 dark:text-slate-200">{totalHits}</strong> results in{' '}
          <strong className="text-indigo-600 dark:text-indigo-400">{processingTimeMs}ms</strong>
        </span>
        {query && query.trim() !== '' && (
          <span className="hidden sm:inline text-slate-400">
            for &quot;<span className="text-slate-600 dark:text-slate-300 font-medium">{query}</span>&quot;
          </span>
        )}
      </div>

      <div className="flex items-center gap-1 text-[11px] font-mono text-slate-400 bg-slate-100 dark:bg-slate-800/80 px-2 py-0.5 rounded-full">
        <CheckCircle2 className="w-3 h-3 text-emerald-500" />
        <span>PostgreSQL pg_trgm + tsvector</span>
      </div>
    </div>
  );
}
