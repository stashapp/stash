export interface ISearchablePackage {
  package_id: string;
  name: string;
  metadata?: { [key: string]: unknown } | null;
}

export interface IPackageMatch {
  /** The metadata value that matched, when the name and id did not */
  via?: string;
}

// metadata keys published by package sources such as CommunityScrapers
const SCRAPES_KEY = "scrapes";
const URLS_KEY_SUFFIX = "_urls";

function stringList(value: unknown): string[] {
  if (!Array.isArray(value)) return [];
  return value.filter((v): v is string => typeof v === "string");
}

function matchesURLPattern(pattern: string, query: string) {
  const p = pattern.toLowerCase();
  // a pasted URL contains the pattern, the same way scrapers match URLs
  return p.includes(query) || query.includes(p);
}

/**
 * Returns undefined if the package does not match the filter. Besides the
 * name and id, matches the studios a scraper covers and its URL patterns
 */
export function matchPackage(
  pkg: ISearchablePackage,
  filter: string
): IPackageMatch | undefined {
  const query = filter.trim().toLowerCase();
  if (!query) return {};

  if (
    pkg.name.toLowerCase().includes(query) ||
    pkg.package_id.toLowerCase().includes(query)
  ) {
    return {};
  }

  const metadata = pkg.metadata ?? {};

  const studio = stringList(metadata[SCRAPES_KEY]).find((s) =>
    s.toLowerCase().includes(query)
  );
  if (studio) return { via: studio };

  for (const [key, value] of Object.entries(metadata)) {
    if (!key.endsWith(URLS_KEY_SUFFIX)) continue;

    const pattern = stringList(value).find((p) => matchesURLPattern(p, query));
    if (pattern) return { via: pattern };
  }

  return undefined;
}
