import {
  createNumberCriterionOption,
  createMandatoryNumberCriterionOption,
  createStringCriterionOption,
  createBooleanCriterionOption,
  createDateCriterionOption,
  createMandatoryTimestampCriterionOption,
  NumberCriterionOption,
  NullNumberCriterionOption,
} from "./criteria/criterion";
import { FavoriteCriterionOption } from "./criteria/favorite";
import { GenderCriterionOption } from "./criteria/gender";
import { PerformerIsMissingCriterionOption } from "./criteria/is-missing";
import { StashIDCriterionOption } from "./criteria/stash-ids";
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
  "measurements",
  "fake_tits",
  "career_length",
  "tattoos",
  "piercings",
  "aliases",
];

const criterionOptions = [
  FavoriteCriterionOption,
  GenderCriterionOption,
  PerformerIsMissingCriterionOption,
  TagsCriterionOption,
  StudiosCriterionOption,
  StashIDCriterionOption,
  createStringCriterionOption("url"),
  new NullNumberCriterionOption("rating", "rating100"),
  createMandatoryNumberCriterionOption("tag_count"),
  createMandatoryNumberCriterionOption("scene_count"),
  createMandatoryNumberCriterionOption("image_count"),
  createMandatoryNumberCriterionOption("gallery_count"),
  createBooleanCriterionOption("ignore_auto_tag"),
  new NumberCriterionOption("height", "height_cm", "height_cm"),
  ...numberCriteria.map((c) => createNumberCriterionOption(c)),
  ...stringCriteria.map((c) => createStringCriterionOption(c)),
  createDateCriterionOption("birthdate"),
  createDateCriterionOption("death_date"),
  createMandatoryTimestampCriterionOption("created_at"),
  createMandatoryTimestampCriterionOption("updated_at"),
];
export const PerformerListFilterOptions = new ListFilterOptions(
  defaultSortBy,
  sortByOptions,
  displayModeOptions,
  criterionOptions
);
