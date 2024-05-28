interface ISortable {
  id: string;
}

// sortByRelevance is a function that sorts an array of objects by relevance to a query string.
// It uses the following priorities:
// 1. Exact matches
// 2. Starts with
// 3. Word matches
// 4. Word starts with
// 5. Includes
// If aliases are provided, they are also checked in the same order, but with lower priority than
// the name of the object.
export function sortByRelevance<T extends ISortable>(
  query: string,
  value: T[],
  getName: (o: T) => string,
  getAliases?: (o: T) => string[] | undefined
) {
  if (!query) {
    return value;
  }

  query = query.toLowerCase();

  interface ICacheEntry {
    aliases?: string[];
    aliasMatch?: boolean;
    aliasStartsWith?: boolean;
    wordIndex?: number;
    wordStartsWithIndex?: number;
    aliasWordIndex?: number;
    aliasWordStartsWithIndex?: number;
    aliasIncludesIndex?: number;
  }

  const cache: Record<string, ICacheEntry> = {};

  function setCache(tag: T, partial: Partial<ICacheEntry>) {
    cache[tag.id] = {
      ...cache[tag.id],
      ...partial,
    };
  }

  function getObjectAliases(o: T) {
    const cached = cache[o.id]?.aliases;

    if (cached !== undefined) {
      return cached;
    }

    if (!getAliases) {
      return [];
    }

    const aliases = getAliases(o)?.map((a) => a.toLowerCase()) ?? [];
    setCache(o, { aliases });

    return aliases;
  }

  function aliasMatches(o: T) {
    const cached = cache[o.id]?.aliasMatch;

    if (cached !== undefined) {
      return cached;
    }

    const aliases = getObjectAliases(o);
    const aliasMatch = aliases.some((a) => a === query);
    setCache(o, { aliasMatch });

    return aliasMatch;
  }

  function aliasStartsWith(o: T) {
    const cached = cache[o.id]?.aliasStartsWith;

    if (cached !== undefined) {
      return cached;
    }

    const aliases = getObjectAliases(o);
    const startsWith = aliases.some((a) => a.startsWith(query));
    setCache(o, { aliasStartsWith: startsWith });

    return startsWith;
  }

  function getWords(o: T) {
    return getName(o).trim().toLowerCase().split(" ");
  }

  function getAliasWords(tag: T) {
    const aliases = getObjectAliases(tag);
    return aliases.map((a) => a.trim().split(" ")).flat();
  }

  function getWordIndex(o: T) {
    const cached = cache[o.id]?.wordIndex;

    if (cached !== undefined) {
      return cached;
    }

    const words = getWords(o);
    const wordIndex = words.findIndex((w) => w === query);
    setCache(o, { wordIndex });

    return wordIndex;
  }

  function getAliasWordIndex(o: T) {
    const cached = cache[o.id]?.aliasWordIndex;

    if (cached !== undefined) {
      return cached;
    }

    const aliasWords = getAliasWords(o);
    const aliasWordIndex = aliasWords.findIndex((w) => w === query);
    setCache(o, { aliasWordIndex });

    return aliasWordIndex;
  }

  function getWordStartsWithIndex(o: T) {
    const cached = cache[o.id]?.wordStartsWithIndex;

    if (cached !== undefined) {
      return cached;
    }

    const words = getWords(o);
    const wordStartsWithIndex = words.findIndex((w) => w.startsWith(query));
    setCache(o, { wordStartsWithIndex });

    return wordStartsWithIndex;
  }

  function getAliasWordStartsWithIndex(o: T) {
    const cached = cache[o.id]?.aliasWordStartsWithIndex;

    if (cached !== undefined) {
      return cached;
    }

    const aliasWords = getAliasWords(o);
    const aliasWordStartsWithIndex = aliasWords.findIndex((w) =>
      w.startsWith(query)
    );
    setCache(o, { aliasWordStartsWithIndex });

    return aliasWordStartsWithIndex;
  }

  function getAliasIncludesIndex(o: T) {
    const cached = cache[o.id]?.aliasIncludesIndex;

    if (cached !== undefined) {
      return cached;
    }

    const aliases = getObjectAliases(o);
    const aliasIncludesIndex = aliases.findIndex((a) => a.includes(query));
    setCache(o, { aliasIncludesIndex });

    return aliasIncludesIndex;
  }

  function compare(a: T, b: T) {
    const aName = getName(a).toLowerCase();
    const bName = getName(b).toLowerCase();

    const aAlias = aliasMatches(a);
    const bAlias = aliasMatches(b);

    // exact matches first
    if (aName === query && bName !== query) {
      return -1;
    }

    if (aName !== query && bName === query) {
      return 1;
    }

    if (aAlias && !bAlias) {
      return -1;
    }

    if (!aAlias && bAlias) {
      return 1;
    }

    // then starts with
    if (aName.startsWith(query) && !bName.startsWith(query)) {
      return -1;
    }

    if (!aName.startsWith(query) && bName.startsWith(query)) {
      return 1;
    }

    const aAliasStartsWith = aliasStartsWith(a);
    const bAliasStartsWith = aliasStartsWith(b);

    if (aAliasStartsWith && !bAliasStartsWith) {
      return -1;
    }

    if (!aAliasStartsWith && bAliasStartsWith) {
      return 1;
    }

    // only check words if the query is a single word
    if (!query.includes(" ")) {
      // word matches
      {
        const aWord = getWordIndex(a);
        const bWord = getWordIndex(b);

        if (aWord !== -1 && bWord === -1) {
          return -1;
        }

        if (aWord === -1 && bWord !== -1) {
          return 1;
        }

        if (aWord !== -1 && bWord !== -1) {
          if (aWord === bWord) {
            return aName.localeCompare(bName);
          }

          return aWord - bWord;
        }

        const aAliasWord = getAliasWordIndex(a);
        const bAliasWord = getAliasWordIndex(b);

        if (aAliasWord !== -1 && bAliasWord === -1) {
          return -1;
        }

        if (aAliasWord === -1 && bAliasWord !== -1) {
          return 1;
        }

        if (aAliasWord !== -1 && bAliasWord !== -1) {
          if (aAliasWord === bAliasWord) {
            return aName.localeCompare(bName);
          }

          return aAliasWord - bAliasWord;
        }
      }

      // then start of word
      {
        const aWord = getWordStartsWithIndex(a);
        const bWord = getWordStartsWithIndex(b);

        if (aWord !== -1 && bWord === -1) {
          return -1;
        }

        if (aWord === -1 && bWord !== -1) {
          return 1;
        }

        if (aWord !== -1 && bWord !== -1) {
          if (aWord === bWord) {
            return aName.localeCompare(bName);
          }

          return aWord - bWord;
        }

        const aAliasWord = getAliasWordStartsWithIndex(a);
        const bAliasWord = getAliasWordStartsWithIndex(b);

        if (aAliasWord !== -1 && bAliasWord === -1) {
          return -1;
        }

        if (aAliasWord === -1 && bAliasWord !== -1) {
          return 1;
        }

        if (aAliasWord !== -1 && bAliasWord !== -1) {
          if (aAliasWord === bAliasWord) {
            return aName.localeCompare(bName);
          }

          return aAliasWord - bAliasWord;
        }
      }
    }

    // then contains
    // performance of this is presumably fast enough to not require caching
    const aNameIncludeIndex = aName.indexOf(query);
    const bNameIncludeIndex = bName.indexOf(query);

    if (aNameIncludeIndex !== -1 && bNameIncludeIndex === -1) {
      return -1;
    }

    if (aNameIncludeIndex === -1 && bNameIncludeIndex !== -1) {
      return 1;
    }

    if (aNameIncludeIndex !== -1 && bNameIncludeIndex !== -1) {
      if (aNameIncludeIndex === bNameIncludeIndex) {
        return aName.localeCompare(bName);
      }

      return aNameIncludeIndex - bNameIncludeIndex;
    }

    const aAliasIncludes = getAliasIncludesIndex(a);
    const bAliasIncludes = getAliasIncludesIndex(b);

    if (aAliasIncludes !== -1 && bAliasIncludes === -1) {
      return -1;
    }

    if (aAliasIncludes === -1 && bAliasIncludes !== -1) {
      return 1;
    }

    if (aAliasIncludes !== -1 && bAliasIncludes !== -1) {
      if (aAliasIncludes === bAliasIncludes) {
        return aName.localeCompare(bName);
      }

      return aAliasIncludes - bAliasIncludes;
    }

    return aName.localeCompare(bName);
  }

  return value.slice().sort(compare);
}
