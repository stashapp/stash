import { CriterionOption } from "./criteria/criterion";
import { DisplayMode } from "./types";

interface ISortByOption {
  messageID: string;
  value: string;
}

export const MediaSortByOptions = [
  "title",
  "path",
  "rating",
  "file_mod_time",
  "tag_count",
  "performer_count",
  "random",
];

export const CreatedSortByOptions = ["created_at", "updated_at"];

export class ListFilterOptions {
  public readonly defaultSortBy: string = "";
  public readonly sortByOptions: ISortByOption[] = [];
  public readonly displayModeOptions: DisplayMode[] = [];
  public readonly criterionOptions: CriterionOption[] = [];

  public static createSortBy(value: string) {
    return {
      messageID: value,
      value,
    };
  }

  constructor(
    defaultSortBy: string,
    sortByOptions: ISortByOption[],
    displayModeOptions: DisplayMode[],
    criterionOptions: CriterionOption[]
  ) {
    this.defaultSortBy = defaultSortBy;
    this.sortByOptions = sortByOptions;
    this.displayModeOptions = displayModeOptions;
    this.criterionOptions = criterionOptions;
  }
}
