import {
  createBooleanCriterionOption,
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
  createStringCriterionOption,
  MandatoryNumberCriterionOption,
  createMandatoryTimestampCriterionOption,
} from "./criteria/criterion";
import { TagIsMissingCriterionOption } from "./criteria/is-missing";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";
import {
  ChildTagsCriterionOption,
  ParentTagsCriterionOption,
} from "./criteria/tags";
import { FavoriteTagCriterionOption } from "./criteria/favorite";

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
      messageID: "group_count",
      value: "groups_count",
    },
    {
      messageID: "marker_count",
      value: "scene_markers_count",
    },
    {
      messageID: "studio_count",
      value: "studios_count",
    },
  ]);

const displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
const criterionOptions = [
  FavoriteTagCriterionOption,
  createMandatoryStringCriterionOption("name"),
  TagIsMissingCriterionOption,
  createStringCriterionOption("aliases"),
  createStringCriterionOption("description"),
  createBooleanCriterionOption("ignore_auto_tag"),
  createMandatoryNumberCriterionOption("scene_count"),
  createMandatoryNumberCriterionOption("image_count"),
  createMandatoryNumberCriterionOption("gallery_count"),
  createMandatoryNumberCriterionOption("performer_count"),
  createMandatoryNumberCriterionOption("studio_count"),
  createMandatoryNumberCriterionOption("group_count"),
  createMandatoryNumberCriterionOption("marker_count"),
  ParentTagsCriterionOption,
  new MandatoryNumberCriterionOption("parent_tag_count", "parent_count"),
  ChildTagsCriterionOption,
  new MandatoryNumberCriterionOption("sub_tag_count", "child_count"),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
];

export const TagListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
