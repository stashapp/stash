import { createCriterionOption } from "./criteria/criterion";
import { MovieIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { StudiosCriterionOption } from "./criteria/studios";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "name";

const sortByOptions = ["name", "random"]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "scene_count",
      value: "scenes_count",
    },
  ]);
const displayModeOptions = [DisplayMode.Grid];
const criterionOptions = [
  NoneCriterionOption,
  StudiosCriterionOption,
  MovieIsMissingCriterionOption,
  createCriterionOption("url"),
];

export const MovieListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
