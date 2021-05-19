import {
  CriterionModifier,
  GenderCriterionInput,
} from "src/core/generated-graphql";
import { getGenderStrings, stringToGender } from "src/core/StashService";
import { CriterionOption, StringCriterion } from "./criterion";

export const GenderCriterionOption = new CriterionOption("gender", "gender");

export class GenderCriterion extends StringCriterion {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];

  constructor() {
    super(GenderCriterionOption, getGenderStrings());
  }

  protected toCriterionInput(): GenderCriterionInput {
    return {
      value: stringToGender(this.value),
      modifier: this.modifier,
    };
  }
}
