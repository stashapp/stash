import {
  Criterion,
  CriterionType,
  ICriterionOption,
  StringCriterion,
} from "./criterion";

export class CountryCriterion extends StringCriterion {
  constructor() {
    super("country", "country");
  }
}

export class CountryCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("performers");
  public value: CriterionType = "country";
}
