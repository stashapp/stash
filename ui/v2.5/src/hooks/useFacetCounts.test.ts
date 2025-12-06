import { describe, it, expect } from "vitest";
import { LabeledFacetCount, FacetCounts } from "./useFacetCounts";

describe("LabeledFacetCount interface", () => {
  it("should store count and label", () => {
    const facet: LabeledFacetCount = {
      count: 42,
      label: "Test Label",
    };

    expect(facet.count).toBe(42);
    expect(facet.label).toBe("Test Label");
  });

  it("should work with Map storage", () => {
    const facetMap = new Map<string, LabeledFacetCount>();
    facetMap.set("123", { count: 10, label: "Item A" });
    facetMap.set("456", { count: 5, label: "Item B" });

    expect(facetMap.size).toBe(2);
    expect(facetMap.get("123")?.count).toBe(10);
    expect(facetMap.get("123")?.label).toBe("Item A");
    expect(facetMap.get("456")?.count).toBe(5);
    expect(facetMap.get("456")?.label).toBe("Item B");
  });
});

describe("FacetCounts interface", () => {
  it("should have proper structure for entity facets", () => {
    const emptyCounts: FacetCounts = {
      tags: new Map<string, LabeledFacetCount>(),
      performers: new Map<string, LabeledFacetCount>(),
      studios: new Map<string, LabeledFacetCount>(),
      groups: new Map<string, LabeledFacetCount>(),
      performerTags: new Map<string, LabeledFacetCount>(),
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      countries: new Map<string, LabeledFacetCount>(),
      circumcised: new Map(),
      ratings: new Map(),
      captions: new Map(),
      booleans: {
        organized: { true: 0, false: 0 },
        interactive: { true: 0, false: 0 },
        favorite: { true: 0, false: 0 },
      },
      parents: new Map<string, LabeledFacetCount>(),
      children: new Map<string, LabeledFacetCount>(),
    };

    // Verify structure exists
    expect(emptyCounts.tags).toBeInstanceOf(Map);
    expect(emptyCounts.performers).toBeInstanceOf(Map);
    expect(emptyCounts.studios).toBeInstanceOf(Map);
    expect(emptyCounts.groups).toBeInstanceOf(Map);
    expect(emptyCounts.booleans.organized).toEqual({ true: 0, false: 0 });
  });

  it("should support labeled facets for entity types", () => {
    const counts: FacetCounts = {
      tags: new Map([
        ["1", { count: 100, label: "Tag A" }],
        ["2", { count: 50, label: "Tag B" }],
      ]),
      performers: new Map([
        ["10", { count: 25, label: "Performer X" }],
      ]),
      studios: new Map([
        ["20", { count: 75, label: "Studio Y" }],
      ]),
      groups: new Map([
        ["30", { count: 15, label: "Group Z" }],
      ]),
      performerTags: new Map(),
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      countries: new Map(),
      circumcised: new Map(),
      ratings: new Map(),
      captions: new Map(),
      booleans: {
        organized: { true: 10, false: 90 },
        interactive: { true: 5, false: 95 },
        favorite: { true: 20, false: 80 },
      },
      parents: new Map(),
      children: new Map(),
    };

    // Verify labeled data
    expect(counts.tags.get("1")).toEqual({ count: 100, label: "Tag A" });
    expect(counts.performers.get("10")?.label).toBe("Performer X");
    expect(counts.studios.get("20")?.count).toBe(75);
    expect(counts.groups.get("30")?.label).toBe("Group Z");
    
    // Verify boolean counts
    expect(counts.booleans.organized.true).toBe(10);
    expect(counts.booleans.favorite.true).toBe(20);
  });
});

describe("Facet count map operations", () => {
  it("should iterate with forEach preserving labels", () => {
    const facets = new Map<string, LabeledFacetCount>([
      ["1", { count: 10, label: "Item 1" }],
      ["2", { count: 20, label: "Item 2" }],
      ["3", { count: 30, label: "Item 3" }],
    ]);

    const results: Array<{ id: string; count: number; label: string }> = [];
    facets.forEach((facetData, id) => {
      results.push({ id, count: facetData.count, label: facetData.label });
    });

    expect(results).toHaveLength(3);
    expect(results[0]).toEqual({ id: "1", count: 10, label: "Item 1" });
    expect(results[1]).toEqual({ id: "2", count: 20, label: "Item 2" });
    expect(results[2]).toEqual({ id: "3", count: 30, label: "Item 3" });
  });

  it("should filter by count correctly", () => {
    const facets = new Map<string, LabeledFacetCount>([
      ["1", { count: 10, label: "Item 1" }],
      ["2", { count: 0, label: "Item 2" }],  // Zero count
      ["3", { count: 30, label: "Item 3" }],
    ]);

    const nonZero: Array<{ id: string; count: number; label: string }> = [];
    facets.forEach((facetData, id) => {
      if (facetData.count > 0) {
        nonZero.push({ id, count: facetData.count, label: facetData.label });
      }
    });

    expect(nonZero).toHaveLength(2);
    expect(nonZero.map(r => r.id)).toEqual(["1", "3"]);
  });

  it("should handle empty maps gracefully", () => {
    const facets = new Map<string, LabeledFacetCount>();
    
    expect(facets.size).toBe(0);
    expect(facets.get("nonexistent")).toBeUndefined();
    
    // Iteration should work without errors
    const results: string[] = [];
    facets.forEach((_, id) => results.push(id));
    expect(results).toHaveLength(0);
  });
});

describe("Conversion from API response format", () => {
  it("should convert array of facet counts to Map with labels", () => {
    // Simulating the toMap function behavior
    const apiResponse = [
      { id: "1", label: "Studio A", count: 100 },
      { id: "2", label: "Studio B", count: 50 },
      { id: "3", label: "Studio C", count: 25 },
    ];

    const map = new Map<string, LabeledFacetCount>(
      apiResponse.map((c) => [c.id, { count: c.count, label: c.label }])
    );

    expect(map.size).toBe(3);
    expect(map.get("1")).toEqual({ count: 100, label: "Studio A" });
    expect(map.get("2")).toEqual({ count: 50, label: "Studio B" });
    expect(map.get("3")).toEqual({ count: 25, label: "Studio C" });
  });

  it("should handle boolean facet counts conversion", () => {
    const apiResponse = [
      { value: true, count: 150 },
      { value: false, count: 350 },
    ];

    const result = { true: 0, false: 0 };
    for (const c of apiResponse) {
      if (c.value) {
        result.true = c.count;
      } else {
        result.false = c.count;
      }
    }

    expect(result.true).toBe(150);
    expect(result.false).toBe(350);
  });

  it("should handle rating facet counts conversion", () => {
    const apiResponse = [
      { rating: 100, count: 50 },
      { rating: 80, count: 120 },
      { rating: 60, count: 200 },
      { rating: 40, count: 100 },
      { rating: 20, count: 30 },
    ];

    const map = new Map<number, number>(
      apiResponse.map((c) => [c.rating, c.count])
    );

    expect(map.size).toBe(5);
    expect(map.get(100)).toBe(50);
    expect(map.get(80)).toBe(120);
    expect(map.get(60)).toBe(200);
  });
});

describe("Lazy loading state preservation", () => {
  /**
   * These tests verify the lazy loading behavior for expensive facets.
   * When performer_tags or captions are lazily loaded, only those fields
   * should be updated while preserving existing data for other facets.
   * This prevents React rendering glitches where labels jump between filters.
   */

  it("should preserve existing facets when lazy loading performer_tags", () => {
    // Simulate existing state before lazy load
    const existingState: FacetCounts = {
      tags: new Map([["1", { count: 100, label: "Tag A" }]]),
      performers: new Map([["2", { count: 50, label: "Performer B" }]]),
      studios: new Map([["3", { count: 25, label: "Studio C" }]]),
      groups: new Map([["4", { count: 10, label: "Group D" }]]),
      performerTags: new Map(), // Empty before lazy load
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      countries: new Map(),
      circumcised: new Map(),
      ratings: new Map([[100, 20]]),
      captions: new Map(),
      booleans: {
        organized: { true: 10, false: 90 },
        interactive: { true: 5, false: 95 },
        favorite: { true: 0, false: 0 },
      },
      parents: new Map(),
      children: new Map(),
    };

    // Simulate lazy load update - only update performerTags
    const lazyLoadUpdate = (prev: FacetCounts): FacetCounts => ({
      ...prev,
      performerTags: new Map([
        ["100", { count: 500, label: "Performer Tag X" }],
        ["101", { count: 300, label: "Performer Tag Y" }],
      ]),
    });

    const newState = lazyLoadUpdate(existingState);

    // Verify performer_tags was updated
    expect(newState.performerTags.size).toBe(2);
    expect(newState.performerTags.get("100")).toEqual({ count: 500, label: "Performer Tag X" });

    // Verify all other facets are PRESERVED (same reference)
    expect(newState.tags).toBe(existingState.tags);
    expect(newState.performers).toBe(existingState.performers);
    expect(newState.studios).toBe(existingState.studios);
    expect(newState.groups).toBe(existingState.groups);
    expect(newState.ratings).toBe(existingState.ratings);
    expect(newState.booleans).toBe(existingState.booleans);

    // Verify data integrity
    expect(newState.tags.get("1")?.label).toBe("Tag A");
    expect(newState.performers.get("2")?.label).toBe("Performer B");
    expect(newState.studios.get("3")?.label).toBe("Studio C");
  });

  it("should preserve existing facets when lazy loading captions", () => {
    const existingState: FacetCounts = {
      tags: new Map([["1", { count: 100, label: "Tag A" }]]),
      performers: new Map([["2", { count: 50, label: "Performer B" }]]),
      studios: new Map(),
      groups: new Map(),
      performerTags: new Map([["10", { count: 200, label: "PT" }]]),
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      countries: new Map(),
      circumcised: new Map(),
      ratings: new Map(),
      captions: new Map(), // Empty before lazy load
      booleans: {
        organized: { true: 0, false: 0 },
        interactive: { true: 0, false: 0 },
        favorite: { true: 0, false: 0 },
      },
      parents: new Map(),
      children: new Map(),
    };

    // Simulate lazy load update - only update captions
    const lazyLoadUpdate = (prev: FacetCounts): FacetCounts => ({
      ...prev,
      captions: new Map([
        ["en", 1000],
        ["es", 500],
        ["de", 250],
      ]),
    });

    const newState = lazyLoadUpdate(existingState);

    // Verify captions was updated
    expect(newState.captions.size).toBe(3);
    expect(newState.captions.get("en")).toBe(1000);

    // Verify all other facets are PRESERVED
    expect(newState.tags).toBe(existingState.tags);
    expect(newState.performers).toBe(existingState.performers);
    expect(newState.performerTags).toBe(existingState.performerTags);
  });

  it("should NOT share references between different facet types", () => {
    // This test ensures that performer tags and performers are completely separate
    const performerTagsData: [string, LabeledFacetCount][] = [
      ["100", { count: 500, label: "No Tattoos" }],
      ["101", { count: 300, label: "Tongue Piercing" }],
    ];

    const performersData: [string, LabeledFacetCount][] = [
      ["1", { count: 50, label: "John Doe" }],
      ["2", { count: 30, label: "Jane Smith" }],
    ];

    const state: FacetCounts = {
      tags: new Map(),
      performers: new Map(performersData),
      studios: new Map(),
      groups: new Map(),
      performerTags: new Map(performerTagsData),
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      countries: new Map(),
      circumcised: new Map(),
      ratings: new Map(),
      captions: new Map(),
      booleans: {
        organized: { true: 0, false: 0 },
        interactive: { true: 0, false: 0 },
        favorite: { true: 0, false: 0 },
      },
      parents: new Map(),
      children: new Map(),
    };

    // Verify performers and performerTags are completely separate
    expect(state.performers).not.toBe(state.performerTags);
    expect(state.performers.size).toBe(2);
    expect(state.performerTags.size).toBe(2);

    // Verify no label contamination
    expect(state.performers.get("1")?.label).toBe("John Doe");
    expect(state.performers.get("100")).toBeUndefined(); // Performer tag ID
    expect(state.performerTags.get("100")?.label).toBe("No Tattoos");
    expect(state.performerTags.get("1")).toBeUndefined(); // Performer ID
  });

  it("should handle simultaneous lazy load of both performer_tags and captions", () => {
    const existingState: FacetCounts = {
      tags: new Map([["1", { count: 100, label: "Tag A" }]]),
      performers: new Map([["2", { count: 50, label: "Performer B" }]]),
      studios: new Map(),
      groups: new Map(),
      performerTags: new Map(), // Empty
      resolutions: new Map(),
      orientations: new Map(),
      genders: new Map(),
      countries: new Map(),
      circumcised: new Map(),
      ratings: new Map(),
      captions: new Map(), // Empty
      booleans: {
        organized: { true: 0, false: 0 },
        interactive: { true: 0, false: 0 },
        favorite: { true: 0, false: 0 },
      },
      parents: new Map(),
      children: new Map(),
    };

    // Simulate lazy load update for both
    const lazyLoadUpdate = (prev: FacetCounts): FacetCounts => ({
      ...prev,
      performerTags: new Map([["100", { count: 500, label: "PT" }]]),
      captions: new Map([["en", 1000]]),
    });

    const newState = lazyLoadUpdate(existingState);

    // Verify both were updated
    expect(newState.performerTags.size).toBe(1);
    expect(newState.captions.size).toBe(1);

    // Verify core facets preserved
    expect(newState.tags).toBe(existingState.tags);
    expect(newState.performers).toBe(existingState.performers);
  });
});

