import { CriterionModifier } from "src/core/generated-graphql";
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
    if (
      this.modifier === CriterionModifier.Equals ||
      this.modifier === CriterionModifier.NotEquals
    ) {
      return getCountryByISO(this.value) ?? this.value;
    }

    return super.getLabelValue();
  }
}
