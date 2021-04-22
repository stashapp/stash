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

export class ImageListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "path";

  constructor() {
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
      new NoneCriterionOption(),
      ListFilterOptions.createCriterionOption("path"),
      new RatingCriterionOption(),
      new OrganizedCriterionOption(),
      ListFilterOptions.createCriterionOption("o_counter"),
      new ResolutionCriterionOption(),
      new ImageIsMissingCriterionOption(),
      new TagsCriterionOption(),
      ListFilterOptions.createCriterionOption("tag_count"),
      new PerformerTagsCriterionOption(),
      new PerformersCriterionOption(),
      ListFilterOptions.createCriterionOption("performer_count"),
      new StudiosCriterionOption(),
    ];

    super(
      ImageListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
