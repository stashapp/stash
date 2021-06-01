import {
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
} from "./criteria/criterion";
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
].map(ListFilterOptions.createSortBy);

const displayModeOptions = [DisplayMode.Grid, DisplayMode.Wall];
const criterionOptions = [
  NoneCriterionOption,
  createStringCriterionOption("path"),
  RatingCriterionOption,
  OrganizedCriterionOption,
  createMandatoryNumberCriterionOption("o_counter"),
  ResolutionCriterionOption,
  ImageIsMissingCriterionOption,
  TagsCriterionOption,
  createMandatoryNumberCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createMandatoryNumberCriterionOption("performer_count"),
  StudiosCriterionOption,
];
export const ImageListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
