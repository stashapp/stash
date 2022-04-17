import {
  createNumberCriterionOption,
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
  createBooleanCriterionOption,
} from "./criteria/criterion";
import { FavoriteCriterionOption } from "./criteria/favorite";
import { GenderCriterionOption } from "./criteria/gender";
import { PerformerIsMissingCriterionOption } from "./criteria/is-missing";
import { RatingCriterionOption } from "./criteria/rating";
import { StudiosCriterionOption } from "./criteria/studios";
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
    {
      messageID: "image_count",
      value: "images_count",
    },
    {
      messageID: "gallery_count",
      value: "galleries_count",
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
  "name",
  "details",
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
  FavoriteCriterionOption,
  GenderCriterionOption,
  PerformerIsMissingCriterionOption,
  TagsCriterionOption,
  RatingCriterionOption,
  StudiosCriterionOption,
  createStringCriterionOption("url"),
  createMandatoryNumberCriterionOption("tag_count"),
  createMandatoryNumberCriterionOption("scene_count"),
  createMandatoryNumberCriterionOption("image_count"),
  createMandatoryNumberCriterionOption("gallery_count"),
  createBooleanCriterionOption("ignore_auto_tag"),
  ...numberCriteria.map((c) => createNumberCriterionOption(c)),
  ...stringCriteria.map((c) => createStringCriterionOption(c)),
];
export const PerformerListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
