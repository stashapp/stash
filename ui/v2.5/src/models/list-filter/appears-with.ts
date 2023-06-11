import { createMandatoryNumberCriterionOption } from "./criteria/criterion";
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

const AppearsWithListCriterionOptions = [
  createMandatoryNumberCriterionOption("appears_with_scene_count"),
  createMandatoryNumberCriterionOption("appears_with_gallery_count"),
  createMandatoryNumberCriterionOption("appears_with_image_count"),
  ...PerformerListFilterOptions.criterionOptions,
];

const criterionOptions = AppearsWithListCriterionOptions.filter(
  (option) => option.type !== "scene_count"
)
  .filter((option) => option.type !== "gallery_count")
  .filter((option) => option.type !== "image_count");

const sortByOptions = AppearsWithListSortByOptions;

const { defaultSortBy } = PerformerListFilterOptions;

const { displayModeOptions } = PerformerListFilterOptions;

export const AppearsWithListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
