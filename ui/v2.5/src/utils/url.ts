/**
 * Extracts the sortable portion of a URL by removing the protocol and www. prefix
 */
const urlSortKey = (url: string): string =>
  url.toLowerCase().replace(/^(https?:\/\/)?(www\.)?/, "");

/**
 * Sorts a list of URLs alphabetically by their base URL,
 * excluding the protocol (http/https) and www. prefix.
 * Returns a new sorted array without mutating the original.
 */
export function sortURLs(urls: string[]): string[] {
  return [...urls].sort((a, b) => urlSortKey(a).localeCompare(urlSortKey(b)));
}
