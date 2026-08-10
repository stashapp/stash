import { describe, expect, it } from "vitest";
import { FilterMode, SortDirectionEnum } from "src/core/generated-graphql";
import { ListFilterModel } from "./list-filter/filter";
import { AudioQueue } from "./audioQueue";

function makeFilter(overrides?: Partial<ListFilterModel>) {
  const filter = new ListFilterModel(FilterMode.Audios);
  filter.sortBy = "title";
  filter.sortDirection = SortDirectionEnum.Desc;
  filter.itemsPerPage = 20;
  filter.currentPage = 1;
  Object.assign(filter, overrides);
  return filter;
}

describe("AudioQueue", () => {
  describe("fromAudioIDList", () => {
    it("makes a link containing every queued id", () => {
      const queue = AudioQueue.fromAudioIDList(["1", "2", "3"]);

      expect(queue.audioIDs).toEqual([1, 2, 3]);
      expect(queue.makeLink("2", {})).toBe("/audios/2?qs=1&qs=2&qs=3");
    });

    it("appends the play options", () => {
      const queue = AudioQueue.fromAudioIDList(["1", "2"]);

      expect(
        queue.makeLink("1", { autoPlay: true, continue: true, start: 30 })
      ).toBe("/audios/1?qs=1&qs=2&autoplay=true&continue=true&t=30");
    });

    it("includes continue=false when explicitly disabled", () => {
      const queue = AudioQueue.fromAudioIDList(["1"]);

      expect(queue.makeLink("1", { continue: false })).toBe(
        "/audios/1?qs=1&continue=false"
      );
    });
  });

  describe("fromListFilterModel", () => {
    it("queues 40 items per page regardless of the list page size", () => {
      const queue = AudioQueue.fromListFilterModel(makeFilter());

      expect(queue.query?.itemsPerPage).toBe(40);
    });

    it("encodes the filter into the link", () => {
      const filter = makeFilter();
      filter.searchTerm = "coltrane";
      const queue = AudioQueue.fromListFilterModel(filter);

      const link = queue.makeLink("5", {});
      expect(link.startsWith("/audios/5?")).toBe(true);

      const params = new URLSearchParams(link.split("?")[1]);
      expect(params.get("qsort")).toBe("title");
      expect(params.get("qsortd")).toBe("desc");
      expect(params.get("qfq")).toBe("coltrane");
      expect(params.get("qfp")).toBe("1");
    });

    it("maps an audio index on a later list page to the queue page", () => {
      // list page 3 of 20 per page => index 5 is overall index 45,
      // which is on queue page 2 (40 per queue page)
      const filter = makeFilter({ currentPage: 3 });
      const queue = AudioQueue.fromListFilterModel(filter);

      const params = new URLSearchParams(
        queue.makeLink("46", { audioIndex: 5 }).split("?")[1]
      );
      expect(params.get("qfp")).toBe("2");
    });

    it("prefers an explicit new page over the audio index", () => {
      const queue = AudioQueue.fromListFilterModel(
        makeFilter({ currentPage: 3 })
      );

      const params = new URLSearchParams(
        queue.makeLink("1", { audioIndex: 5, newPage: 7 }).split("?")[1]
      );
      expect(params.get("qfp")).toBe("7");
    });
  });

  describe("fromQueryParameters", () => {
    it("round-trips a filter queue", () => {
      const filter = makeFilter();
      filter.searchTerm = "ellington";
      const queue = AudioQueue.fromListFilterModel(filter);

      const link = queue.makeLink("1", { audioIndex: 0 });
      const restored = AudioQueue.fromQueryParameters(
        new URLSearchParams(link.split("?")[1])
      );

      expect(restored.audioIDs).toBeUndefined();
      expect(restored.query?.mode).toBe(FilterMode.Audios);
      expect(restored.query?.sortBy).toBe("title");
      expect(restored.query?.sortDirection).toBe(SortDirectionEnum.Desc);
      expect(restored.query?.searchTerm).toBe("ellington");
      expect(restored.query?.currentPage).toBe(1);
    });

    it("round-trips an id list queue", () => {
      const queue = AudioQueue.fromAudioIDList(["3", "1", "2"]);

      const link = queue.makeLink("3", {});
      const restored = AudioQueue.fromQueryParameters(
        new URLSearchParams(link.split("?")[1])
      );

      expect(restored.query).toBeUndefined();
      expect(restored.audioIDs).toEqual([3, 1, 2]);
    });

    it("returns an empty queue when there are no queue parameters", () => {
      const restored = AudioQueue.fromQueryParameters(
        new URLSearchParams("autoplay=true")
      );

      expect(restored.query).toBeUndefined();
      expect(restored.audioIDs).toBeUndefined();
      expect(restored.makeLink("1", {})).toBe("/audios/1?");
    });
  });
});
