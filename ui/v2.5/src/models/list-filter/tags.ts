import { createCriterionOption } from "./criteria/criterion";
import { TagIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "name";
// scene markers count has been disabled for now due to performance
// issues
const sortByOptions = [
  "name",
  "scenes_count",
  "images_count",
  "galleries_count",
  "performers_count",
  "random",
  /* "scene_markers_count" */
];
const displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
const criterionOptions = [
  NoneCriterionOption,
  TagIsMissingCriterionOption,
  createCriterionOption("scene_count"),
  createCriterionOption("image_count"),
  createCriterionOption("gallery_count"),
  createCriterionOption("performer_count"),
  // marker count has been disabled for now due to performance issues
  // ListFilterModel.createCriterionOption("marker_count"),
];

export const TagListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
