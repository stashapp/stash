import { describe, it, expect } from "vitest";
import {
  buildFacetCandidates,
  buildEnumCandidates,
  FacetCandidateOptions,
} from "./facetCandidateUtils";
import { Option } from "./SidebarListFilter";
import { LabeledFacetCount } from "src/hooks/useFacetCounts";

describe("buildFacetCandidates", () => {
  const modifierOption: Option = {
    id: "any",
    label: "(Any)",
    className: "modifier-object",
  };

  const searchCandidates: Option[] = [
    modifierOption,
    { id: "1", label: "Performer A" },
    { id: "2", label: "Performer B" },
    { id: "3", label: "Performer C" },
  ];

  // Helper to create facet counts with labels
  const makeFacetCounts = (data: [string, number, string][]): Map<string, LabeledFacetCount> => {
    return new Map(data.map(([id, count, label]) => [id, { count, label }]));
  };

  describe("when facets are loaded and no search query", () => {
    it("should use facet results as candidates with labels from facets", () => {
      const facetCounts = makeFacetCounts([
        ["10", 50, "Studio X"], // Different IDs than search
        ["20", 30, "Studio Y"],
        ["30", 10, "Studio Z"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Should include modifier option
      expect(result[0]).toEqual(modifierOption);

      // Should use facet IDs and labels from facets
      const candidateIds = result.slice(1).map((c) => c.id);
      expect(candidateIds).toEqual(["10", "20", "30"]);

      // Should use labels from facets
      expect(result[1].label).toBe("Studio X");
      expect(result[2].label).toBe("Studio Y");
      expect(result[3].label).toBe("Studio Z");

      // Should be sorted by count descending
      expect(result[1].count).toBe(50);
      expect(result[2].count).toBe(30);
      expect(result[3].count).toBe(10);
    });

    it("should exclude selected items from candidates", () => {
      const facetCounts = makeFacetCounts([
        ["10", 50, "Studio X"],
        ["20", 30, "Studio Y"],
        ["30", 10, "Studio Z"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(["20"]), // ID 20 is selected
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);
      const candidateIds = result.slice(1).map((c) => c.id);

      // Should not include selected ID
      expect(candidateIds).not.toContain("20");
      expect(candidateIds).toEqual(["10", "30"]);
    });

    it("should exclude zero-count items from candidates", () => {
      const facetCounts = makeFacetCounts([
        ["10", 50, "Studio X"],
        ["20", 0, "Studio Y"], // Zero count
        ["30", 10, "Studio Z"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);
      const candidateIds = result.slice(1).map((c) => c.id);

      // Should not include zero-count item
      expect(candidateIds).not.toContain("20");
      expect(candidateIds).toEqual(["10", "30"]);
    });
  });

  describe("when there is a search query", () => {
    it("should use search results merged with facet counts", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Performer A"],
        ["2", 30, "Performer B"],
        // ID 3 is not in facets (undefined count)
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "search term",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Should include all search candidates
      expect(result.length).toBe(4); // modifier + 3 search results

      // Should have counts where available
      const item1 = result.find((c) => c.id === "1");
      expect(item1?.count).toBe(50);

      // Items not in facets should have undefined count but still be included
      const item3 = result.find((c) => c.id === "3");
      expect(item3).toBeDefined();
      expect(item3?.count).toBeUndefined();
    });

    it("should filter out zero-count search results", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Performer A"],
        ["2", 0, "Performer B"], // Zero count
        ["3", 10, "Performer C"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "search term",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);
      const candidateIds = result.filter((c) => c.className !== "modifier-object").map((c) => c.id);

      // Should not include zero-count item
      expect(candidateIds).not.toContain("2");
    });
  });

  describe("when facets are loading", () => {
    it("should return all search candidates without filtering", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Performer A"],
        // Only ID 1 in facets
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: true, // Still loading
      };

      const result = buildFacetCandidates(options);

      // Should return all search candidates (fallback behavior)
      expect(result.length).toBe(4);
      expect(result[0]).toEqual(modifierOption);
    });
  });

  describe("when facets are empty", () => {
    it("should return all search candidates", () => {
      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts: new Map(),
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Should return all search candidates
      expect(result.length).toBe(4);
    });
  });
});

describe("buildEnumCandidates", () => {
  const enumOptions: Option[] = [
    { id: "LANDSCAPE", label: "Landscape" },
    { id: "PORTRAIT", label: "Portrait" },
    { id: "SQUARE", label: "Square" },
  ];

  it("should filter out options with undefined counts", () => {
    const counts = new Map<string, number>([
      ["LANDSCAPE", 100],
      ["PORTRAIT", 50],
      // SQUARE is not in counts (undefined)
    ]);

    const result = buildEnumCandidates(
      enumOptions,
      new Set(),
      counts,
      false
    );

    expect(result.map((r) => r.id)).toEqual(["LANDSCAPE", "PORTRAIT"]);
    expect(result.map((r) => r.id)).not.toContain("SQUARE");
  });

  it("should filter out options with zero counts", () => {
    const counts = new Map<string, number>([
      ["LANDSCAPE", 100],
      ["PORTRAIT", 0], // Zero count
      ["SQUARE", 50],
    ]);

    const result = buildEnumCandidates(
      enumOptions,
      new Set(),
      counts,
      false
    );

    expect(result.map((r) => r.id)).toEqual(["LANDSCAPE", "SQUARE"]);
    expect(result.map((r) => r.id)).not.toContain("PORTRAIT");
  });

  it("should exclude selected items", () => {
    const counts = new Map<string, number>([
      ["LANDSCAPE", 100],
      ["PORTRAIT", 50],
      ["SQUARE", 25],
    ]);

    const result = buildEnumCandidates(
      enumOptions,
      new Set(["PORTRAIT"]), // PORTRAIT is selected
      counts,
      false
    );

    expect(result.map((r) => r.id)).not.toContain("PORTRAIT");
  });

  it("should return all options when counts are loading", () => {
    const counts = new Map<string, number>([
      ["LANDSCAPE", 100],
    ]);

    const result = buildEnumCandidates(
      enumOptions,
      new Set(),
      counts,
      true // Loading
    );

    expect(result.length).toBe(3);
  });

  it("should return all options when counts are undefined", () => {
    const result = buildEnumCandidates(
      enumOptions,
      new Set(),
      undefined,
      false
    );

    expect(result.length).toBe(3);
  });
});

// Regression tests for specific bug fixes
describe("Regression tests", () => {
  const makeFacetCounts = (data: [string, number, string][]): Map<string, LabeledFacetCount> => {
    return new Map(data.map(([id, count, label]) => [id, { count, label }]));
  };

  describe("Bug fix: Labels showing as IDs instead of names", () => {
    it("should use label from facets, not fallback to ID", () => {
      // This tests the fix where studios/performers were showing IDs
      // because labels were being discarded in toMap()
      const facetCounts = makeFacetCounts([
        ["123", 50, "My Studio Name"],
        ["456", 30, "Another Studio"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates: [], // No search candidates
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Should use label from facets, NOT the ID
      expect(result.find(c => c.id === "123")?.label).toBe("My Studio Name");
      expect(result.find(c => c.id === "456")?.label).toBe("Another Studio");
      
      // Should NOT show ID as label
      expect(result.find(c => c.label === "123")).toBeUndefined();
      expect(result.find(c => c.label === "456")).toBeUndefined();
    });
  });

  describe("Bug fix: Stale facet counts filtering candidates incorrectly", () => {
    it("should not filter by stale counts when facets are loading", () => {
      // Simulates scenario where filter is applied, facets are loading,
      // and stale counts shouldn't be used to filter candidates
      const staleFacetCounts = makeFacetCounts([
        ["1", 100, "Old Item"], // Stale data
      ]);

      const searchCandidates: Option[] = [
        { id: "1", label: "Item 1" },
        { id: "2", label: "Item 2" }, // Not in stale facets
        { id: "3", label: "Item 3" }, // Not in stale facets
      ];

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts: staleFacetCounts,
        facetsLoading: true, // Key: facets are still loading
      };

      const result = buildFacetCandidates(options);

      // Should return ALL candidates because facets are loading
      // Not just the ones in stale facets
      expect(result).toHaveLength(3);
      expect(result.map(c => c.id)).toContain("2");
      expect(result.map(c => c.id)).toContain("3");
    });

    it("should filter correctly when facets are fully loaded", () => {
      const facetCounts = makeFacetCounts([
        ["1", 100, "Item 1"],
        ["2", 0, "Item 2"], // Zero count - should be filtered
        ["3", 50, "Item 3"],
      ]);

      const searchCandidates: Option[] = [
        { id: "1", label: "Item 1" },
        { id: "2", label: "Item 2" },
        { id: "3", label: "Item 3" },
      ];

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "search",
        facetCounts,
        facetsLoading: false, // Facets fully loaded
      };

      const result = buildFacetCandidates(options);

      // Should filter out zero-count item
      expect(result.map(c => c.id)).not.toContain("2");
    });
  });

  describe("Bug fix: Mismatch between search results and facet results", () => {
    it("when no search query, should use facet results (TOP N) as candidates", () => {
      // Facets return TOP 100 by count
      const facetCounts = makeFacetCounts([
        ["top1", 1000, "Most Popular"],
        ["top2", 500, "Second Popular"],
        ["top3", 250, "Third Popular"],
      ]);

      // Search might return different set (sorted by relevance)
      const searchCandidates: Option[] = [
        { id: "random1", label: "Random Item" },
        { id: "random2", label: "Another Random" },
      ];

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "", // No search query
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Should use facet results, NOT search results
      expect(result.map(c => c.id)).toContain("top1");
      expect(result.map(c => c.id)).toContain("top2");
      expect(result.map(c => c.id)).toContain("top3");
      
      // Should NOT include random search results
      expect(result.map(c => c.id)).not.toContain("random1");
      expect(result.map(c => c.id)).not.toContain("random2");
    });

    it("when search query exists, should include search results not in TOP N", () => {
      // Facets only have TOP 3
      const facetCounts = makeFacetCounts([
        ["1", 1000, "Item 1"],
        ["2", 500, "Item 2"],
        ["3", 250, "Item 3"],
      ]);

      // Search found a specific item not in TOP N
      const searchCandidates: Option[] = [
        { id: "1", label: "Item 1" },
        { id: "99", label: "Specific Item Found" }, // Not in facets
      ];

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "Specific", // User searched for something
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Should include the searched item even though it's not in TOP N facets
      const item99 = result.find(c => c.id === "99");
      expect(item99).toBeDefined();
      expect(item99?.label).toBe("Specific Item Found");
      expect(item99?.count).toBeUndefined(); // No count because not in facets
    });
  });

  describe("Modifier options handling", () => {
    it("should preserve modifier options at the start of candidates", () => {
      const modifierOptions: Option[] = [
        { id: "any", label: "(Any)", className: "modifier-object" },
        { id: "none", label: "(None)", className: "modifier-object" },
      ];

      const facetCounts = makeFacetCounts([
        ["1", 50, "Item 1"],
        ["2", 30, "Item 2"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates: modifierOptions,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      // Modifier options should be preserved
      expect(result.filter(c => c.className === "modifier-object")).toHaveLength(2);
      
      // They should be at the start
      expect(result[0].className).toBe("modifier-object");
      expect(result[1].className).toBe("modifier-object");
    });
  });
});

