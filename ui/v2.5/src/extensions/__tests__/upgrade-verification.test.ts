/**
 * Upgrade Verification Tests
 *
 * These tests verify that critical extension utilities and types work correctly.
 * Run these after every upstream upgrade to catch issues early.
 *
 * Usage: yarn test --run
 */

import { describe, it, expect } from "vitest";
import { LabeledFacetCount } from "src/extensions/hooks/useFacetCounts";
import {
  buildFacetCandidates,
  buildEnumCandidates,
  FacetCandidateOptions,
} from "src/extensions/filters/facetCandidateUtils";
import { Option } from "src/components/List/Filters/SidebarListFilter";

describe("Upgrade Verification - Critical Functions", () => {
  describe("Facet candidate utilities exist", () => {
    it("should have buildFacetCandidates function", () => {
      expect(typeof buildFacetCandidates).toBe("function");
    });

    it("should have buildEnumCandidates function", () => {
      expect(typeof buildEnumCandidates).toBe("function");
    });
  });

  describe("Facet utilities work correctly", () => {
    const makeFacetCounts = (
      data: [string, number, string][]
    ): Map<string, LabeledFacetCount> => {
      return new Map(data.map(([id, count, label]) => [id, { count, label }]));
    };

    it("buildFacetCandidates should create candidates from facets", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Item A"],
        ["2", 30, "Item B"],
        ["3", 10, "Item C"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates: [],
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      expect(result).toHaveLength(3);
      expect(result[0]).toEqual({ id: "1", count: 50, label: "Item A" });
      expect(result[1]).toEqual({ id: "2", count: 30, label: "Item B" });
      expect(result[2]).toEqual({ id: "3", count: 10, label: "Item C" });
    });

    it("buildFacetCandidates should filter zero counts", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Item A"],
        ["2", 0, "Item B"], // Zero count
        ["3", 10, "Item C"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates: [],
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      expect(result).toHaveLength(2);
      expect(result.map((c) => c.id)).not.toContain("2");
    });

    it("buildEnumCandidates should filter by counts", () => {
      const enumOptions: Option[] = [
        { id: "LANDSCAPE", label: "Landscape" },
        { id: "PORTRAIT", label: "Portrait" },
        { id: "SQUARE", label: "Square" },
      ];

      const counts = new Map<string, number>([
        ["LANDSCAPE", 100],
        ["PORTRAIT", 50],
        // SQUARE not in counts - should be filtered
      ]);

      const result = buildEnumCandidates(enumOptions, new Set(), counts, false);

      // Should only include items with counts
      expect(result.map((r) => r.id)).toContain("LANDSCAPE");
      expect(result.map((r) => r.id)).toContain("PORTRAIT");
      expect(result.map((r) => r.id)).not.toContain("SQUARE");
    });
  });
});

describe("Upgrade Verification - Type Structures", () => {
  describe("LabeledFacetCount type works", () => {
    it("should create valid LabeledFacetCount objects", () => {
      const facet: LabeledFacetCount = {
        count: 42,
        label: "Test Label",
      };

      expect(facet.count).toBe(42);
      expect(facet.label).toBe("Test Label");
    });

    it("should work in Map structures", () => {
      const map = new Map<string, LabeledFacetCount>();
      map.set("test-id", { count: 10, label: "Test" });

      const retrieved = map.get("test-id");
      expect(retrieved?.count).toBe(10);
      expect(retrieved?.label).toBe("Test");
    });
  });
});

describe("Upgrade Verification - Import Paths", () => {
  /**
   * These tests verify critical import paths work.
   * Static imports at top of file already verify the main paths.
   * This section tests additional paths that don't have React dependencies.
   */

  it("should import facetCandidateUtils without error", async () => {
    const module = await import("src/extensions/filters/facetCandidateUtils");
    expect(module.buildFacetCandidates).toBeDefined();
    expect(module.buildEnumCandidates).toBeDefined();
  });

  it("LabeledFacetCount type should be importable", () => {
    // This is tested by the static import at the top of this file
    // If this test runs, the import worked
    const testFacet: LabeledFacetCount = { count: 1, label: "test" };
    expect(testFacet).toBeDefined();
  });
});

describe("Upgrade Verification - Regression Prevention", () => {
  const makeFacetCounts = (
    data: [string, number, string][]
  ): Map<string, LabeledFacetCount> => {
    return new Map(data.map(([id, count, label]) => [id, { count, label }]));
  };

  describe("Labels must come from facet data, not IDs", () => {
    it("should preserve labels from facet data", () => {
      const facetCounts = makeFacetCounts([
        ["123", 50, "Actual Label Name"],
      ]);

      const options: FacetCandidateOptions = {
        searchCandidates: [],
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts,
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      expect(result[0].label).toBe("Actual Label Name");
      expect(result[0].label).not.toBe("123"); // Must NOT use ID as label
    });
  });

  describe("Empty facet maps should not crash", () => {
    it("should handle empty facet maps gracefully", () => {
      const options: FacetCandidateOptions = {
        searchCandidates: [],
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts: new Map(),
        facetsLoading: false,
      };

      const result = buildFacetCandidates(options);

      expect(result).toHaveLength(0);
      expect(Array.isArray(result)).toBe(true);
    });
  });

  describe("Facets loading state should return all candidates", () => {
    it("should not filter when facets are loading", () => {
      const searchCandidates: Option[] = [
        { id: "1", label: "Item 1" },
        { id: "2", label: "Item 2" },
        { id: "3", label: "Item 3" },
      ];

      const options: FacetCandidateOptions = {
        searchCandidates,
        selectedIds: new Set(),
        searchQuery: "",
        facetCounts: new Map(), // Empty
        facetsLoading: true, // But loading
      };

      const result = buildFacetCandidates(options);

      // Should return all candidates when loading
      expect(result).toHaveLength(3);
    });
  });
});
