import { IntlShape } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { getCountryByISO } from "src/utils/country";
import { StringCriterion, StringCriterionOption } from "./criterion";

export const CountryCriterionOption = new StringCriterionOption({
  messageID: "country",
  type: "country",
  makeCriterion: () => new CountryCriterion(),
});

export class CountryCriterion extends StringCriterion {
  constructor() {
    super(CountryCriterionOption);
  }

  protected getLabelValue(intl: IntlShape) {
    if (
      this.modifier === CriterionModifier.Equals ||
      this.modifier === CriterionModifier.NotEquals
    ) {
      return getCountryByISO(this.value, intl.locale) ?? this.value;
    }

    return super.getLabelValue(intl);
  }
}
