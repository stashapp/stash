import {
  ConfigDataFragment,
  FilterMode,
  FindFilterType,
  SavedFilterDataFragment,
  SortDirectionEnum,
} from "src/core/generated-graphql";
import { Criterion, CriterionValue } from "./criteria/criterion";
import { getFilterOptions } from "./factory";
import { CriterionType, DisplayMode } from "./types";
import * as GQL from "src/core/generated-graphql";
import { useContext, useMemo } from "react";
import { ConfigurationContext } from "src/hooks/Config";
import { View } from "src/components/List/views";
import { DefaultFilters, IUIConfig } from "src/core/config";

interface IDecodedParams {
  perPage?: number;
  sortby?: string;
  sortdir?: string;
  disp?: DisplayMode;
  q?: string;
  p?: number;
  z?: number;
  c?: string[];
}

interface IEncodedParams {
  perPage?: string | null;
  sortby?: string | null;
  sortdir?: string | null;
  disp?: string | null;
  q?: string | null;
  p?: string | null;
  z?: string | null;
  c?: string[];
}

const DEFAULT_PARAMS = {
  sortDirection: SortDirectionEnum.Asc,
  displayMode: DisplayMode.Grid,
  currentPage: 1,
  itemsPerPage: 40,
};

// TODO: handle customCriteria
export class ListFilterModel {
  public mode: FilterMode;
  private config?: ConfigDataFragment;
  public searchTerm: string = "";
  public currentPage = DEFAULT_PARAMS.currentPage;
  public itemsPerPage = DEFAULT_PARAMS.itemsPerPage;
  public sortDirection: SortDirectionEnum = DEFAULT_PARAMS.sortDirection;
  public sortBy?: string;
  public displayMode: DisplayMode = DEFAULT_PARAMS.displayMode;
  public zoomIndex: number = 1;
  public criteria: Array<Criterion<CriterionValue>> = [];
  public randomSeed = -1;
  private defaultZoomIndex: number = 1;

  public constructor(
    mode: FilterMode,
    config?: ConfigDataFragment,
    defaultSort?: string,
    defaultDisplayMode?: DisplayMode,
    defaultZoomIndex?: number
  ) {
    this.mode = mode;
    this.config = config;
    this.sortBy = defaultSort;
    if (this.sortBy === "date") {
      this.sortDirection = SortDirectionEnum.Desc;
    }
    if (defaultDisplayMode !== undefined) {
      this.displayMode = defaultDisplayMode;
    }
    if (defaultZoomIndex !== undefined) {
      this.defaultZoomIndex = defaultZoomIndex;
      this.zoomIndex = defaultZoomIndex;
    }
  }

  public clone() {
    return Object.assign(new ListFilterModel(this.mode, this.config), this);
  }

  // returns the number of filters applied
  public count() {
    // don't include search term
    return this.criteria.length;
  }

  public configureFromDecodedParams(
    params: IDecodedParams,
    defaultFilter: ListFilterModel | undefined = undefined
  ) {
    if (params.perPage !== undefined) {
      this.itemsPerPage = params.perPage;
    } else if (
      defaultFilter !== undefined &&
      defaultFilter.itemsPerPage !== undefined
    ) {
      this.itemsPerPage = defaultFilter.itemsPerPage;
    }
    if (params.sortby !== undefined) {
      this.sortBy = params.sortby;

      // parse the random seed if provided
      const match = this.sortBy.match(/^random_(\d+)$/);
      if (match) {
        this.sortBy = "random";
        this.randomSeed = Number.parseInt(match[1], 10);
      }
    } else if (
      defaultFilter !== undefined &&
      defaultFilter.sortBy !== undefined
    ) {
      this.sortBy = defaultFilter.sortBy;

      // parse the random seed if provided
      const match = this.sortBy.match(/^random_(\d+)$/);
      if (match) {
        this.sortBy = "random";
        this.randomSeed = Number.parseInt(match[1], 10);
      }
    }
    if (params.sortdir !== undefined) {
      this.sortDirection =
        params.sortdir === "desc"
          ? SortDirectionEnum.Desc
          : SortDirectionEnum.Asc;
    } else if (
      defaultFilter !== undefined &&
      defaultFilter.sortDirection !== undefined
    ) {
      this.sortDirection = defaultFilter.sortDirection;
    } else {
      // #3193 - sortdir undefined means asc
      // #3559 - unless sortby is date, then desc
      this.sortDirection =
        params.sortby === "date"
          ? SortDirectionEnum.Desc
          : SortDirectionEnum.Asc;
    }
    if (params.disp !== undefined) {
      this.displayMode = params.disp;
    } else if (
      defaultFilter !== undefined &&
      defaultFilter.displayMode !== undefined
    ) {
      this.displayMode = defaultFilter.displayMode;
    }
    if (params.q !== undefined) {
      this.searchTerm = params.q;
    } else if (
      defaultFilter !== undefined &&
      defaultFilter.searchTerm !== undefined
    ) {
      this.searchTerm = defaultFilter.searchTerm;
    }
    this.currentPage = params.p ?? defaultFilter?.currentPage ?? 1;
    if (params.z !== undefined) {
      this.zoomIndex = params.z;
    } else if (
      defaultFilter !== undefined &&
      defaultFilter.zoomIndex !== undefined
    ) {
      this.zoomIndex = defaultFilter.zoomIndex;
    }

    this.criteria = [];
    if (params.c !== undefined) {
      for (const jsonString of params.c) {
        try {
          const encodedCriterion = JSON.parse(jsonString);
          const criterion = this.makeCriterion(encodedCriterion.type);
          criterion.setFromEncodedCriterion(encodedCriterion);
          this.criteria.push(criterion);
        } catch (err) {
          // eslint-disable-next-line no-console
          console.error("Failed to parse encoded criterion:", err);
        }
      }
    }
    if (defaultFilter !== undefined && defaultFilter.criteria !== undefined) {
      for (const criterion of defaultFilter.criteria) {
        this.criteria.push(criterion);
      }
    }
  }

  // Does not decode any URL-encoding, only type conversions
  public static decodeParams(params: IEncodedParams): IDecodedParams {
    const ret: IDecodedParams = {};

    if (params.perPage) {
      ret.perPage = Number.parseInt(params.perPage, 10);
    }
    if (params.sortby) {
      ret.sortby = params.sortby;
    }
    if (params.sortdir) {
      ret.sortdir = params.sortdir;
    }
    if (params.disp) {
      ret.disp = Number.parseInt(params.disp, 10);
    }
    if (params.q) {
      ret.q = params.q.trim();
    }
    if (params.p) {
      ret.p = Number.parseInt(params.p, 10);
    }
    if (params.z) {
      const zoomIndex = Number.parseInt(params.z, 10);
      if (zoomIndex >= 0) {
        ret.z = zoomIndex;
      }
    }

    if (params.c && params.c.length !== 0) {
      ret.c = params.c.map((jsonString) =>
        ListFilterModel.translateJSON(jsonString, true)
      );
    }

    return ret;
  }

  private static translateJSON(jsonString: string, decoding: boolean) {
    let inString = false;
    let escape = false;
    return [...jsonString]
      .map((c) => {
        if (escape) {
          // this character has been escaped, skip
          escape = false;
          return c;
        }

        switch (c) {
          case "\\":
            // escape the next character if in a string
            if (inString) {
              escape = true;
            }
            break;
          case '"':
            // unescaped quote, toggle inString
            inString = !inString;
            break;
          case "(":
            // decode only: restore ( to { if not in a string
            if (decoding && !inString) {
              return "{";
            }
            break;
          case ")":
            // decode only: restore ) to } if not in a string
            if (decoding && !inString) {
              return "}";
            }
            break;
          case "{":
            // encode only: replace { with ( if not in a string
            if (!decoding && !inString) {
              return "(";
            }
            break;
          case "}":
            // encode only: replace } with ) if not in a string
            if (!decoding && !inString) {
              return ")";
            }
            break;
        }

        return c;
      })
      .join("");
  }

  public configureFromQueryString(
    queryString: string,
    defaultFilter: ListFilterModel | undefined = undefined
  ) {
    const query = new URLSearchParams(queryString);
    const params = {
      perPage: query.get("perPage"),
      sortby: query.get("sortby"),
      sortdir: query.get("sortdir"),
      disp: query.get("disp"),
      q: query.get("q"),
      p: query.get("p"),
      z: query.get("z"),
      c: query.getAll("c"),
    };
    const decoded = ListFilterModel.decodeParams(params);
    this.configureFromDecodedParams(decoded, defaultFilter);
  }

  public configureFromSavedFilter(savedFilter: SavedFilterDataFragment) {
    const {
      find_filter: findFilter,
      object_filter: objectFilter,
      ui_options: uiOptions,
    } = savedFilter;

    this.itemsPerPage = findFilter?.per_page ?? this.itemsPerPage;
    this.sortBy = findFilter?.sort ?? this.sortBy;
    // parse the random seed if provided
    const match = this.sortBy?.match(/^random_(\d+)$/);
    if (match) {
      this.sortBy = "random";
      this.randomSeed = Number.parseInt(match[1], 10);
    }
    this.sortDirection =
      (findFilter?.direction as SortDirectionEnum) ?? this.sortDirection;
    this.searchTerm = findFilter?.q ?? this.searchTerm;

    this.displayMode = uiOptions?.display_mode ?? this.displayMode;
    this.zoomIndex = uiOptions?.zoom_index ?? this.zoomIndex;

    this.currentPage = 1;

    this.criteria = [];
    if (objectFilter) {
      Object.keys(objectFilter).forEach((key) => {
        const criterion = this.makeCriterion(key as CriterionType);
        criterion.setFromEncodedCriterion(objectFilter[key]);
        this.criteria.push(criterion);
      });
    }
  }

  private setRandomSeed() {
    if (this.sortBy === "random") {
      // #321 - set the random seed if it is not set
      if (this.randomSeed === -1) {
        // generate 8-digit seed
        this.randomSeed = Math.floor(Math.random() * 10 ** 8);
      }
    } else {
      this.randomSeed = -1;
    }
  }

  private getSortBy(): string | undefined {
    this.setRandomSeed();

    if (this.sortBy === "random") {
      return `random_${this.randomSeed.toString()}`;
    }

    return this.sortBy;
  }

  // Returns query parameters with necessary parts URL-encoded
  public getEncodedParams(): IEncodedParams {
    const encodedCriteria: string[] = this.criteria.map((criterion) => {
      let str = ListFilterModel.translateJSON(criterion.toJSON(), false);

      // URL-encode other characters
      str = encodeURI(str);

      // only the reserved characters ?#&;=+ need to be URL-encoded
      // as they have special meaning in query strings
      str = str.replaceAll("?", encodeURIComponent("?"));
      str = str.replaceAll("#", encodeURIComponent("#"));
      str = str.replaceAll("&", encodeURIComponent("&"));
      str = str.replaceAll(";", encodeURIComponent(";"));
      str = str.replaceAll("=", encodeURIComponent("="));
      str = str.replaceAll("+", encodeURIComponent("+"));

      return str;
    });

    return {
      perPage:
        this.itemsPerPage !== DEFAULT_PARAMS.itemsPerPage
          ? String(this.itemsPerPage)
          : undefined,
      sortby: this.getSortBy(),
      sortdir:
        this.sortBy === "date"
          ? this.sortDirection === SortDirectionEnum.Asc
            ? "asc"
            : undefined
          : this.sortDirection === SortDirectionEnum.Desc
          ? "desc"
          : undefined,
      disp:
        this.displayMode !== DEFAULT_PARAMS.displayMode
          ? String(this.displayMode)
          : undefined,
      q: this.searchTerm ? encodeURIComponent(this.searchTerm) : undefined,
      p:
        this.currentPage !== DEFAULT_PARAMS.currentPage
          ? String(this.currentPage)
          : undefined,
      z:
        this.zoomIndex !== this.defaultZoomIndex
          ? String(this.zoomIndex)
          : undefined,
      c: encodedCriteria,
    };
  }

  public makeSavedFilterJSON() {
    const encodedCriteria: string[] = this.criteria.map((criterion) =>
      criterion.toJSON()
    );

    const result = {
      perPage: this.itemsPerPage,
      sortby: this.getSortBy(),
      sortdir:
        this.sortBy === "date"
          ? this.sortDirection === SortDirectionEnum.Asc
            ? "asc"
            : undefined
          : this.sortDirection === SortDirectionEnum.Desc
          ? "desc"
          : undefined,
      disp: this.displayMode,
      q: this.searchTerm || undefined,
      z: this.zoomIndex,
      c: encodedCriteria,
    };

    return JSON.stringify(result);
  }

  public makeQueryParameters(): string {
    const query: string[] = [];
    const params = this.getEncodedParams();

    if (params.q) {
      query.push(`q=${params.q}`);
    }
    if (params.c) {
      for (const c of params.c) {
        query.push(`c=${c}`);
      }
    }
    if (params.sortby) {
      query.push(`sortby=${params.sortby}`);
    }
    if (params.sortdir) {
      query.push(`sortdir=${params.sortdir}`);
    }
    if (params.perPage) {
      query.push(`perPage=${params.perPage}`);
    }
    if (params.disp) {
      query.push(`disp=${params.disp}`);
    }
    if (params.z) {
      query.push(`z=${params.z}`);
    }
    if (params.p) {
      query.push(`p=${params.p}`);
    }

    return query.join("&");
  }

  public makeCriterion(type: CriterionType) {
    const { criterionOptions } = getFilterOptions(this.mode);

    const option = criterionOptions.find((o) => o.type === type);

    if (!option) {
      throw new Error(`Unknown criterion parameter name: ${type}`);
    }

    return option.makeCriterion(this.config);
  }

  // TODO: These don't support multiple of the same criteria, only the last one set is used.

  public makeFindFilter(): FindFilterType {
    return {
      q: this.searchTerm,
      page: this.currentPage,
      per_page: this.itemsPerPage,
      sort: this.getSortBy(),
      direction: this.sortDirection,
    };
  }

  public makeFilter() {
    const output: Record<string, unknown> = {};
    this.criteria.forEach((criterion) => {
      criterion.apply(output);
    });

    return output;
  }

  public makeSavedFindFilter() {
    const output: Record<string, { value: CriterionValue; modifier: string }> =
      {};
    this.criteria.forEach((criterion) => {
      criterion.toSavedFilter(output);
    });

    return output;
  }

  public makeUIOptions(): Record<string, unknown> {
    return {
      display_mode: this.displayMode,
      zoom_index: this.zoomIndex,
    };
  }
}

// const recursiveRenameToSnakeCase = (
//   camelCaseObject: Record<string, unknown>
// ) => {
//   let camelCaseKeys = Object.keys(camelCaseObject);
//   camelCaseKeys.forEach((key) => {
//     if (
//       typeof camelCaseObject[key] === "object" &&
//       !Array.isArray(camelCaseObject[key]) &&
//       camelCaseObject[key] !== null
//     ) {
//       let cco: Record<string, unknown> = Object(camelCaseObject[key]);
//       camelCaseObject[key] = recursiveRenameToSnakeCase(cco);
//     }
//     let snakeCaseKey = key.replace(/([a-z])([A-Z])/g, "$1_$2").toLowerCase();
//     if (snakeCaseKey !== key) {
//       camelCaseObject[snakeCaseKey] = camelCaseObject[key];
//       delete camelCaseObject[key];
//     }
//   });
//   return camelCaseObject;
// };

export const useDefaultFilter = (mode: GQL.FilterMode, view?: View) => {
  let { configuration: config, loading } = useContext(ConfigurationContext);
  const { defaultFilters } = config?.ui as IUIConfig;

  const defaultFilter = useMemo(() => {
    // TODO - this is a horrible temporary workaround for viper
    let parsed: DefaultFilters;

    try {
      parsed = JSON.parse(defaultFilters ?? "{}");
    } catch (e) {
      // eslint-disable-next-line no-console
      console.error("Failed to parse default filters:", e);
      return undefined;
    }

    const savedFilter = view ? parsed[view] : undefined;
    if (!view || !savedFilter) return undefined;

    let filter = new ListFilterModel(mode, config);
    filter.configureFromSavedFilter(savedFilter);
    filter.currentPage = 1;
    filter.randomSeed = -1;
    return filter;
  }, [config, mode, view, defaultFilters]);

  return { defaultFilter, loading };
};
