import {
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
  NullNumberCriterionOption,
  createDateCriterionOption,
  createMandatoryTimestampCriterionOption,
} from "./criteria/criterion";
import { PerformerFavoriteCriterionOption } from "./criteria/favorite";
import { GalleryIsMissingCriterionOption } from "./criteria/is-missing";
import { OrganizedCriterionOption } from "./criteria/organized";
import { HasChaptersCriterionOption } from "./criteria/has-chapters";
import { PerformersCriterionOption } from "./criteria/performers";
import { AverageResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterionOption } from "./criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions, MediaSortByOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "path";

const sortByOptions = ["date", ...MediaSortByOptions]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "image_count",
      value: "images_count",
    },
    {
      messageID: "zip_file_count",
      value: "file_count",
    },
  ]);

const displayModeOptions = [
  DisplayMode.Grid,
  DisplayMode.List,
  DisplayMode.Wall,
];

const criterionOptions = [
  createStringCriterionOption("title"),
  createStringCriterionOption("details"),
  createStringCriterionOption("path"),
  createStringCriterionOption(
    "galleryChecksum",
    "media_info.checksum",
    "checksum"
  ),
  new NullNumberCriterionOption("rating", "rating100"),
  OrganizedCriterionOption,
  AverageResolutionCriterionOption,
  GalleryIsMissingCriterionOption,
  TagsCriterionOption,
  HasChaptersCriterionOption,
  createStringCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createStringCriterionOption("performer_count"),
  createMandatoryNumberCriterionOption("performer_age"),
  PerformerFavoriteCriterionOption,
  createStringCriterionOption("image_count"),
  StudiosCriterionOption,
  createStringCriterionOption("url"),
  createMandatoryNumberCriterionOption("file_count", "zip_file_count"),
  createDateCriterionOption("date"),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
];

export const GalleryListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
