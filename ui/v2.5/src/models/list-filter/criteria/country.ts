import { IntlShape } from "react-intl";
import { CriterionModifier } from "src/core/generated-graphql";
import { getCountryByISO } from "src/utils/country";
import {
  CriterionOption,
  StringCriterion,
  StringCriterionOption,
} from "./criterion";

export const CountryCriterionOption = new CriterionOption({
  messageID: "country",
  type: "country",
  modifierOptions: StringCriterionOption.modifierOptions,
  defaultModifier: StringCriterionOption.defaultModifier,
  makeCriterion: () => new CountryCriterion(),
  inputType: StringCriterionOption.inputType,
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
