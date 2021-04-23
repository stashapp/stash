import { CriterionOption, StringCriterion } from "./criterion";

export class CountryCriterionOption extends CriterionOption {
  constructor() {
    super("country", "country");
  }
}

export class CountryCriterion extends StringCriterion {
  constructor() {
    super(new CountryCriterionOption());
  }
}
