import { CriterionModifier } from "src/core/generated-graphql";
import {
  Criterion,
  CriterionType,
  ICriterionOption,
  NumberCriterion,
} from "./criterion";

export class RatingCriterion extends NumberCriterion {
  public modifier = CriterionModifier.Equals;
  public modifierOptions = [
    Criterion.getModifierOption(CriterionModifier.Equals),
    Criterion.getModifierOption(CriterionModifier.NotEquals),
    Criterion.getModifierOption(CriterionModifier.GreaterThan),
    Criterion.getModifierOption(CriterionModifier.LessThan),
    Criterion.getModifierOption(CriterionModifier.IsNull),
    Criterion.getModifierOption(CriterionModifier.NotNull),
  ];

  constructor() {
    super("rating", "rating", [1, 2, 3, 4, 5]);
  }
}

export class RatingCriterionOption implements ICriterionOption {
  public label: string = Criterion.getLabel("rating");
  public value: CriterionType = "rating";
}
