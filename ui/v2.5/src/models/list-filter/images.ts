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
      NoneCriterionOption,
      ListFilterOptions.createCriterionOption("path"),
      RatingCriterionOption,
      OrganizedCriterionOption,
      ListFilterOptions.createCriterionOption("o_counter"),
      ResolutionCriterionOption,
      ImageIsMissingCriterionOption,
      TagsCriterionOption,
      ListFilterOptions.createCriterionOption("tag_count"),
      PerformerTagsCriterionOption,
      PerformersCriterionOption,
      ListFilterOptions.createCriterionOption("performer_count"),
      StudiosCriterionOption,
    ];

    super(
      ImageListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
