import { CriterionOption, StringCriterion } from "./criterion";

const countryCriterionOption = new CriterionOption("country", "country");

export class CountryCriterion extends StringCriterion {
  constructor() {
    super(countryCriterionOption);
  }
}
