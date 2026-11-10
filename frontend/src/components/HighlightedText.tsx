import React from 'react';

interface HighlightedTextProps {
  html: string;
  className?: string;
}

export function HighlightedText({ html, className = '' }: HighlightedTextProps) {
  // If no mark tags exist, render plain text
  if (!html.includes('<mark>')) {
    return <span className={className}>{html}</span>;
  }

  // Parse HTML string with <mark> tags into React elements
  const parts: React.ReactNode[] = [];
  const regex = /<mark>(.*?)<\/mark>/gi;
  let lastIndex = 0;
  let match: RegExpExecArray | null;

  while ((match = regex.exec(html)) !== null) {
    if (match.index > lastIndex) {
      parts.push(html.substring(lastIndex, match.index));
    }
    parts.push(
      <mark
        key={match.index}
        className="bg-amber-200/80 dark:bg-amber-500/30 text-amber-950 dark:text-amber-200 font-semibold px-1 py-0.5 rounded transition-colors"
      >
        {match[1]}
      </mark>
    );
    lastIndex = regex.lastIndex;
  }

  if (lastIndex < html.length) {
    parts.push(html.substring(lastIndex));
  }

  return <span className={className}>{parts}</span>;
}
