"use client";

import "./globals.css";

interface GlobalErrorProps {
  error: Error & { digest?: string };
  reset: () => void;
}

/**
 * Replaces the root layout when the layout itself fails, so it can't use
 * Providers, antd or the header: plain markup on the global styles only.
 */
export default function GlobalError({ error, reset }: GlobalErrorProps) {
  const detail = error.digest ? `digest ${error.digest}` : error.message;

  return (
    <html lang="en" className="dark">
      <body className="min-h-dvh bg-background font-sans text-foreground antialiased">
        <main className="mx-auto max-w-xl px-4 pt-24 md:px-6">
          <h1 className="mb-3 text-2xl font-semibold">XWave failed to load</h1>
          <p className="mb-4 text-base text-neutral-400">
            Something broke before the page could render. Reloading usually
            fixes it.
          </p>
          {detail && (
            <pre className="mb-8 overflow-x-auto rounded border border-border bg-surface px-3 py-2 font-mono text-sm text-neutral-400">
              {detail}
            </pre>
          )}
          <button
            type="button"
            onClick={reset}
            className="rounded bg-primary px-4 py-1.5 text-white hover:brightness-110"
          >
            Reload
          </button>
        </main>
      </body>
    </html>
  );
}
