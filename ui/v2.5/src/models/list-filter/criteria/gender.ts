import { CriterionModifier } from "src/core/generated-graphql";
import { getGenderStrings } from "src/core/StashService";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  StringCriterion,
} from "./criterion";

export class GenderCriterion extends StringCriterion {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [];

  constructor() {
    super("gender", "gender", getGenderStrings());
  }
}

export class GenderCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("gender");
  public value: CriterionType = "gender";
}
