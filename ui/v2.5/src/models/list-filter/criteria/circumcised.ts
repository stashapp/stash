import {
  CircumcisionCriterionInput,
  CircumisedEnum,
  CriterionModifier,
} from "src/core/generated-graphql";
import { circumcisedStrings, stringToCircumcised } from "src/utils/circumcised";
import { ModifierCriterionOption, MultiStringCriterion } from "./criterion";

export const CircumcisedCriterionOption = new ModifierCriterionOption({
  messageID: "circumcised",
  type: "circumcised",
  modifierOptions: [
    CriterionModifier.Includes,
    CriterionModifier.Excludes,
    CriterionModifier.IsNull,
    CriterionModifier.NotNull,
  ],
  defaultModifier: CriterionModifier.Includes,
  options: circumcisedStrings,
  makeCriterion: () => new CircumcisedCriterion(),
});

export class CircumcisedCriterion extends MultiStringCriterion {
  constructor() {
    super(CircumcisedCriterionOption);
  }

  public toCriterionInput(): CircumcisionCriterionInput {
    const value = this.value.map((v) =>
      stringToCircumcised(v)
    ) as CircumisedEnum[];

    return {
      value,
      modifier: this.modifier,
    };
  }
}
