import { createCriterionOption } from "./criteria/criterion";
import { ImageIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { OrganizedCriterionOption } from "./criteria/organized";
import { PerformersCriterionOption } from "./criteria/performers";
import { RatingCriterionOption } from "./criteria/rating";
import { ResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterionOption } from "./criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "path";

const sortByOptions = [
  "title",
  "path",
  "rating",
  "o_counter",
  "filesize",
  "file_mod_time",
  "tag_count",
  "performer_count",
  "random",
];

const displayModeOptions = [DisplayMode.Grid, DisplayMode.Wall];
const criterionOptions = [
  NoneCriterionOption,
  createCriterionOption("path"),
  RatingCriterionOption,
  OrganizedCriterionOption,
  createCriterionOption("o_counter"),
  ResolutionCriterionOption,
  ImageIsMissingCriterionOption,
  TagsCriterionOption,
  createCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createCriterionOption("performer_count"),
  StudiosCriterionOption,
];
export const ImageListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
