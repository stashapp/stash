import { createStringCriterionOption } from "./criteria/criterion";
import { GalleryIsMissingCriterionOption } from "./criteria/is-missing";
import { OrganizedCriterionOption } from "./criteria/organized";
import { PerformersCriterionOption } from "./criteria/performers";
import { RatingCriterionOption } from "./criteria/rating";
import { AverageResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterionOption } from "./criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "path";

const sortByOptions = [
  "date",
  "path",
  "file_mod_time",
  "tag_count",
  "performer_count",
  "title",
  "random",
]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "image_count",
      value: "images_count",
    },
  ]);

const displayModeOptions = [
  DisplayMode.Grid,
  DisplayMode.List,
  DisplayMode.Wall,
];

const criterionOptions = [
  createStringCriterionOption("path"),
  RatingCriterionOption,
  OrganizedCriterionOption,
  AverageResolutionCriterionOption,
  GalleryIsMissingCriterionOption,
  TagsCriterionOption,
  createStringCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createStringCriterionOption("performer_count"),
  createStringCriterionOption("image_count"),
  StudiosCriterionOption,
  createStringCriterionOption("url"),
];

export const GalleryListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
