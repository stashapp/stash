import queryString, { ParsedQuery } from "query-string";
import clone from "lodash-es/clone";
import {
  FilterMode,
  FindFilterType,
  SortDirectionEnum,
} from "src/core/generated-graphql";
import { Criterion, CriterionValue } from "./criteria/criterion";
import { makeCriteria } from "./criteria/factory";
import { DisplayMode } from "./types";

interface IQueryParameters {
  perPage?: string;
  sortby?: string;
  sortdir?: string;
  disp?: string;
  q?: string;
  p?: string;
  c?: string[];
  z?: string;
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
  public searchTerm?: string;
  public currentPage = DEFAULT_PARAMS.currentPage;
  public itemsPerPage = DEFAULT_PARAMS.itemsPerPage;
  public sortDirection: SortDirectionEnum = SortDirectionEnum.Asc;
  public sortBy?: string;
  public displayMode: DisplayMode = DEFAULT_PARAMS.displayMode;
  public zoomIndex: number = 1;
  public criteria: Array<Criterion<CriterionValue>> = [];
  public randomSeed = -1;
  private defaultZoomIndex: number = 1;

  public constructor(
    mode: FilterMode,
    defaultSort?: string,
    defaultDisplayMode?: DisplayMode,
    defaultZoomIndex?: number
  ) {
    this.mode = mode;
    this.sortBy = defaultSort;
    if (defaultDisplayMode !== undefined) this.displayMode = defaultDisplayMode;
    if (defaultZoomIndex !== undefined) {
      this.defaultZoomIndex = defaultZoomIndex;
      this.zoomIndex = defaultZoomIndex;
    }
  }

  public clone() {
    return Object.assign(new ListFilterModel(this.mode), this);
  }

  // Does not decode any URL-encoding in parameters
  public configureFromQueryParameters(params: IQueryParameters) {
    if (params.sortby !== undefined) {
      this.sortBy = params.sortby;

      // parse the random seed if provided
      const randomPrefix = "random_";
      if (this.sortBy && this.sortBy.startsWith(randomPrefix)) {
        const seedStr = this.sortBy.substring(randomPrefix.length);

        this.sortBy = "random";
        try {
          this.randomSeed = Number.parseInt(seedStr, 10);
        } catch (err) {
          // ignore
        }
      }
    }
    if (params.sortdir !== undefined) {
      this.sortDirection =
        params.sortdir === "desc"
          ? SortDirectionEnum.Desc
          : SortDirectionEnum.Asc;
    }
    if (params.disp !== undefined) {
      this.displayMode = Number.parseInt(params.disp, 10);
    }
    if (params.q) {
      this.searchTerm = params.q.trim();
    }
    this.currentPage = params.p ? Number.parseInt(params.p, 10) : 1;
    if (params.perPage) this.itemsPerPage = Number.parseInt(params.perPage, 10);
    if (params.z !== undefined) {
      const zoomIndex = Number.parseInt(params.z, 10);
      if (zoomIndex >= 0 && !Number.isNaN(zoomIndex)) {
        this.zoomIndex = zoomIndex;
      }
    }

    this.criteria = [];
    if (params.c !== undefined) {
      params.c.forEach((jsonString) => {
        try {
          const encodedCriterion = JSON.parse(jsonString);
          const criterion = makeCriteria(encodedCriterion.type);
          // it's possible that we have unsupported criteria. Just skip if so.
          if (criterion) {
            if (encodedCriterion.value !== undefined) {
              criterion.value = encodedCriterion.value;
            }
            criterion.modifier = encodedCriterion.modifier;
            this.criteria.push(criterion);
          }
        } catch (err) {
          // eslint-disable-next-line no-console
          console.error("Failed to parse encoded criterion:", err);
        }
      });
    }
  }

  public static decodeQueryParameters(
    parsedQuery: ParsedQuery<string>
  ): IQueryParameters {
    const params = clone(parsedQuery);
    if (params.q) {
      let searchTerm: string;
      if (params.q instanceof Array) {
        searchTerm = params.q[0];
      } else {
        searchTerm = params.q;
      }

      // See https://developer.mozilla.org/en-US/docs/Web/JavaScript/Reference/Global_Objects/decodeURIComponent#decoding_query_parameters_from_a_url
      searchTerm = searchTerm.replaceAll("+", " ");
      params.q = decodeURIComponent(searchTerm);
    }
    if (params.c !== undefined) {
      let jsonParameters: string[];
      if (params.c instanceof Array) {
        jsonParameters = params.c;
      } else {
        jsonParameters = [params.c!];
      }
      params.c = jsonParameters.map((jsonString) => {
        let decodedJson = jsonString;
        // replace () back to {}
        decodedJson = decodedJson.replaceAll("(", "{");
        decodedJson = decodedJson.replaceAll(")", "}");
        // decode all other characters
        decodedJson = decodeURIComponent(decodedJson);
        return decodedJson;
      });
    }
    return params;
  }

  public configureFromQueryString(query: string) {
    const parsed = queryString.parse(query, { decode: false });
    const decoded = ListFilterModel.decodeQueryParameters(parsed);
    this.configureFromQueryParameters(decoded);
  }

  public configureFromJSON(json: string) {
    this.configureFromQueryParameters(JSON.parse(json));
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
      return `${this.sortBy}_${this.randomSeed.toString()}`;
    }

    return this.sortBy;
  }

  // Returns query parameters with necessary parts encoded
  public getQueryParameters(): IQueryParameters {
    const encodedCriteria: string[] = this.criteria.map((criterion) => {
      let str = criterion.toJSON();
      // URL-encode other characters
      str = encodeURI(str);
      // force URL-encode existing ()
      str = str.replaceAll("(", "%28");
      str = str.replaceAll(")", "%29");
      // replace JSON '{'(%7B) '}'(%7D) with explicitly unreserved ()
      str = str.replaceAll("%7B", "(");
      str = str.replaceAll("%7D", ")");
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
      sortby: this.getSortBy() ?? undefined,
      sortdir:
        this.sortDirection === SortDirectionEnum.Desc ? "desc" : undefined,
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
      sortby: this.getSortBy() ?? undefined,
      sortdir:
        this.sortDirection === SortDirectionEnum.Desc ? "desc" : undefined,
      disp: this.displayMode,
      q: this.searchTerm,
      z: this.zoomIndex,
      c: encodedCriteria,
    };

    return JSON.stringify(result);
  }

  public makeQueryParameters(): string {
    return queryString.stringify(this.getQueryParameters(), { encode: false });
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
    // eslint-disable-next-line @typescript-eslint/no-explicit-any
    const output: Record<string, any> = {};
    this.criteria.forEach((criterion) => {
      criterion.apply(output);
    });

    return output;
  }
}
