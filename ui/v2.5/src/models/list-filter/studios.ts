import { StudioIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { RatingCriterionOption } from "./criteria/rating";
import { ParentStudiosCriterionOption } from "./criteria/studios";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

export class StudioListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "name";

  constructor() {
    const sortByOptions = [
      "name",
      "scenes_count",
      "images_count",
      "galleries_count",
      "random",
      "rating",
    ];

    const displayModeOptions = [DisplayMode.Grid];
    const criterionOptions = [
      new NoneCriterionOption(),
      new ParentStudiosCriterionOption(),
      new StudioIsMissingCriterionOption(),
      new RatingCriterionOption(),
      ListFilterOptions.createCriterionOption("scene_count"),
      ListFilterOptions.createCriterionOption("image_count"),
      ListFilterOptions.createCriterionOption("gallery_count"),
      ListFilterOptions.createCriterionOption("url"),
      ListFilterOptions.createCriterionOption("stash_id"),
    ];

    super(
      StudioListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
