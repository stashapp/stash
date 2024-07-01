import {
  ConfigDataFragment,
  FilterMode,
  FindFilterType,
  SavedFilterDataFragment,
  SortDirectionEnum,
} from "src/core/generated-graphql";
import {
  Criterion,
  CriterionValue,
  ISavedCriterion,
} from "./criteria/criterion";
import { getFilterOptions } from "./factory";
import {
  CriterionType,
  DisplayMode,
  SavedObjectFilter,
  SavedUIOptions,
} from "./types";
import { ListFilterOptions } from "./filter-options";

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
  public readonly mode: FilterMode;
  public readonly options: ListFilterOptions;
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
    defaultZoomIndex?: number
  ) {
    this.mode = mode;
    this.config = config;
    this.options = getFilterOptions(mode);
    const { defaultSortBy, displayModeOptions } = this.options;

    this.sortBy = defaultSortBy;
    if (this.sortBy === "date") {
      this.sortDirection = SortDirectionEnum.Desc;
    }
    this.displayMode = displayModeOptions[0];
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

  public configureFromDecodedParams(params: IDecodedParams) {
    if (params.perPage !== undefined) {
      this.itemsPerPage = params.perPage;
    }
    if (params.sortby !== undefined) {
      this.sortBy = params.sortby;

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
    }
    if (params.q !== undefined) {
      this.searchTerm = params.q;
    }
    this.currentPage = params.p ?? 1;
    if (params.z !== undefined) {
      this.zoomIndex = params.z;
    }

    this.criteria = [];
    if (params.c !== undefined) {
      for (const jsonString of params.c) {
        try {
          const { type: criterionType, ...savedCriterion } =
            JSON.parse(jsonString);

          const criterion = this.makeCriterion(criterionType);
          criterion.setFromSavedCriterion(savedCriterion);

          this.criteria.push(criterion);
        } catch (err) {
          // eslint-disable-next-line no-console
          console.error("Failed to parse encoded criterion:", err);
        }
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

  public configureFromQueryString(queryString: string) {
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
    this.configureFromDecodedParams(decoded);
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
    this.sortDirection = findFilter?.direction ?? this.sortDirection;
    this.searchTerm = findFilter?.q ?? this.searchTerm;

    this.displayMode = uiOptions?.display_mode ?? this.displayMode;
    this.zoomIndex = uiOptions?.zoom_index ?? this.zoomIndex;

    this.currentPage = 1;

    this.criteria = [];
    if (objectFilter) {
      for (const [k, v] of Object.entries(objectFilter)) {
        const criterion = this.makeCriterion(k as CriterionType);
        criterion.setFromSavedCriterion(v as ISavedCriterion<CriterionValue>);
        this.criteria.push(criterion);
      }
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
    for (const c of this.criteria) {
      output[c.criterionOption.type] = c.toCriterionInput();
    }
    return output;
  }

  public makeSavedFilter() {
    const output: SavedObjectFilter = {};
    for (const c of this.criteria) {
      output[c.criterionOption.type] = c.toSavedCriterion();
    }
    return output;
  }

  public makeSavedUIOptions(): SavedUIOptions {
    return {
      display_mode: this.displayMode,
      zoom_index: this.zoomIndex,
    };
  }
}
