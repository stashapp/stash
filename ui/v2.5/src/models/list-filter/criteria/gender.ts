import { GenderCriterionInput } from "src/core/generated-graphql";
import { genderStrings, stringToGender } from "src/utils/gender";
import { CriterionOption, StringCriterion } from "./criterion";

export const GenderCriterionOption = new CriterionOption({
  messageID: "gender",
  type: "gender",
  options: genderStrings,
  makeCriterion: () => new GenderCriterion(),
});

export class GenderCriterion extends StringCriterion {
  constructor() {
    super(GenderCriterionOption);
  }

  protected toCriterionInput(): GenderCriterionInput {
    return {
      value: stringToGender(this.value),
      modifier: this.modifier,
    };
  }
}
