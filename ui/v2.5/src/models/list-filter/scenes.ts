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

export class SceneListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "date";

  constructor() {
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
    ];

    const displayModeOptions = [
      DisplayMode.Grid,
      DisplayMode.List,
      DisplayMode.Wall,
      DisplayMode.Tagger,
    ];

    const criterionOptions = [
      new NoneCriterionOption(),
      ListFilterOptions.createCriterionOption("path"),
      new RatingCriterionOption(),
      new OrganizedCriterionOption(),
      ListFilterOptions.createCriterionOption("o_counter"),
      new ResolutionCriterionOption(),
      ListFilterOptions.createCriterionOption("duration"),
      new HasMarkersCriterionOption(),
      new SceneIsMissingCriterionOption(),
      new TagsCriterionOption(),
      ListFilterOptions.createCriterionOption("tag_count"),
      new PerformerTagsCriterionOption(),
      new PerformersCriterionOption(),
      ListFilterOptions.createCriterionOption("performer_count"),
      new StudiosCriterionOption(),
      new MoviesCriterionOption(),
      ListFilterOptions.createCriterionOption("url"),
      ListFilterOptions.createCriterionOption("stash_id"),
    ];

    super(
      SceneListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
