import queryString from "query-string";
import { RouteComponentProps } from "react-router-dom";
import { FilterMode } from "src/core/generated-graphql";
import { ListFilterModel } from "./list-filter/filter";
import { SceneListFilterOptions } from "./list-filter/scenes";

interface IQueryParameters {
  qsort?: string;
  qsortd?: string;
  qfq?: string;
  qfp?: string;
  qfc?: string[];
  qs?: string[];
}

export interface IPlaySceneOptions {
  sceneIndex?: number;
  newPage?: number;
  autoPlay?: boolean;
  continue?: boolean;
}

export class SceneQueue {
  public query?: ListFilterModel;
  public sceneIDs?: number[];
  private originalQueryPage?: number;
  private originalQueryPageSize?: number;

  public static fromListFilterModel(filter: ListFilterModel) {
    const ret = new SceneQueue();

    const filterCopy = filter.clone();
    filterCopy.itemsPerPage = 40;

    ret.originalQueryPage = filter.currentPage;
    ret.originalQueryPageSize = filter.itemsPerPage;

    ret.query = filterCopy;
    return ret;
  }

  public static fromSceneIDList(sceneIDs: string[]) {
    const ret = new SceneQueue();
    ret.sceneIDs = sceneIDs.map((v) => Number(v));
    return ret;
  }

  private makeQueryParameters(sceneIndex?: number, page?: number) {
    if (this.query) {
      const queryParams = this.query.getQueryParameters();
      const translatedParams = {
        qfp: queryParams.p ?? 1,
        qfc: queryParams.c,
        qfq: queryParams.q,
        qsort: queryParams.sortby,
        qsortd: queryParams.sortdir,
      };

      if (page !== undefined) {
        translatedParams.qfp = page;
      } else if (
        sceneIndex !== undefined &&
        this.originalQueryPage !== undefined &&
        this.originalQueryPageSize !== undefined
      ) {
        // adjust page to be correct for the index
        const filterIndex =
          sceneIndex +
          (this.originalQueryPage - 1) * this.originalQueryPageSize;
        const newPage = Math.floor(filterIndex / this.query.itemsPerPage) + 1;
        translatedParams.qfp = newPage;
      }

      return queryString.stringify(translatedParams, { encode: false });
    }

    if (this.sceneIDs && this.sceneIDs.length > 0) {
      const params = {
        qs: this.sceneIDs,
      };
      return queryString.stringify(params, { encode: false });
    }

    return "";
  }

  public static fromQueryParameters(params: string) {
    const ret = new SceneQueue();
    const parsed = queryString.parse(params) as IQueryParameters;
    const translated = {
      sortby: parsed.qsort,
      sortdir: parsed.qsortd,
      q: parsed.qfq,
      p: parsed.qfp,
      c: parsed.qfc,
    };

    if (parsed.qfp) {
      const query = new ListFilterModel(
        FilterMode.Scenes,
        translated as queryString.ParsedQuery,
        SceneListFilterOptions.defaultSortBy
      );
      ret.query = query;
    } else if (parsed.qs) {
      // must be scene list
      ret.sceneIDs = Array.isArray(parsed.qs)
        ? parsed.qs.map((v) => Number(v))
        : [Number(parsed.qs)];
    }

    return ret;
  }

  public playScene(
    history: RouteComponentProps["history"],
    sceneID: string,
    options?: IPlaySceneOptions
  ) {
    history.replace(this.makeLink(sceneID, options));
  }

  public makeLink(sceneID: string, options?: IPlaySceneOptions) {
    const params = [
      this.makeQueryParameters(options?.sceneIndex, options?.newPage),
      options?.autoPlay ? "autoplay=true" : "",
      options?.continue ? "continue=true" : "",
    ].filter((param) => !!param);
    return `/scenes/${sceneID}${params.length ? "?" + params.join("&") : ""}`;
  }
}

export default SceneQueue;
