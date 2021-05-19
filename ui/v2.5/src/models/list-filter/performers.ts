import { FavoriteCriterionOption } from "./criteria/favorite";
import { GenderCriterionOption } from "./criteria/gender";
import { PerformerIsMissingCriterionOption } from "./criteria/is-missing";
import { NoneCriterionOption } from "./criteria/none";
import { TagsCriterionOption } from "./criteria/tags";
import { ListFilterOptions } from "./filter-options";
import { CriterionType, DisplayMode } from "./types";

export class PerformerListFilterOptions extends ListFilterOptions {
  public static readonly defaultSortBy = "name";

  constructor() {
    const sortByOptions = [
      "name",
      "height",
      "birthdate",
      "scenes_count",
      "tag_count",
      "random",
      "rating",
    ];
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
      ListFilterOptions.createCriterionOption("url"),
      ListFilterOptions.createCriterionOption("tag_count"),
      ListFilterOptions.createCriterionOption("scene_count"),
      ListFilterOptions.createCriterionOption("image_count"),
      ListFilterOptions.createCriterionOption("gallery_count"),
      ...numberCriteria
        .concat(stringCriteria)
        .map((c) => ListFilterOptions.createCriterionOption(c)),
    ];

    super(
      PerformerListFilterOptions.defaultSortBy,
      sortByOptions,
      displayModeOptions,
      criterionOptions
    );
  }
}
