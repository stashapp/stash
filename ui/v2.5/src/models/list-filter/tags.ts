import { TagIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

export class TagListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "name";

  constructor() {
    // scene markers count has been disabled for now due to performance
    // issues
    const sortByOptions = [
      "name",
      "scenes_count",
      "images_count",
      "galleries_count",
      "performers_count",
      "random",
      /* "scene_markers_count" */
    ];
    const displayModeOptions = [DisplayMode.Grid, DisplayMode.List];
    const criterionOptions = [
      NoneCriterionOption,
      TagIsMissingCriterionOption,
      ListFilterOptions.createCriterionOption("scene_count"),
      ListFilterOptions.createCriterionOption("image_count"),
      ListFilterOptions.createCriterionOption("gallery_count"),
      ListFilterOptions.createCriterionOption("performer_count"),
      // marker count has been disabled for now due to performance issues
      // ListFilterModel.createCriterionOption("marker_count"),
    ];

    super(
      TagListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
