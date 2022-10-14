import { getCountryByISO } from "src/utils";
import { StringCriterion, StringCriterionOption } from "./criterion";

const countryCriterionOption = new StringCriterionOption(
  "country",
  "country",
  "country"
);

export class CountryCriterion extends StringCriterion {
  constructor() {
    super(countryCriterionOption);
  }

  public getLabelValue() {
    return getCountryByISO(this.value) ?? this.value;
  }
}
