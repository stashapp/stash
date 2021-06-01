import {
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
} from "./criteria/criterion";
import { HasMarkersCriterionOption } from "./criteria/has-markers";
import { SceneIsMissingCriterionOption } from "./criteria/is-missing";
import { MoviesCriterionOption } from "./criteria/movies";
import { NoneCriterionOption } from "./criteria/none";
import { OrganizedCriterionOption } from "./criteria/organized";
import { PerformersCriterionOption } from "./criteria/performers";
import { RatingCriterionOption } from "./criteria/rating";
import { ResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterionOption } from "./criteria/studios";
import { InteractiveCriterionOption } from "./criteria/interactive";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

const defaultSortBy = "date";
const sortByOptions = [
  "title",
  "path",
  "rating",
  "organized",
  "o_counter",
  "date",
  "filesize",
  "file_mod_time",
  "duration",
  "framerate",
  "bitrate",
  "tag_count",
  "performer_count",
  "random",
  "movie_scene_number",
  "interactive",
].map(ListFilterOptions.createSortBy);

const displayModeOptions = [
  DisplayMode.Grid,
  DisplayMode.List,
  DisplayMode.Wall,
  DisplayMode.Tagger,
];

const criterionOptions = [
  NoneCriterionOption,
  createStringCriterionOption("path"),
  RatingCriterionOption,
  OrganizedCriterionOption,
  createMandatoryNumberCriterionOption("o_counter"),
  ResolutionCriterionOption,
  createMandatoryNumberCriterionOption("duration"),
  HasMarkersCriterionOption,
  SceneIsMissingCriterionOption,
  TagsCriterionOption,
  createMandatoryNumberCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createMandatoryNumberCriterionOption("performer_count"),
  StudiosCriterionOption,
  MoviesCriterionOption,
  createStringCriterionOption("url"),
  createStringCriterionOption("stash_id"),
  InteractiveCriterionOption,
];

export const SceneListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
