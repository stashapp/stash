import { CriterionModifier } from "src/core/generated-graphql";
import { Criterion, CriterionOption, NumberCriterion } from "./criterion";

export class RatingCriterionOption extends CriterionOption {
  constructor() {
    super("rating", "rating");
  }
}

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
    super(new RatingCriterionOption(), [1, 2, 3, 4, 5]);
  }
}
