import { createCriterionOption } from "./criteria/criterion";
import { HasMarkersCriterionOption } from "./criteria/has-markers";
import { SceneIsMissingCriterionOption } from "./criteria/is-missing";
import { MoviesCriterionOption } from "./criteria/movies";
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
].map(ListFilterOptions.createSortBy);

const displayModeOptions = [
  DisplayMode.Grid,
  DisplayMode.List,
  DisplayMode.Wall,
  DisplayMode.Tagger,
];

const criterionOptions = [
  NoneCriterionOption,
  createCriterionOption("path"),
  RatingCriterionOption,
  OrganizedCriterionOption,
  createCriterionOption("o_counter"),
  ResolutionCriterionOption,
  createCriterionOption("duration"),
  HasMarkersCriterionOption,
  SceneIsMissingCriterionOption,
  TagsCriterionOption,
  createCriterionOption("tag_count"),
  PerformerTagsCriterionOption,
  PerformersCriterionOption,
  createCriterionOption("performer_count"),
  StudiosCriterionOption,
  MoviesCriterionOption,
  createCriterionOption("url"),
  createCriterionOption("stash_id"),
];

export const SceneListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
