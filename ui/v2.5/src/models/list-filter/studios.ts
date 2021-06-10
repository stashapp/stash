import { createCriterionOption } from "./criteria/criterion";
import { StudioIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
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
  NoneCriterionOption,
  ParentStudiosCriterionOption,
  StudioIsMissingCriterionOption,
  RatingCriterionOption,
  createCriterionOption("scene_count"),
  createCriterionOption("image_count"),
  createCriterionOption("gallery_count"),
  createCriterionOption("url"),
  createCriterionOption("stash_id"),
];

export const StudioListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
