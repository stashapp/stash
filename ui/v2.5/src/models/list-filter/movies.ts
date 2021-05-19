import { MovieIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { StudiosCriterionOption } from "./criteria/studios";
import { ListFilterOptions } from "./filter-options";
import { DisplayMode } from "./types";

export class MovieListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "name";

  constructor() {
    const sortByOptions = ["name", "scenes_count", "random"];
    const displayModeOptions = [DisplayMode.Grid];
    const criterionOptions = [
      NoneCriterionOption,
      StudiosCriterionOption,
      MovieIsMissingCriterionOption,
      ListFilterOptions.createCriterionOption("url"),
    ];

    super(
      MovieListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
