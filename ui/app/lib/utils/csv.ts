/** Shared CSV serialization helpers. */

export function escapeCsv(value: string): string {
  return /[",\n]/.test(value) ? `"${value.replace(/"/g, '""')}"` : value;
}

export function csvCell(value: number | string | undefined | null): string {
  if (value === undefined || value === null) return "";
  return escapeCsv(String(value));
}

/** Join a header plus rows into a CSV document. */
export function toCsv(
  header: string[],
  rows: (number | string | undefined | null)[][]
): string {
  return [
    header.map(csvCell).join(","),
    ...rows.map((row) => row.map(csvCell).join(",")),
  ].join("\n");
}

/** Trigger a client-side download of a CSV string as a file. */
export function downloadCsv(filename: string, csv: string): void {
  if (typeof window === "undefined") return;
  const blob = new Blob([csv], { type: "text/csv;charset=utf-8;" });
  const url = URL.createObjectURL(blob);
  const a = document.createElement("a");
  a.href = url;
  a.download = filename;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}
