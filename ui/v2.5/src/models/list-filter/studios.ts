import {
  createBooleanCriterionOption,
  createMandatoryNumberCriterionOption,
  createMandatoryStringCriterionOption,
  createStringCriterionOption,
} from "./criteria/criterion";
import { StudioIsMissingCriterionOption } from "./criteria/is-missing";
import { RatingCriterionOption } from "./criteria/rating";
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
  ]);

const displayModeOptions = [DisplayMode.Grid];
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
  createStringCriterionOption("stash_id"),
  createStringCriterionOption("aliases"),
];

export const StudioListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
