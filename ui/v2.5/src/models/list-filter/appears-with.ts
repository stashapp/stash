import { PerformerListFilterOptions } from "./performers";
import { ListFilterOptions } from "./filter-options";

const AppearsWithListSortByOptions =
  PerformerListFilterOptions.sortByOptions.map((option) => {
    if (option.value === "scenes_count") {
      return {
        ...option,
        value: "appears_with_scenes",
      };
    } else if (option.value === "images_count") {
      return {
        ...option,
        value: "appears_with_images",
      };
    } else if (option.value === "galleries_count") {
      return {
        ...option,
        value: "appears_with_galleries",
      };
    } else {
      return option;
    }
  });

const sortByOptions = AppearsWithListSortByOptions;

const defaultSortBy = "name";

const { displayModeOptions } = PerformerListFilterOptions;

const { criterionOptions } = PerformerListFilterOptions;

export const AppearsWithListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
