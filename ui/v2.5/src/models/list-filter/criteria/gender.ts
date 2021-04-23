import {
  CriterionModifier,
  GenderCriterionInput,
} from "src/core/generated-graphql";
import { getGenderStrings, stringToGender } from "src/core/StashService";
import { CriterionOption, StringCriterion } from "./criterion";

export class GenderCriterionOption extends CriterionOption {
  constructor() {
    super("gender", "gender");
  }
}

export class GenderCriterion extends StringCriterion {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];

  constructor() {
    super(new GenderCriterionOption(), getGenderStrings());
  }

  protected toCriterionInput(): GenderCriterionInput {
    return {
      value: stringToGender(this.value),
      modifier: this.modifier,
    };
  }
}
