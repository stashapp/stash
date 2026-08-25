import { Audio, FilterMode } from "src/core/generated-graphql";
import { ListFilterModel } from "./list-filter/filter";
import { INamedObject } from "src/utils/navigation";

export type QueuedAudio = Pick<Audio, "id" | "title" | "date" | "paths"> & {
  performers?: INamedObject[] | null;
  studio?: INamedObject | null;
};

export interface IPlayAudioOptions {
  audioIndex?: number;
  newPage?: number;
  autoPlay?: boolean;
  continue?: boolean;
  start?: number;
}

export class AudioQueue {
  public query?: ListFilterModel;
  public audioIDs?: number[];
  private originalQueryPage?: number;
  private originalQueryPageSize?: number;

  public static fromListFilterModel(filter: ListFilterModel) {
    const ret = new AudioQueue();

    const filterCopy = filter.clone();
    filterCopy.itemsPerPage = 40;

    ret.originalQueryPage = filter.currentPage;
    ret.originalQueryPageSize = filter.itemsPerPage;

    ret.query = filterCopy;
    return ret;
  }

  public static fromAudioIDList(audioIDs: string[]) {
    const ret = new AudioQueue();
    ret.audioIDs = audioIDs.map((v) => Number(v));
    return ret;
  }

  private makeQueryParameters(audioIndex?: number, page?: number) {
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
        audioIndex !== undefined &&
        this.originalQueryPage !== undefined &&
        this.originalQueryPageSize !== undefined
      ) {
        // adjust page to be correct for the index
        const filterIndex =
          audioIndex +
          (this.originalQueryPage - 1) * this.originalQueryPageSize;
        const newPage = Math.floor(filterIndex / this.query.itemsPerPage) + 1;
        qfp = String(newPage);
      }
      ret.push(`qfp=${qfp}`);
    } else if (this.audioIDs && this.audioIDs.length > 0) {
      for (const id of this.audioIDs) {
        ret.push(`qs=${id}`);
      }
    }

    return ret.join("&");
  }

  public static fromQueryParameters(params: URLSearchParams) {
    const ret = new AudioQueue();

    if (params.has("qfp")) {
      const translated = {
        sortby: params.get("qsort"),
        sortdir: params.get("qsortd"),
        q: params.get("qfq"),
        p: params.get("qfp"),
        c: params.getAll("qfc"),
      };
      const decoded = ListFilterModel.decodeParams(translated);
      const query = new ListFilterModel(FilterMode.Audios);
      query.configureFromDecodedParams(decoded);
      ret.query = query;
    } else if (params.has("qs")) {
      // must be audio list
      ret.audioIDs = params.getAll("qs").map((v) => Number(v));
    }

    return ret;
  }

  public makeLink(audioID: string, options: IPlayAudioOptions) {
    const params = [
      this.makeQueryParameters(options.audioIndex, options.newPage),
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
    return `/audios/${audioID}${params.length ? "?" + params.join("&") : ""}`;
  }
}

export default AudioQueue;
