import { createCriterionOption } from "./criteria/criterion";
import { MovieIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { StudiosCriterionOption } from "./criteria/studios";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "name";

const sortByOptions = ["name", "scenes_count", "random"];
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
