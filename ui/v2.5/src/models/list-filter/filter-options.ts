import { CriterionOption } from "./criteria/criterion";
import { CriterionType, DisplayMode } from "./types";

export class ListFilterOptions {
  public defaultSortBy: string = "";
  public sortByOptions: string[] = [];
  public displayModeOptions: DisplayMode[] = [];
  public criterionOptions: CriterionOption[] = [];

  protected static createCriterionOption(value: CriterionType) {
    return new CriterionOption(value, value);
  }

  constructor(
    defaultSortBy: string,
    sortByOptions: string[],
    displayModeOptions: DisplayMode[],
    criterionOptions: CriterionOption[]
  ) {
    this.defaultSortBy = defaultSortBy;
    this.sortByOptions = [...sortByOptions, "created_at", "updated_at"];
    this.displayModeOptions = displayModeOptions;
    this.criterionOptions = criterionOptions;
  }
}
