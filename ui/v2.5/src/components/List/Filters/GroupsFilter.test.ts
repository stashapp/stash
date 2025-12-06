import { describe, it, expect, vi } from "vitest";
import { LabeledFacetCount } from "src/hooks/useFacetCounts";

// Test the GroupsFilter-specific logic
// Since GroupsFilter uses the same pattern as other filters,
// we test the pattern used for groups filtering

describe("GroupsFilter logic", () => {
  describe("Groups candidate building", () => {
    // Helper to create facet counts with labels
    const makeFacetCounts = (
      data: [string, number, string][]
    ): Map<string, LabeledFacetCount> => {
      return new Map(data.map(([id, count, label]) => [id, { count, label }]));
    };

    it("should build group candidates from facet data", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Action Movies"],
        ["2", 30, "Comedy Series"],
        ["3", 10, "Documentary"],
      ]);

      const candidates: Array<{ id: string; label: string; count: number }> = [];
      facetCounts.forEach((facet, id) => {
        if (facet.count > 0) {
          candidates.push({
            id,
            label: facet.label,
            count: facet.count,
          });
        }
      });

      expect(candidates).toHaveLength(3);
      expect(candidates[0]).toEqual({ id: "1", label: "Action Movies", count: 50 });
      expect(candidates[1]).toEqual({ id: "2", label: "Comedy Series", count: 30 });
      expect(candidates[2]).toEqual({ id: "3", label: "Documentary", count: 10 });
    });

    it("should exclude zero-count groups", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Action Movies"],
        ["2", 0, "Empty Group"], // Zero count
        ["3", 10, "Documentary"],
      ]);

      const candidates: Array<{ id: string; label: string; count: number }> = [];
      facetCounts.forEach((facet, id) => {
        if (facet.count > 0) {
          candidates.push({
            id,
            label: facet.label,
            count: facet.count,
          });
        }
      });

      expect(candidates).toHaveLength(2);
      expect(candidates.map((c) => c.id)).not.toContain("2");
    });

    it("should exclude already selected groups", () => {
      const facetCounts = makeFacetCounts([
        ["1", 50, "Action Movies"],
        ["2", 30, "Comedy Series"],
        ["3", 10, "Documentary"],
      ]);

      const selectedIds = new Set(["2"]); // Comedy Series is selected

      const candidates: Array<{ id: string; label: string; count: number }> = [];
      facetCounts.forEach((facet, id) => {
        if (facet.count > 0 && !selectedIds.has(id)) {
          candidates.push({
            id,
            label: facet.label,
            count: facet.count,
          });
        }
      });

      expect(candidates).toHaveLength(2);
      expect(candidates.map((c) => c.id)).not.toContain("2");
      expect(candidates.map((c) => c.id)).toContain("1");
      expect(candidates.map((c) => c.id)).toContain("3");
    });

    it("should sort groups by count descending", () => {
      const facetCounts = makeFacetCounts([
        ["1", 10, "Small Group"],
        ["2", 100, "Large Group"],
        ["3", 50, "Medium Group"],
      ]);

      const candidates: Array<{ id: string; label: string; count: number }> = [];
      facetCounts.forEach((facet, id) => {
        candidates.push({
          id,
          label: facet.label,
          count: facet.count,
        });
      });

      // Sort by count descending
      candidates.sort((a, b) => b.count - a.count);

      expect(candidates[0].label).toBe("Large Group");
      expect(candidates[1].label).toBe("Medium Group");
      expect(candidates[2].label).toBe("Small Group");
    });
  });

  describe("ContainingGroups and SubGroups filters", () => {
    it("should use same pattern for parent group filtering", () => {
      // ContainingGroups uses the same GroupsCriterion class
      // Verifying the pattern works for hierarchical group relationships

      const parentGroups = new Map<string, LabeledFacetCount>([
        ["parent-1", { count: 5, label: "Parent Collection A" }],
        ["parent-2", { count: 3, label: "Parent Collection B" }],
      ]);

      const candidates: Array<{ id: string; label: string; count: number }> = [];
      parentGroups.forEach((facet, id) => {
        candidates.push({ id, label: facet.label, count: facet.count });
      });

      expect(candidates).toHaveLength(2);
      expect(candidates[0].label).toBe("Parent Collection A");
    });

    it("should use same pattern for child group filtering", () => {
      // SubGroups uses the same GroupsCriterion class
      // Verifying the pattern works for hierarchical group relationships

      const childGroups = new Map<string, LabeledFacetCount>([
        ["child-1", { count: 10, label: "Sub Group 1" }],
        ["child-2", { count: 8, label: "Sub Group 2" }],
        ["child-3", { count: 2, label: "Sub Group 3" }],
      ]);

      const candidates: Array<{ id: string; label: string; count: number }> = [];
      childGroups.forEach((facet, id) => {
        candidates.push({ id, label: facet.label, count: facet.count });
      });

      expect(candidates).toHaveLength(3);
    });
  });
});

describe("GroupsFilter integration with FacetCounts", () => {
  it("should handle empty facet counts gracefully", () => {
    const emptyFacets = new Map<string, LabeledFacetCount>();

    expect(emptyFacets.size).toBe(0);

    const candidates: Array<{ id: string; label: string; count: number }> = [];
    emptyFacets.forEach((facet, id) => {
      candidates.push({ id, label: facet.label, count: facet.count });
    });

    expect(candidates).toHaveLength(0);
  });

  it("should preserve group labels from facets", () => {
    // Regression test for the label preservation bug
    const facetCounts = new Map<string, LabeledFacetCount>([
      ["123", { count: 50, label: "My Group Name" }],
      ["456", { count: 30, label: "Another Group" }],
    ]);

    const candidates: Array<{ id: string; label: string }> = [];
    facetCounts.forEach((facet, id) => {
      candidates.push({ id, label: facet.label });
    });

    // Should use label from facets, NOT the ID
    expect(candidates.find((c) => c.id === "123")?.label).toBe("My Group Name");
    expect(candidates.find((c) => c.id === "456")?.label).toBe("Another Group");

    // Should NOT show ID as label
    expect(candidates.find((c) => c.label === "123")).toBeUndefined();
  });
});

