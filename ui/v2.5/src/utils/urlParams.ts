/**
 * Utility functions for handling URL parameters consistently
 */

/**
 * Adds or updates a returnTo parameter to a URL with the current location
 * @param pathname The base pathname to navigate to
 * @param returnLocation The location to return to (usually current location)
 * @returns Object with pathname and search params
 */
export function withReturnTo(pathname: string, returnLocation: string) {
  return {
    pathname,
    search: `?returnTo=${encodeURIComponent(returnLocation)}`,
  };
}

/**
 * Gets the returnTo parameter from the search string
 * @param search The search string to parse
 * @returns The decoded returnTo URL or null if not found
 */
export function getReturnTo(search: string): string | null {
  const params = new URLSearchParams(search);
  const returnTo = params.get("returnTo");
  return returnTo ? decodeURIComponent(returnTo) : null;
}

/**
 * Preserves non-filter query parameters when updating filter parameters
 * @param newParams The new filter parameters
 * @param currentSearch The current search string
 * @returns The combined search string
 */
export function preserveNonFilterParams(
  newParams: string,
  currentSearch: string
): string {
  const currentParams = new URLSearchParams(currentSearch);
  const newSearchParams = new URLSearchParams(newParams);

  // List of parameters that should be preserved when changing filters
  const preserveParams = ["returnTo"];

  for (const param of preserveParams) {
    const value = currentParams.get(param);
    if (value !== null) {
      newSearchParams.set(param, value);
    }
  }

  return newSearchParams.toString();
}

/**
 * Determines if there are actual filter parameters in the URL
 * @param search The current search string
 * @returns Whether there are any filter parameters
 */
export function hasFilterParams(search: string): boolean {
  const params = new URLSearchParams(search);
  return Array.from(params.keys()).some((key) => key !== "returnTo");
}
