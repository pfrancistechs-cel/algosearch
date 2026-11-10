import type { Metadata } from "next";
import "./globals.css";

export const metadata: Metadata = {
  title: "AlgoSearch",
  description: "Instant typo-tolerant search and indexing engine",
};

export default function RootLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <html lang="en">
      <body>{children}</body>
    </html>
  );
}

