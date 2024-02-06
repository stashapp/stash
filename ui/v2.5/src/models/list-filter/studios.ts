import {
  createBooleanCriterionOption,
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
  createStringCriterionOption,
  createMandatoryTimestampCriterionOption,
} from "./criteria/criterion";
import { StudioIsMissingCriterionOption } from "./criteria/is-missing";
import { RatingCriterionOption } from "./criteria/rating";
import { StashIDCriterionOption } from "./criteria/stash-ids";
import { ParentStudiosCriterionOption } from "./criteria/studios";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "name";
const sortByOptions = ["name", "random", "rating"]
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
      messageID: "scene_count",
      value: "scenes_count",
    },
    {
      messageID: "subsidiary_studio_count",
      value: "child_count",
    },
  ]);

const displayModeOptions = [DisplayMode.Grid, DisplayMode.Tagger];
const criterionOptions = [
  createMandatoryStringCriterionOption("name"),
  createStringCriterionOption("details"),
  ParentStudiosCriterionOption,
  StudioIsMissingCriterionOption,
  RatingCriterionOption,
  createBooleanCriterionOption("ignore_auto_tag"),
  createMandatoryNumberCriterionOption("scene_count"),
  createMandatoryNumberCriterionOption("image_count"),
  createMandatoryNumberCriterionOption("gallery_count"),
  createStringCriterionOption("url"),
  StashIDCriterionOption,
  createStringCriterionOption("aliases"),
  createMandatoryNumberCriterionOption(
    "child_count",
    "subsidiary_studio_count"
  ),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
];

export const StudioListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
