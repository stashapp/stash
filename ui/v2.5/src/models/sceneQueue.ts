import queryString, { ParsedQuery } from "query-string";
import { FilterMode, Scene } from "src/core/generated-graphql";
import { ListFilterModel } from "./list-filter/filter";
import { SceneListFilterOptions } from "./list-filter/scenes";

export type QueuedScene = Pick<Scene, "id" | "title" | "paths">;

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

  public static fromQueryParameters(params: ParsedQuery<string>) {
    const ret = new SceneQueue();
    const translated = {
      sortby: params.qsort,
      sortdir: params.qsortd,
      q: params.qfq,
      p: params.qfp,
      c: params.qfc,
    };

    if (params.qfp) {
      const decoded = ListFilterModel.decodeQueryParameters(translated);
      const query = new ListFilterModel(
        FilterMode.Scenes,
        undefined,
        SceneListFilterOptions.defaultSortBy
      );
      query.configureFromQueryParameters(decoded);
      ret.query = query;
    } else if (params.qs) {
      // must be scene list
      ret.sceneIDs = Array.isArray(params.qs)
        ? params.qs.map((v) => Number(v))
        : [Number(params.qs)];
    }

    return ret;
  }

  public makeLink(sceneID: string, options: IPlaySceneOptions) {
    let params = [
      this.makeQueryParameters(options.sceneIndex, options.newPage),
    ];
    if (options.autoPlay !== undefined) {
      params.push("autoplay=" + options.autoPlay);
    }
    if (options.continue !== undefined) {
      params.push("continue=" + options.continue);
    }
    return `/scenes/${sceneID}${params.length ? "?" + params.join("&") : ""}`;
  }
}

export default SceneQueue;
