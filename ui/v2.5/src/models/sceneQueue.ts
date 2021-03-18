import { ListFilterModel } from "./list-filter/filter";
import queryString from "query-string";
import { FilterMode } from "./list-filter/types";

interface IQueryParameters {
  qsort?: string;
  qsortd?: string;
  qfq?: string;
  qfp?: string;
  qfc?: string[];
  qs?: string[];
}

export class SceneQueue {
  public query?: ListFilterModel;
  public sceneIDs?: string[];

  public static fromListFilterModel(filter: ListFilterModel, currentSceneIndex?: number) {
    const ret = new SceneQueue();

    const filterCopy = Object.assign(new ListFilterModel(filter.filterMode), filter);
    filterCopy.itemsPerPage = 40;

    // adjust page to be correct for the index
    const filterIndex = currentSceneIndex !== undefined ? currentSceneIndex + ((filter.currentPage - 1) * filter.itemsPerPage) : 0;
    const newPage = Math.floor(filterIndex / filterCopy.itemsPerPage) + 1;
    filterCopy.currentPage = newPage;

    ret.query = filterCopy;
    return ret;
  }

  public makeQueryParameters(page?: number) {
    if (this.query) {
      const queryParams = this.query.getQueryParameters();
      const translatedParams = {
        qfp: queryParams.p ?? 1,
        qfc: queryParams.c,
        qfq: queryParams.q,
        qsort: queryParams.sortby,
        qsortd: queryParams.sortdir,
      }

      if (page !== undefined) {
        translatedParams.qfp = page;
      }
      
      return queryString.stringify(translatedParams, { encode: false });
    } else if (this.sceneIDs && this.sceneIDs.length > 0) {
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
      const query = new ListFilterModel(FilterMode.Scenes, translated as queryString.ParsedQuery);
      ret.query = query;
    } else if (parsed.qs) {
      // must be scene list
      ret.sceneIDs = parsed.qs;
    }

    return ret;
  }
}