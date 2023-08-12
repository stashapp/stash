import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";
import {
  createDateCriterionOption,
  createMandatoryTimestampCriterionOption,
} from "./criteria/criterion";

const defaultSortBy = "id";
const sortByOptions = ["scene_id", "random", "scenes_updated_at"].map(
  ListFilterOptions.createSortBy
);
const displayModeOptions = [DisplayMode.Wall];
const criterionOptions = [
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
  createDateCriterionOption("scene_date"),
  createMandatoryTimestampCriterionOption("scene_created_at"),
  createMandatoryTimestampCriterionOption("scene_updated_at"),
];

export const SceneFilterListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
