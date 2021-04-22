import { GalleryIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { OrganizedCriterionOption } from "./criteria/organized";
import { PerformersCriterionOption } from "./criteria/performers";
import { RatingCriterionOption } from "./criteria/rating";
import { AverageResolutionCriterionOption } from "./criteria/resolution";
import { StudiosCriterionOption } from "./criteria/studios";
import {
  PerformerTagsCriterionOption,
  TagsCriterionOption,
} from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

export class GalleryListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "path";

  constructor() {
    const sortByOptions = [
      "date",
      "path",
      "file_mod_time",
      "images_count",
      "tag_count",
      "performer_count",
      "title",
      "random",
    ];
    const displayModeOptions = [DisplayMode.Grid, DisplayMode.List, DisplayMode.Wall];
    const criterionOptions = [
      new NoneCriterionOption(),
      ListFilterOptions.createCriterionOption("path"),
      new RatingCriterionOption(),
      new OrganizedCriterionOption(),
      new AverageResolutionCriterionOption(),
      new GalleryIsMissingCriterionOption(),
      new TagsCriterionOption(),
      ListFilterOptions.createCriterionOption("tag_count"),
      new PerformerTagsCriterionOption(),
      new PerformersCriterionOption(),
      ListFilterOptions.createCriterionOption("performer_count"),
      ListFilterOptions.createCriterionOption("image_count"),
      new StudiosCriterionOption(),
      ListFilterOptions.createCriterionOption("url"),
    ];

    super(
      GalleryListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
