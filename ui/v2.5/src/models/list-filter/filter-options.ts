import {
  Criterion,
  CriterionOption,
  CriterionType,
} from "./criteria/criterion";
import { DisplayMode } from "./types";

export class ListFilterOptions {
  public defaultSortBy: string = "";
  public sortByOptions: string[] = [];
  public displayModeOptions: DisplayMode[] = [];
  public criterionOptions: CriterionOption[] = [];

  protected static createCriterionOption(criterion: CriterionType) {
    return new CriterionOption(Criterion.getLabel(criterion), criterion);
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
