import { createCriterionOption } from "./criteria/criterion";
import { FavoriteCriterionOption } from "./criteria/favorite";
import { GenderCriterionOption } from "./criteria/gender";
import { PerformerIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { TagsCriterionOption } from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { CriterionType, DisplayMode } from "./types";

const defaultSortBy = "name";
const sortByOptions = [
  "name",
  "height",
  "birthdate",
  "tag_count",
  "random",
  "rating",
]
  .map(ListFilterOptions.createSortBy)
  .concat([
    {
      messageID: "scene_count",
      value: "scenes_count",
    },
  ]);

const displayModeOptions = [
  DisplayMode.Grid,
  DisplayMode.List,
  DisplayMode.Tagger,
];

const numberCriteria: CriterionType[] = [
  "birth_year",
  "death_year",
  "age",
  "weight",
];

const stringCriteria: CriterionType[] = [
  "ethnicity",
  "country",
  "hair_color",
  "eye_color",
  "height",
  "measurements",
  "fake_tits",
  "career_length",
  "tattoos",
  "piercings",
  "aliases",
  "stash_id",
];

const criterionOptions = [
  NoneCriterionOption,
  FavoriteCriterionOption,
  GenderCriterionOption,
  PerformerIsMissingCriterionOption,
  TagsCriterionOption,
  createCriterionOption("url"),
  createCriterionOption("tag_count"),
  createCriterionOption("scene_count"),
  createCriterionOption("image_count"),
  createCriterionOption("gallery_count"),
  ...numberCriteria.concat(stringCriteria).map((c) => createCriterionOption(c)),
];
export const PerformerListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
