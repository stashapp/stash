import { FilterMode, Scene } from "src/core/generated-graphql";
import { ListFilterModel } from "./list-filter/filter";
import { INamedObject } from "src/utils/navigation";

export type QueuedScene = Pick<Scene, "id" | "title" | "date" | "paths"> & {
  performers?: INamedObject[] | null;
  studio?: INamedObject | null;
};

export interface IPlaySceneOptions {
  sceneIndex?: number;
  newPage?: number;
  autoPlay?: boolean;
  continue?: boolean;
  start?: number;
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
    const ret: string[] = [];

    if (this.query) {
      const queryParams = this.query.getEncodedParams();

      if (queryParams.sortby) {
        ret.push(`qsort=${queryParams.sortby}`);
      }
      if (queryParams.sortdir) {
        ret.push(`qsortd=${queryParams.sortdir}`);
      }
      if (queryParams.q) {
        ret.push(`qfq=${queryParams.q}`);
      }
      for (const c of queryParams.c ?? []) {
        ret.push(`qfc=${c}`);
      }

      let qfp = queryParams.p ?? "1";
      if (page !== undefined) {
        qfp = String(page);
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
        qfp = String(newPage);
      }
      ret.push(`qfp=${qfp}`);
    } else if (this.sceneIDs && this.sceneIDs.length > 0) {
      for (const id of this.sceneIDs) {
        ret.push(`qs=${id}`);
      }
    }

    return ret.join("&");
  }

  public static fromQueryParameters(params: URLSearchParams) {
    const ret = new SceneQueue();

    if (params.has("qfp")) {
      const translated = {
        sortby: params.get("qsort"),
        sortdir: params.get("qsortd"),
        q: params.get("qfq"),
        p: params.get("qfp"),
        c: params.getAll("qfc"),
      };
      const decoded = ListFilterModel.decodeParams(translated);
      const query = new ListFilterModel(FilterMode.Scenes);
      query.configureFromDecodedParams(decoded);
      ret.query = query;
    } else if (params.has("qs")) {
      // must be scene list
      ret.sceneIDs = params.getAll("qs").map((v) => Number(v));
    }

    return ret;
  }

  public makeLink(sceneID: string, options: IPlaySceneOptions) {
    let params = [
      this.makeQueryParameters(options.sceneIndex, options.newPage),
    ];
    if (options.autoPlay) {
      params.push("autoplay=true");
    }
    if (options.continue !== undefined) {
      params.push("continue=" + options.continue);
    }
    if (options.start !== undefined) {
      params.push("t=" + options.start);
    }
    return `/scenes/${sceneID}${params.length ? "?" + params.join("&") : ""}`;
  }
}

export default SceneQueue;
