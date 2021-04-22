import {
  Criterion,
  CriterionOption,
  CriterionType,
  ICriterionOption,
} from "./criteria/criterion";
import { DisplayMode } from "./types";

export class ListFilterOptions {
  public defaultSortBy: string = "";
  public sortByOptions: string[] = [];
  public displayModeOptions: DisplayMode[] = [];
  public criterionOptions: ICriterionOption[] = [];
  // eslint-disable-next-line @typescript-eslint/no-explicit-any

  protected static createCriterionOption(criterion: CriterionType) {
    return new CriterionOption(Criterion.getLabel(criterion), criterion);
  }

  constructor(
    defaultSortBy: string,
    sortByOptions: string[],
    displayModeOptions: DisplayMode[],
    criterionOptions: ICriterionOption[]
  ) {
    this.defaultSortBy = defaultSortBy;
    this.sortByOptions = [...this.sortByOptions, "created_at", "updated_at"];
    this.displayModeOptions = displayModeOptions;
    this.criterionOptions = criterionOptions;
  }
}
