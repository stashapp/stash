import { describe, it, expect } from "vitest";
import { formatDuration, formatFileSize, formatDate, getResolutionLabel } from "../components/shared";

describe("formatDuration", () => {
  it("returns 0:00 for zero seconds", () => {
    expect(formatDuration(0)).toBe("0:00");
  });

  it("returns 0:00 for negative values", () => {
    expect(formatDuration(-5)).toBe("0:00");
  });

  it("formats seconds only", () => {
    expect(formatDuration(45)).toBe("0:45");
  });

  it("formats minutes and seconds", () => {
    expect(formatDuration(125)).toBe("2:05");
  });

  it("pads seconds with zero", () => {
    expect(formatDuration(61)).toBe("1:01");
  });

  it("formats hours, minutes, and seconds", () => {
    expect(formatDuration(3661)).toBe("1:01:01");
  });

  it("pads minutes and seconds in hour format", () => {
    expect(formatDuration(3600)).toBe("1:00:00");
  });

  it("handles large durations", () => {
    expect(formatDuration(36000)).toBe("10:00:00");
  });
});

describe("formatFileSize", () => {
  it("returns 0 B for zero bytes", () => {
    expect(formatFileSize(0)).toBe("0 B");
  });

  it("formats bytes", () => {
    expect(formatFileSize(500)).toBe("500 B");
  });

  it("formats kilobytes", () => {
    expect(formatFileSize(1024)).toBe("1 KB");
  });

  it("formats megabytes", () => {
    expect(formatFileSize(1048576)).toBe("1 MB");
  });

  it("formats gigabytes with decimal", () => {
    expect(formatFileSize(1610612736)).toBe("1.5 GB");
  });

  it("formats terabytes", () => {
    expect(formatFileSize(1099511627776)).toBe("1 TB");
  });
});

describe("formatDate", () => {
  it("returns empty string for undefined", () => {
    expect(formatDate(undefined)).toBe("");
  });

  it("returns empty string for empty string", () => {
    expect(formatDate("")).toBe("");
  });

  it("formats a valid date string", () => {
    const result = formatDate("2024-01-15");
    expect(result).toBeTruthy();
    expect(typeof result).toBe("string");
  });

  it("returns Invalid Date string for unparseable dates", () => {
    expect(formatDate("not-a-date")).toBe("Invalid Date");
  });
});

describe("getResolutionLabel", () => {
  it("returns null for very small resolutions", () => {
    expect(getResolutionLabel(100, 80)).toBeNull();
  });

  it("returns 144p for 144 height", () => {
    expect(getResolutionLabel(256, 144)).toBe("144p");
  });

  it("returns 240p", () => {
    expect(getResolutionLabel(426, 240)).toBe("240p");
  });

  it("returns 360p", () => {
    expect(getResolutionLabel(640, 360)).toBe("360p");
  });

  it("returns 480p", () => {
    expect(getResolutionLabel(854, 480)).toBe("480p");
  });

  it("returns 540p", () => {
    expect(getResolutionLabel(960, 540)).toBe("540p");
  });

  it("returns 720p", () => {
    expect(getResolutionLabel(1280, 720)).toBe("720p");
  });

  it("returns 1080p", () => {
    expect(getResolutionLabel(1920, 1080)).toBe("1080p");
  });

  it("returns 1440p", () => {
    expect(getResolutionLabel(2560, 1440)).toBe("1440p");
  });

  it("returns 4K", () => {
    expect(getResolutionLabel(3840, 1920)).toBe("4K");
  });

  it("returns 8K", () => {
    expect(getResolutionLabel(7680, 4320)).toBe("8K");
  });

  it("handles portrait orientation (width < height)", () => {
    expect(getResolutionLabel(1080, 1920)).toBe("1080p");
  });

  it("returns HUGE for very large resolutions", () => {
    expect(getResolutionLabel(12000, 8000)).toBe("HUGE");
  });
});
