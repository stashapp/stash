import { describe, it, expect } from "vitest";
import {
  normalizeRatingSystemType,
  normalizeRatingStarPrecision,
  normalizeRatingOptions,
  getRatingPrecision,
  getRatingStep,
  getRatingMax,
  convertToRatingFormat,
  convertFromRatingFormat,
  formatDisplayRating,
  getRatingInputLabel,
  getRatingBannerColor,
  defaultRatingSystemOptions,
} from "../components/Rating";

describe("normalizeRatingSystemType", () => {
  it("returns 'stars' for null", () => {
    expect(normalizeRatingSystemType(null)).toBe("stars");
  });

  it("returns 'stars' for undefined", () => {
    expect(normalizeRatingSystemType(undefined)).toBe("stars");
  });

  it("returns 'decimal' for 'decimal'", () => {
    expect(normalizeRatingSystemType("decimal")).toBe("decimal");
  });

  it("returns 'decimal' case-insensitive", () => {
    expect(normalizeRatingSystemType("Decimal")).toBe("decimal");
  });

  it("returns 'stars' for any other value", () => {
    expect(normalizeRatingSystemType("other")).toBe("stars");
  });
});

describe("normalizeRatingStarPrecision", () => {
  it("returns 'full' for null", () => {
    expect(normalizeRatingStarPrecision(null)).toBe("full");
  });

  it("returns 'half'", () => {
    expect(normalizeRatingStarPrecision("half")).toBe("half");
  });

  it("returns 'quarter'", () => {
    expect(normalizeRatingStarPrecision("quarter")).toBe("quarter");
  });

  it("returns 'tenth'", () => {
    expect(normalizeRatingStarPrecision("tenth")).toBe("tenth");
  });

  it("returns 'full' for unknown values", () => {
    expect(normalizeRatingStarPrecision("whatever")).toBe("full");
  });
});

describe("normalizeRatingOptions", () => {
  it("returns defaults for null", () => {
    expect(normalizeRatingOptions(null)).toEqual(defaultRatingSystemOptions);
  });

  it("normalizes partial options", () => {
    expect(normalizeRatingOptions({ type: "decimal" })).toEqual({
      type: "decimal",
      starPrecision: "full",
    });
  });
});

describe("getRatingPrecision", () => {
  it("returns 1 for full", () => {
    expect(getRatingPrecision("full")).toBe(1);
  });

  it("returns 0.5 for half", () => {
    expect(getRatingPrecision("half")).toBe(0.5);
  });

  it("returns 0.25 for quarter", () => {
    expect(getRatingPrecision("quarter")).toBe(0.25);
  });

  it("returns 0.1 for tenth", () => {
    expect(getRatingPrecision("tenth")).toBe(0.1);
  });
});

describe("getRatingStep", () => {
  it("returns 1 for default star options", () => {
    expect(getRatingStep()).toBe(1);
  });

  it("returns 0.1 for decimal type", () => {
    expect(getRatingStep({ type: "decimal", starPrecision: "full" })).toBe(0.1);
  });

  it("returns 0.5 for half precision stars", () => {
    expect(getRatingStep({ type: "stars", starPrecision: "half" })).toBe(0.5);
  });
});

describe("getRatingMax", () => {
  it("returns 5 for stars", () => {
    expect(getRatingMax()).toBe(5);
  });

  it("returns 10 for decimal", () => {
    expect(getRatingMax({ type: "decimal", starPrecision: "full" })).toBe(10);
  });
});

describe("convertToRatingFormat", () => {
  it("returns null for null value", () => {
    expect(convertToRatingFormat(null)).toBeNull();
  });

  it("returns null for undefined value", () => {
    expect(convertToRatingFormat(undefined)).toBeNull();
  });

  it("returns null for 0", () => {
    expect(convertToRatingFormat(0)).toBeNull();
  });

  it("converts 100 to 5 stars", () => {
    expect(convertToRatingFormat(100)).toBe(5);
  });

  it("converts 50 to 2.5 stars (half precision)", () => {
    expect(convertToRatingFormat(50, { type: "stars", starPrecision: "half" })).toBe(2.5);
  });

  it("converts 100 to 10.0 decimal", () => {
    expect(convertToRatingFormat(100, { type: "decimal", starPrecision: "full" })).toBe(10);
  });

  it("converts 75 to 7.5 decimal", () => {
    expect(convertToRatingFormat(75, { type: "decimal", starPrecision: "full" })).toBe(7.5);
  });
});

describe("convertFromRatingFormat", () => {
  it("converts 5 stars to 100", () => {
    expect(convertFromRatingFormat(5)).toBe(100);
  });

  it("converts 3 stars to 60", () => {
    expect(convertFromRatingFormat(3)).toBe(60);
  });

  it("converts 10.0 decimal to 100", () => {
    expect(convertFromRatingFormat(10, { type: "decimal", starPrecision: "full" })).toBe(100);
  });

  it("converts 7.5 decimal to 75", () => {
    expect(convertFromRatingFormat(7.5, { type: "decimal", starPrecision: "full" })).toBe(75);
  });
});

describe("formatDisplayRating", () => {
  it("returns null for null", () => {
    expect(formatDisplayRating(null)).toBeNull();
  });

  it("formats 100 as 5 stars", () => {
    expect(formatDisplayRating(100)).toBe("5");
  });

  it("formats 50 as 2.5 with half precision", () => {
    expect(formatDisplayRating(50, { type: "stars", starPrecision: "half" })).toBe("2.5");
  });

  it("formats 75 as 7.5 decimal", () => {
    expect(formatDisplayRating(75, { type: "decimal", starPrecision: "full" })).toBe("7.5");
  });
});

describe("getRatingInputLabel", () => {
  it("returns star label for default", () => {
    expect(getRatingInputLabel()).toBe("Rating (0-5 stars)");
  });

  it("returns decimal label", () => {
    expect(getRatingInputLabel({ type: "decimal", starPrecision: "full" })).toBe("Rating (0-10.0)");
  });

  it("includes step for fractional precision", () => {
    expect(getRatingInputLabel({ type: "stars", starPrecision: "half" })).toBe("Rating (0-5 stars, step 0.5)");
  });
});

describe("getRatingBannerColor", () => {
  it("returns gray-ish for low rating", () => {
    const color = getRatingBannerColor(10);
    expect(color).toMatch(/^#[0-9a-fA-F]{6}$/);
  });

  it("returns red-ish for high rating", () => {
    const color = getRatingBannerColor(100);
    expect(color).toBe("#ff0000");
  });

  it("returns a valid hex color for mid rating", () => {
    const color = getRatingBannerColor(50);
    expect(color).toMatch(/^#[0-9a-fA-F]{6}$/);
  });

  it("returns different colors for different ratings", () => {
    const c1 = getRatingBannerColor(20);
    const c2 = getRatingBannerColor(80);
    expect(c1).not.toBe(c2);
  });
});
