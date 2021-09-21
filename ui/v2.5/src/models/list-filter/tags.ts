import {
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
  createStringCriterionOption,
} from "./criteria/criterion";
import { TagIsMissingCriterionOption } from "./criteria/is-missing";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";
import {
  ChildTagsCriterionOption,
  ParentTagsCriterionOption,
} from "./criteria/tags";

const defaultSortBy = "name";
const sortByOptions = ["name", "random"]
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
    {
      messageID: "marker_count",
      value: "scene_markers_count",
    },
  ]);

const displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
const criterionOptions = [
  createMandatoryStringCriterionOption("name"),
  TagIsMissingCriterionOption,
  createStringCriterionOption("aliases"),
  createMandatoryNumberCriterionOption("scene_count"),
  createMandatoryNumberCriterionOption("image_count"),
  createMandatoryNumberCriterionOption("gallery_count"),
  createMandatoryNumberCriterionOption("performer_count"),
  createMandatoryNumberCriterionOption("marker_count"),
  ParentTagsCriterionOption,
  ChildTagsCriterionOption,
];

export const TagListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
