import { NoneCriterionOption } from "./criteria/none";
import { PerformersCriterionOption } from "./criteria/performers";
import { SceneTagsCriterionOption, TagsCriterionOption } from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

export class SceneMarkerListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "title";

  constructor() {
    const sortByOptions = [
      "title",
      "seconds",
      "scene_id",
      "random",
      "scenes_updated_at",
    ];
    const displayModeOptions = [DisplayMode.Wall];
    const criterionOptions = [
      NoneCriterionOption,
      TagsCriterionOption,
      SceneTagsCriterionOption,
      PerformersCriterionOption,
    ];

    super(
      SceneMarkerListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
