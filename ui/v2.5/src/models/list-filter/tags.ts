import {
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
} from "./criteria/criterion";
import { TagIsMissingCriterionOption } from "./criteria/is-missing";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "name";
// scene markers count has been disabled for now due to performance
// issues
const sortByOptions = [
  "name",
  "random",
  /* "scene_markers_count" */
]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "gallery_count",
      value: "galleries_count",
    },
    {
      messageID: "image_count",
      value: "images_count",
    },
    {
      messageID: "performer_count",
      value: "performers_count",
    },
    {
      messageID: "scene_count",
      value: "scenes_count",
    },
  ]);

const displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
const criterionOptions = [
  createMandatoryStringCriterionOption("name"),
  TagIsMissingCriterionOption,
  createMandatoryNumberCriterionOption("scene_count"),
  createMandatoryNumberCriterionOption("image_count"),
  createMandatoryNumberCriterionOption("gallery_count"),
  createMandatoryNumberCriterionOption("performer_count"),
  // marker count has been disabled for now due to performance issues
  // ListFilterModel.createCriterionOption("marker_count"),
];

export const TagListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
