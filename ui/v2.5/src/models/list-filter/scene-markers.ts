import { NoneCriterionOption } from "./criteria/none";
import { PerformersCriterionOption } from "./criteria/performers";
import { SceneTagsCriterionOption, TagsCriterionOption } from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "title";
const sortByOptions = [
  "title",
  "seconds",
  "scene_id",
  "random",
  "scenes_updated_at",
].map(ListFilterOptions.createSortBy);
const displayModeOptions = [DisplayMode.Wall];
const criterionOptions = [
  NoneCriterionOption,
  TagsCriterionOption,
  SceneTagsCriterionOption,
  PerformersCriterionOption,
];

export const SceneMarkerListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
