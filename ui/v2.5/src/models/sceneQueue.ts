import queryString from "query-string";
import { RouteComponentProps } from "react-router-dom";
import { ListFilterModel } from "./list-filter/filter";
import { FilterMode } from "./list-filter/types";

interface IQueryParameters {
  qsort?: string;
  qsortd?: string;
  qfq?: string;
  qfp?: string;
  qfc?: string[];
  qs?: string[];
}

export interface IPlaySceneOptions {
  newPage?: number;
  autoPlay?: boolean;
}

export class SceneQueue {
  public query?: ListFilterModel;
  public sceneIDs?: number[];

  public static fromListFilterModel(
    filter: ListFilterModel,
    currentSceneIndex?: number
  ) {
    const ret = new SceneQueue();

    const filterCopy = Object.assign(
      new ListFilterModel(filter.filterMode),
      filter
    );
    filterCopy.itemsPerPage = 40;

    // adjust page to be correct for the index
    const filterIndex =
      currentSceneIndex !== undefined
        ? currentSceneIndex + (filter.currentPage - 1) * filter.itemsPerPage
        : 0;
    const newPage = Math.floor(filterIndex / filterCopy.itemsPerPage) + 1;
    filterCopy.currentPage = newPage;

    ret.query = filterCopy;
    return ret;
  }

  public static fromSceneIDList(sceneIDs: string[]) {
    const ret = new SceneQueue();
    ret.sceneIDs = sceneIDs.map((v) => Number(v));
    return ret;
  }

  private makeQueryParameters(page?: number) {
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
        translated as queryString.ParsedQuery
      );
      ret.query = query;
    } else if (parsed.qs) {
      // must be scene list
      ret.sceneIDs = parsed.qs.map((v) => Number(v));
    }

    return ret;
  }

  public playScene(
    history: RouteComponentProps["history"],
    sceneID: string,
    options?: IPlaySceneOptions
  ) {
    const paramStr = this.makeQueryParameters(options?.newPage);
    const autoplayParam = options?.autoPlay ? "&autoplay=true" : "";
    history.push(`/scenes/${sceneID}?${paramStr}${autoplayParam}`);
  }
}
