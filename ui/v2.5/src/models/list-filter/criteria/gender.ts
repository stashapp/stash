import { CriterionModifier, GenderCriterionInput, GenderEnum } from "src/core/generated-graphql";
import { genderStrings, stringToGender } from "src/utils/gender";
import { CriterionOption, MultiStringCriterion } from "./criterion";

export const GenderCriterionOption = new CriterionOption({
  messageID: "gender",
  type: "gender",
  options: genderStrings,
  modifierOptions: [
    CriterionModifier.Equals,
    CriterionModifier.IsNull,
  ],
  makeCriterion: () => new GenderCriterion(),
});

export class GenderCriterion extends MultiStringCriterion {
  constructor() {
    super(GenderCriterionOption);
  }

  protected toCriterionInput(): GenderCriterionInput {
    const value = this.value.map((v) =>
      stringToGender(v)
    ) as GenderEnum[];

    return {
      value,
      modifier: this.modifier,
    };
  }
}
